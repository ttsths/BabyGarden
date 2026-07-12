#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scripts/orca-issue-start.sh <issue-number> [--dry-run]

Creates one independent Orca worktree from origin/main, links the GitHub Issue,
and starts Codex with a bounded implementation prompt. The agent may create a
Draft PR but must never merge or deploy production.
EOF
}

if [[ $# -lt 1 || "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

ISSUE_NUMBER="$1"
DRY_RUN=false
if [[ "${2:-}" == "--dry-run" ]]; then
  DRY_RUN=true
elif [[ $# -gt 1 ]]; then
  echo "Unknown option: $2" >&2
  usage >&2
  exit 2
fi

if [[ ! "$ISSUE_NUMBER" =~ ^[0-9]+$ ]]; then
  echo "Issue number must be numeric" >&2
  exit 2
fi

for command_name in git gh jq orca; do
  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "Missing required command: $command_name" >&2
    exit 1
  fi
done

REPO_ROOT="$(git rev-parse --show-toplevel)"
REPO="$(gh repo view --json nameWithOwner --jq .nameWithOwner)"
ISSUE_JSON="$(gh issue view "$ISSUE_NUMBER" --repo "$REPO" \
  --json number,title,body,state,labels,url)"

if [[ "$(jq -r .state <<<"$ISSUE_JSON")" != "OPEN" ]]; then
  echo "Issue #$ISSUE_NUMBER is not open" >&2
  exit 1
fi

TITLE="$(jq -r .title <<<"$ISSUE_JSON")"
BODY="$(jq -r .body <<<"$ISSUE_JSON")"
URL="$(jq -r .url <<<"$ISSUE_JSON")"
LABELS="$(jq -r '[.labels[].name] | join(", ")' <<<"$ISSUE_JSON")"
SLUG="$(printf '%s' "$TITLE" | tr '[:upper:]' '[:lower:]' \
  | sed -E 's/[^[:alnum:]]+/-/g; s/^-+|-+$//g' | cut -c1-38)"
WORKTREE_NAME="issue-${ISSUE_NUMBER}-${SLUG:-task}"

PROMPT="$(cat <<EOF
Implement GitHub Issue #$ISSUE_NUMBER in the BabyGarden repository.

Issue: $URL
Title: $TITLE
Labels: $LABELS

Issue body:
$BODY

Execution contract:
1. Start from the worktree's origin/main baseline and read AGENTS.md,
   WORKFLOW.md, and relevant design/spec documents before editing.
2. Stay inside this Issue's scope. Do not modify unrelated user files.
3. Keep SMS, AI, photos/R2, realtime, push and native App out of MVP unless the
   Issue explicitly says otherwise.
4. Update the Orca worktree comment and workspace status at meaningful gates:
   in-progress while implementing, in-review when the Draft PR is ready.
5. Run the smallest relevant build, lint and tests, then record exact commands
   and results in the Draft PR body.
6. Push only this issue branch and create a Draft PR that references #$ISSUE_NUMBER.
7. Never merge the PR, close dependency Issues, or deploy production.
8. If blocked or a human decision is required, stop safely and write the
   blocker in the Orca comment and GitHub Issue/PR; do not guess.
EOF
)"

if [[ "$DRY_RUN" == true ]]; then
  jq -n \
    --arg repo "$REPO" \
    --arg issue "$ISSUE_NUMBER" \
    --arg worktree "$WORKTREE_NAME" \
    --arg base "origin/main" \
    --arg prompt "$PROMPT" \
    '{dryRun:true, repo:$repo, issue:($issue|tonumber), worktree:$worktree, base:$base, prompt:$prompt}'
  exit 0
fi

if orca worktree show --worktree "issue:$ISSUE_NUMBER" --json >/dev/null 2>&1; then
  echo "Issue #$ISSUE_NUMBER already has an Orca worktree; reuse it instead of creating a duplicate" >&2
  exit 1
fi

orca repo set-base-ref --repo "path:$REPO_ROOT" --ref origin/main --json >/dev/null

orca worktree create \
  --repo "path:$REPO_ROOT" \
  --name "$WORKTREE_NAME" \
  --issue "$ISSUE_NUMBER" \
  --base-branch origin/main \
  --no-parent \
  --agent codex \
  --prompt "$PROMPT" \
  --comment "Issue #$ISSUE_NUMBER: Codex implementing; Draft PR only" \
  --json
