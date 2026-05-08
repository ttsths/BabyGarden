---
tracker:
  kind: github
  repo: ttsths/BabyGarden
  # Issue 需要同时有这两个标签才触发
  required_labels: [bug, agent]
  # 只有 OPEN 状态的 Issue 才派发
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
  # 单个 Issue 最大修复轮次
  max_turns: 10
  # 单轮最大耗时（毫秒）
  turn_timeout_ms: 600000
  # 停滞检测（毫秒，0 禁用）
  stall_timeout_ms: 300000
  # LLM 配置
  coding_provider: kimi-coding
  coding_model: k2p6
  # 审批策略：always / on_change / never
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

  # ── 仓库元信息 ──
  title: "小园子的数字花园 — Bug Fix 工作流"
  description: >
    本文件定义了 BabyGarden 项目的自动 Bug 修复策略。
    Symphony-inspired: 仓库内配置，版本随代码演进。
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
