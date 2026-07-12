# BabyGarden Orca Issue Workflow

## Purpose

This workflow provides a controlled, observable path from one GitHub Issue to
one isolated Orca worktree and one Draft PR. GitHub Issues remain the source of
scope and acceptance criteria. Orca owns the local worktree, terminal and card
status. A human owns review and merge.

## Safety Boundary

- Trigger is an explicit command. The optional `orca-ready` label signals human
  intent but no background scanner consumes it in v1.
- Every run starts from `origin/main` in an independent top-level worktree.
- The agent may push its issue branch and create a Draft PR.
- The agent must not merge, auto-close dependency Issues or deploy production.
- Existing `bug + agent` and `feature + agent` automations remain separate.
- Maximum recommended concurrency is two active implementation worktrees.

## Start an Issue

```bash
scripts/orca-issue-start.sh 84 --dry-run
scripts/orca-issue-start.sh 84
```

The dry run validates GitHub/Orca access and prints the exact bounded prompt
without creating a worktree or terminal.

## State Model

| Orca status   | Meaning                                            | Required evidence        |
| ------------- | -------------------------------------------------- | ------------------------ |
| `todo`        | Worktree exists but implementation has not started | Linked Issue             |
| `in-progress` | Agent is analyzing or editing                      | Current worktree comment |
| `in-review`   | Draft PR exists                                    | Test results in PR body  |
| `completed`   | Human merged or explicitly ended the task          | Merged/closed reference  |

Update state and the short card comment at meaningful gates:

```bash
orca worktree set --worktree active \
  --workspace-status in-progress \
  --comment "Issue #84: preparing Stitch prompt pack" --json

orca worktree set --worktree active \
  --workspace-status in-review \
  --comment "Draft PR ready; waiting for visual review" --json
```

## Agent Contract

The generated prompt always includes the complete Issue body, labels, URL,
repository rules, MVP exclusions, verification requirements and Draft-PR-only
boundary. If the Issue requires a human choice, the agent records the blocker
instead of expanding scope.

## Draft PR Contract

The PR body must contain:

1. Issue reference and scope summary.
2. What changed and what was intentionally excluded.
3. Exact validation commands and results.
4. Risks, blockers and human review points.
5. `Closes #N` only when merging this PR should close that Issue.

## Failure Recovery

1. Inspect the Orca card comment and terminal output.
2. If the terminal handle is stale, reacquire it with `orca terminal list`.
3. If the branch is useful, keep the worktree and restart Codex in that same
   worktree; do not create a duplicate Issue worktree.
4. If no useful changes exist, remove the worktree only after confirming the
   branch and logs are no longer needed.
5. Never use a force push or automatic merge as recovery.

## Pilot

Issue #84 is the initial low-risk pilot because its first deliverable is a
reviewable Stitch prompt pack rather than production business behavior. The
pilot is successful when the worktree is linked, the prompt contains the Issue
contract, validation is recorded, and a Draft PR reaches `in-review` without
automatic merge.
