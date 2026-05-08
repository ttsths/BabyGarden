# 小园子的数字花园 — 自动化工作流 (WORKFLOW.md)

> 本文件定义了 BabyGarden 项目的自动化工作流策略。
> 包含：Bug Fix 流水线 + Feature 开发流水线
> Symphony-inspired: 仓库内配置，版本随代码演进。

---
tracker:
  kind: github
  repo: ttsths/BabyGarden
  # Bug Fix 触发标签
  bug_labels: [bug, agent]
  # Feature 开发触发标签
  feature_labels: [feature, agent]
  # 只有 OPEN 状态的 Issue 才触发
  active_states: [OPEN]
  # 这些状态的 Issue 会被清理工作区
  terminal_states: [CLOSED, MERGED]

poll:
  # 扫描间隔（秒）
  interval_seconds: 900
  # 检查 WORKFLOW.md 变化的间隔
  hot_reload_check_seconds: 60

agent:
  # 最大并发数（防止同时起太多 Agent）
  max_concurrent_agents: 2
  # 单个 Issue 最大轮次
  max_turns: 10
  # 单轮最大耗时（毫秒）
  turn_timeout_ms: 600000
  # 停滞检测（毫秒，0 禁用）
  stall_timeout_ms: 300000
  # LLM 配置
  coding_provider: kimi-coding
  coding_model: k2p6
  # 审批策略：always / on_change / never
  # Bug Fix: always（每次修复后需审批）
  # Feature Dev: on_change（设计阶段需审批，编码后可选）
  approval_policy: always

retry:
  # 指数退避重试
  strategy: exponential_backoff
  # 基础延迟（毫秒）
  base_ms: 10000
  # 最大退避上限（毫秒）
  max_ms: 300000
  # 最大重试次数
  max_attempts: 3

workspace:
  # 工作区根目录（相对路径相对于此 WORKFLOW.md 所在目录）
  root: /tmp/bug-workspaces
  # 工作区生命周期钩子
  hooks:
    after_create: |
      echo "[hook] after_create: 初始化工作区..."
      # 克隆仓库（如果没指定 branch，由 orchestrator 传入）
      if [ -d "{{ .RepoPath }}" ]; then
        cp -r "{{ .RepoPath }}" "{{ .WorkspacePath }}/repo"
      fi

    before_run: |
      echo "[hook] before_run: 拉取最新代码..."
      cd "{{ .WorkspacePath }}/repo"
      git fetch origin main
      git checkout main
      git pull origin main

    after_run: |
      echo "[hook] after_run: 清理临时文件..."

    before_remove: |
      echo "[hook] before_remove: 归档日志..."
      mkdir -p "{{ .RepoPath }}/.symphony/logs"
      cp "{{ .WorkspacePath }}"/*.log "{{ .RepoPath }}/.symphony/logs/" 2>/dev/null || true

observability:
  # 结构化日志
  structured_logs: true
  # 日志保留天数
  log_retention_days: 7
  # 工作区审计
  audit_workspaces: true

---

# Feature 开发流水线

## 触发条件
- **标签**: `feature` + `agent`
- **手动触发**: `lobster run feature-dev --issue N`

## 步骤

### 1. preflight — 环境检查
- 检查 Issue 状态（open）
- 创建隔离工作区 `/tmp/feature-workspaces/issue-N/`
- 克隆仓库，切换到 main
- 检查必要 secrets/环境变量

### 2. design — AI 设计评审（审批门禁）
- 读取 WORKFLOW.md 中的 feature 设计模板
- PI 分析需求，输出技术方案
- 生成 DESIGN.md（架构决策、API 设计、数据模型）
- **哥哥确认设计后继续**

### 3. implement — PI 编码实现
- 基于 DESIGN.md 实现代码
- 遵循项目代码规范
- 运行本地测试（编译/单元测试）
- **重试策略**: 指数退避 10s → 40s → 160s，最多 3 次

### 4. review — 代码审查
- Codex CLI review --base main
- 检查：规范、类型安全、测试覆盖
- 关键问题阻塞，非关键问题记录

### 5. create-pr — 创建 PR
- 推送分支 `feat/issue-N-feature-name`
- 创建 PR 到 main
- PR 描述包含：设计摘要、测试说明、截图（如有 UI）
- **关联 Issue #N**
- **持续监控 CI checks，必须全部通过**

### 6. notify — 通知完成
- 发送汇报消息（Telegram）
- 更新 pending-tasks.md

---

# BabyGarden Bug Fix Agent

你是 **小园子的数字花园** 项目的 Bug 修复 Agent。

## 你的任务

**{{ issue.title }}**

{{ issue.description }}

## 工作流程

你有一个目标：**修复这个 Bug，确保所有测试通过，代码可编译，创建一个 PR**。

你可以自由使用以下能力来完成这个目标：

### 可用工具
1. **gh CLI** — 操作 GitHub Issues / PRs
2. **go / npm** — 编译、运行测试、检查类型
3. **git** — 分支管理、提交、推送
4. **文件读写** — 查看和修改代码

### 工作步骤
你可以自主决定执行顺序，但最终必须产出：

```
1. 📖 分析代码 — 找到 Bug 根源
2. 🔧 实现修复 — 只改必要的文件
3. 🧪 运行测试 — 单元测试 + 编译验证
4. 🔄 如有失败 — 分析 → 修复 → 重测（最多 {{ retry.max_attempts }} 轮）
5. 📦 提交代码 — 在 fix/issue-{{ issue.number }} 分支
6. 📝 创建 PR — 包含完整的修复说明
```

### 约束
- ✅ 只修改必要文件，不改无关代码
- ✅ 保持向后兼容
- ✅ 添加回归测试覆盖
- ✅ Go: go build + go vet 必须通过
- ✅ TypeScript: tsc --noEmit 必须通过
- ❌ 不要修改 WORKFLOW.md 本身
- ❌ 不要修改 CI 配置文件
- ❌ 不要合并或关闭 Issue（由人工完成）

### 项目结构
```
{{ .RepoPath }}/
├── yuanzi-backend/     # Go + Gin + GORM + TiDB
│   ├── handler/        # HTTP 处理器
│   │   └── admin/      # 管理后台 API
│   ├── model/          # 数据模型
│   ├── middleware/      # JWT/CORS 中间件
│   ├── router/         # 路由定义
│   └── e2e/            # E2E 冒烟测试
└── yuanzi-frontend/    # React + TypeScript + Ant Design
    └── src/
        ├── admin/      # 管理后台页面
        ├── pages/      # 用户端页面
        ├── components/ # 通用组件
        └── services/   # API 调用层
```

### 完成后
1. 确认所有测试通过
2. 提交代码到 `fix/issue-{{ issue.number }}-<timestamp>` 分支
3. 创建 PR，关联 Issue #{{ issue.number }}
4. 在你的回复中输出 PR 链接
