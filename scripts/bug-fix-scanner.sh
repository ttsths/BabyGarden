#!/bin/bash
# Bug Fix Scanner — 扫描 GitHub Issues 并触发 Bug Fix 流水线
# 触发条件: 标签 = bug + agent, 状态 = open

set -euo pipefail

REPO="ttsths/BabyGarden"
WORKFLOW_ROOT="/Users/ttsths/Desktop/AGI/AGI/BabyGarden"
LOG_DIR="$WORKFLOW_ROOT/.symphony/logs"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
LOG_FILE="$LOG_DIR/scanner-$TIMESTAMP.log"

# 创建日志目录
mkdir -p "$LOG_DIR"

exec > >(tee -a "$LOG_FILE")
exec 2>&1

echo "=========================================="
echo "🐛 Bug Fix Scanner — $(date)"
echo "=========================================="
echo ""

# 检查 gh CLI
cd "$WORKFLOW_ROOT"

# 获取带 bug + agent 标签的开放 Issue
echo "🔍 扫描 GitHub Issues ($REPO)..."
ISSUES=$(gh issue list --repo "$REPO" --label bug,agent --state open --json number,title,labels,createdAt 2>/dev/null || echo "[]")

# 检查是否有 Issue
COUNT=$(echo "$ISSUES" | jq 'length')

if [ "$COUNT" -eq 0 ]; then
    echo "✅ 未发现需要处理的 Bug Issue（标签: bug + agent）"
    echo "   当前开放 Issue 总数: $(gh issue list --repo "$REPO" --state open --json number | jq 'length')"
    echo ""
    echo "扫描完成 — $(date)"
    exit 0
fi

echo "🚨 发现 $COUNT 个 Bug Issue 需要处理:"
echo "$ISSUES" | jq -r '.[] | "  - Issue #\(.number): \(.title)"'
echo ""

# 逐个处理 Issue（或交给工作流调度器）
# 当前实现仅报告，实际触发由外部 cron/调度器控制
echo "📋 待处理 Issue 列表:"
echo "$ISSUES" | jq -r '.[] | "  Issue #\(.number): \(.title) [\(.createdAt)]"'
echo ""

echo "⚡ 触发方式:"
echo "  lobster run bug-fix --issue <N>"
echo ""

echo "扫描完成 — $(date)"
