#!/usr/bin/env bash
set -euo pipefail

# wait-for-workflow.sh — Poll a GitHub Actions workflow run until completion
# Usage: wait-for-workflow.sh <owner> <repo> <run-id> [poll-interval] [timeout]
# Exit codes: 0 = success, 1 = failure, 2 = timeout

OWNER="${1:?Usage: wait-for-workflow.sh <owner> <repo> <run-id> [poll] [timeout]}"
REPO="${2:?Usage: wait-for-workflow.sh <owner> <repo> <run-id> [poll] [timeout]}"
RUN_ID="${3:?Usage: wait-for-workflow.sh <owner> <repo> <run-id> [poll] [timeout]}"
POLL="${4:-10}"
TIMEOUT="${5:-600}"

ELAPSED=0

while [ "${ELAPSED}" -lt "${TIMEOUT}" ]; do
  STATUS=$(gh api "repos/${OWNER}/${REPO}/actions/runs/${RUN_ID}" \
    --jq '.status' 2>/dev/null || echo "unknown")

  CONCLUSION=$(gh api "repos/${OWNER}/${REPO}/actions/runs/${RUN_ID}" \
    --jq '.conclusion // empty' 2>/dev/null || echo "")

  if [ "${STATUS}" = "completed" ]; then
    if [ "${CONCLUSION}" = "success" ]; then
      echo "✅ Workflow completed successfully"
      exit 0
    else
      echo "❌ Workflow completed with conclusion: ${CONCLUSION}" >&2
      exit 1
    fi
  fi

  echo "   Status: ${STATUS} (${ELAPSED}s / ${TIMEOUT}s)"
  sleep "${POLL}"
  ELAPSED=$((ELAPSED + POLL))
done

echo "⏰ Timeout after ${TIMEOUT}s waiting for workflow" >&2
exit 2
