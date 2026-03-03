---
name: release
description: >-
  Release the gh-copilot-scratch extension: commit, push, tag, wait for CI, and
  create a GitHub release via the tag-triggered workflow.
  Use this skill when the user asks to "release", "ship it", "push and release",
  "cut a release", "tag a release", or any variation of committing + pushing + tagging.
---

# Release Skill

Automates the full release pipeline for gh-copilot-scratch: commit → push → tag → CI → release.

## When to use

Any time the user wants to ship changes. This includes partial flows (e.g., "just push and wait for CI") — adapt by running only the relevant steps.

## Scripts

Three shell scripts in `scripts/` handle the mechanical parts:

| Script | Purpose |
|---|---|
| `scripts/release.sh` | Full orchestrated release (push → tag → CI → release) |
| `scripts/find-workflow-run.sh` | Finds the GH Actions run triggered by a specific tag/commit |
| `scripts/wait-for-workflow.sh` | Polls a workflow run until completion (success/failure/timeout) |

## Full release flow

### Step 1: Commit

Use the `git-commit` skill if available, or commit directly. Ensure the `Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>` trailer is included.

### Step 2: Determine version

Check existing tags to determine the next version:

```bash
git tag --sort=-v:refname | head -5
```

Follow semver (v0.x.y for pre-1.0, vX.Y.Z after).

### Step 3: Run the release script

```bash
chmod +x .github/skills/release/scripts/*.sh
.github/skills/release/scripts/release.sh ekroon gh-copilot-scratch main <version>
```

This will:
1. Push to main
2. Create and push the version tag
3. Wait for the Release workflow to complete
4. Report the release URL

### Step 4: Report results

Tell the user the final status — the release tag, the workflow URLs, and whether everything succeeded.

## Running individual scripts

```bash
SCRIPTS=".github/skills/release/scripts"

# Find the workflow run for a tag push
RUN_ID=$("$SCRIPTS/find-workflow-run.sh" ekroon gh-copilot-scratch <sha> release.yml)

# Wait for it to complete (10s poll, 600s timeout)
"$SCRIPTS/wait-for-workflow.sh" ekroon gh-copilot-scratch "$RUN_ID" 10 600
```

## Error handling

- **CI failure**: The release script exits with code 1 and prints the workflow conclusion. Show the user the workflow URL and offer to investigate.
- **Timeout**: Exits with code 2. Default timeout is 600s. Override with `TIMEOUT=900` env var.
- **GPG signing hang**: Use `-c commit.gpgsign=false` when committing.

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `POLL_INTERVAL` | `10` | Seconds between CI status checks |
| `TIMEOUT` | `600` | Max seconds to wait for each workflow |
