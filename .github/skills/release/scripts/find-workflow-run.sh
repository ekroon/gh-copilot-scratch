#!/usr/bin/env bash
set -euo pipefail

# find-workflow-run.sh — Find a GitHub Actions workflow run for a commit
# Usage: find-workflow-run.sh <owner> <repo> <sha> <workflow-file>
# Outputs: run ID

OWNER="${1:?Usage: find-workflow-run.sh <owner> <repo> <sha> <workflow-file>}"
REPO="${2:?Usage: find-workflow-run.sh <owner> <repo> <sha> <workflow-file>}"
SHA="${3:?Usage: find-workflow-run.sh <owner> <repo> <sha> <workflow-file>}"
WORKFLOW="${4:?Usage: find-workflow-run.sh <owner> <repo> <sha> <workflow-file>}"

MAX_ATTEMPTS=30
POLL=5

for i in $(seq 1 "${MAX_ATTEMPTS}"); do
  RUN_ID=$(gh api "repos/${OWNER}/${REPO}/actions/workflows/${WORKFLOW}/runs?head_sha=${SHA}&per_page=1" \
    --jq '.workflow_runs[0].id // empty' 2>/dev/null || echo "")

  if [ -n "${RUN_ID}" ]; then
    echo "${RUN_ID}"
    exit 0
  fi

  sleep "${POLL}"
done

echo "ERROR: No workflow run found for ${SHA} after ${MAX_ATTEMPTS} attempts" >&2
exit 1
