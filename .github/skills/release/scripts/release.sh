#!/usr/bin/env bash
set -euo pipefail

# release.sh — Full release orchestration for gh-copilot-scratch
# Usage: release.sh <owner> <repo> <branch> <version>
# Example: release.sh ekroon gh-copilot-scratch main v0.1.0

OWNER="${1:?Usage: release.sh <owner> <repo> <branch> <version>}"
REPO="${2:?Usage: release.sh <owner> <repo> <branch> <version>}"
BRANCH="${3:?Usage: release.sh <owner> <repo> <branch> <version>}"
VERSION="${4:?Usage: release.sh <owner> <repo> <branch> <version>}"

SCRIPTS_DIR="$(cd "$(dirname "$0")" && pwd)"
POLL_INTERVAL="${POLL_INTERVAL:-10}"
TIMEOUT="${TIMEOUT:-600}"

echo "==> Pushing to ${BRANCH}..."
git push origin "${BRANCH}"

echo "==> Creating tag ${VERSION}..."
git tag "${VERSION}"
git push origin "${VERSION}"

SHA="$(git rev-parse HEAD)"
echo "==> Commit: ${SHA}"

echo "==> Waiting for Release workflow..."
sleep 5  # Give GitHub Actions a moment to pick up the tag push

RUN_ID=$("${SCRIPTS_DIR}/find-workflow-run.sh" "${OWNER}" "${REPO}" "${SHA}" "release.yml")
echo "==> Found workflow run: ${RUN_ID}"
echo "    https://github.com/${OWNER}/${REPO}/actions/runs/${RUN_ID}"

"${SCRIPTS_DIR}/wait-for-workflow.sh" "${OWNER}" "${REPO}" "${RUN_ID}" "${POLL_INTERVAL}" "${TIMEOUT}"

echo ""
echo "==> Release ${VERSION} complete!"
echo "    https://github.com/${OWNER}/${REPO}/releases/tag/${VERSION}"
