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

  # ═════════════════════════════════════
  # 模型降级链（按优先级排序）
  # 每个模型失败/限流时自动降级到下一个
  # ═════════════════════════════════════
  model_chain:
    - provider: openai-codex
      model: gpt-5.5
      description: "首选 — 最强编码能力"
    - provider: deepseek
      model: deepseek-v4-pro
      description: "降级 1 — 备用编码能力"
    - provider: zai
      model: glm-5.1
      description: "降级 2 — 国产备用"
    - provider: kimi-coding
      model: k2p6
      description: "兜底 — 最后兜底"

  # 限流检测关键词（模型返回中包含这些即视为限流）
  rate_limit_keywords:
    - "rate limit"
    - "too many requests"
    - "429"
    - "quota exceeded"
    - "请求过于频繁"

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
  # 最大重试次数（每个模型级别）
  max_attempts: 3

workspace:
  # Bug Fix 工作区根目录
  root: /tmp/bug-workspaces
  # Feature 开发工作区根目录
  feature_root: /tmp/feature-workspaces
  # 工作区生命周期钩子
  hooks:
    after_create: |
      echo "[hook] after_create: 初始化工作区..."
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

# ═════════════════════════════════════
# CI 门禁要求
# ═════════════════════════════════════
ci_checks:
  required:
    - name: Codecov
      description: "测试覆盖率报告"
      threshold: "不降低现有覆盖率"
    - name: Backend Quality
      description: "go test + golangci-lint"
      steps:
        - "go test ./... -count=1"
        - "go build ./..."
        - "go vet ./..."
        - "golangci-lint run"
    - name: Frontend Quality
      description: "ESLint + TypeScript 检查"
      steps:
        - "npx eslint src/ --max-warnings 0"
        - "npx tsc --noEmit"

  # 本地预检要求（提交前必须本地通过）
  local_preflight:
    - "cd yuanzi-backend && go build ./... && go vet ./... && go test ./... -count=1 -short"
    - "cd yuanzi-frontend && npx eslint src/ --max-warnings 0 && npx tsc --noEmit"

observability:
  # 结构化日志
  structured_logs: true
  # 日志保留天数
  log_retention_days: 7
  # 工作区审计
  audit_workspaces: true

---

# PI 驱动开发规范

## 核心原则

**所有代码开发都必须通过 PI（Programmatic Intelligence）实现。**
通过设计文档与 PI 沟通，将 Bug Fix 分析或 Feature Design 落到文档中，
告诉 PI 让它实现。减少上下文传递成本。

## 工作模式

```
Issue 创建
  ↓
分析/设计 → 写入 DESIGN.md 或 ANALYSIS.md
  ↓
PI 读取设计文档 → 编码实现
  ↓
本地预检（build + lint + test）
  ↓
PI 修复 CI/Lint 问题（如有）
  ↓
PR 创建
```

## PI 调用规范

1. **优先用文档沟通** — 将需求、设计、约束写入 Markdown 文件，PI 读取后实现
2. **单次任务聚焦** — 每个 PI 调用只做一个明确的任务（修复一个 Bug / 实现一个 Feature）
3. **输出验证** — PI 完成后必须验证编译、lint、测试
4. **失败重试** — 按模型降级链自动重试（见上 `model_chain`）

---

# Bug Fix 流水线 (bug-fix.lobster)

## 触发条件
- **标签**: `bug` + `agent`
- **手动触发**: `lobster run bug-fix --issue N`
- **Cron 扫描**: 每 15 分钟

## 步骤

### 1. preflight — 环境检查 + 本地预检
- 检查 Issue 状态（open）
- 创建隔离工作区 `/tmp/bug-workspaces/issue-N/`
- 克隆仓库，切换到 main
- **本地预检**: go build + go vet + go test + eslint + tsc
  - 确保本地基础校验通过后再让 PI 修复，提高 CI 成功率

### 2. analyze — PI 分析 + 输出分析文档
- 模型优先级: gpt-5.5 → deepseek-v4-pro → glm-5.1 → k2p6
- PI 读取 Issue 和代码，输出 `ANALYSIS.md`
  - Bug 根因分析
  - 影响范围
  - 修复方案
- 限流检测: 检测返回内容是否包含限流关键词，自动降级
- 全部模型失败 → 暂停，提示稍后重试

### 3. implement — PI 编码修复
- 基于 `ANALYSIS.md` 实现修复
- 模型降级链（同上）
- **必须添加测试用例**（后端：单元测试/集成测试；前端：组件测试）
- 本地验证: build + lint + test

### 4. review — 代码审查
- 确定性检查（保护文件、密钥泄露、文件规模）
- 检查后端测试覆盖（新增代码必须有测试）
- 检查前端 ESLint 通过

### 5. approve 🔒 — 人工审批

### 6. create-pr — 创建 PR
- 推送分支 `fix/issue-N-timestamp`
- PR 关联 Issue
- PR 描述包含: ANALYSIS.md 摘要、测试结果

### 7. wait-ci — 等待 CI 通过
- 监控 GitHub Checks: Codecov + Backend Quality + Frontend Quality
- 失败 → 回到 Step 3（PI 修复）

### 8. e2e-cascade — E2E 级联测试
- PR merge + CI 通过后 → 通知 TG
- 自动 spawn subagent 编写/更新 E2E 测试
- 等 Cloudflare 部署后执行 E2E
- E2E 失败 → 自动创建 `bug+agent+e2e` Issue

### 9. notify — 通知结果

---

# Feature 开发流水线 (feature-dev.lobster)

## 触发条件
- **标签**: `feature` + `agent`
- **手动触发**: `lobster run feature-dev --issue N`
- **Cron 扫描**: 每 15 分钟

## 步骤

### 1. preflight — 环境检查 + 本地预检
- 检查 Issue 状态（open）
- 创建隔离工作区 `/tmp/feature-workspaces/issue-N/`
- 克隆仓库，切换到 main
- **本地预检**: 确保当前 main 分支基础校验通过

### 2. design — PI 设计评审（审批门禁）
- 模型优先级: gpt-5.5 → deepseek-v4-pro → glm-5.1 → k2p6
- PI 分析 Feature Issue，输出 `DESIGN.md`:
  - 需求分析
  - 技术方案（架构决策、API 设计、数据模型）
  - 文件变更清单
  - 测试计划
- 限流检测 & 自动降级
- **全部模型失败 → 暂停，提示稍后重试**
- **哥哥确认设计后继续** 🔒

### 3. implement — PI 编码实现
- 基于 `DESIGN.md` 实现代码
- 模型降级链（同上）
- **必须包含测试用例**
- 本地验证: build + lint + test
- 失败 → 指数退避重试

### 4. review — 代码审查
- 确定性检查
- 测试覆盖检查
- CI 规范检查

### 5. create-pr — 创建 PR
- 推送分支 `feat/issue-N-feature-name`
- PR 关联 Issue
- PR 描述包含: DESIGN.md 摘要、实现说明、测试结果
- **持续监控 CI Checks: Codecov + Backend Quality + Frontend Quality**

### 6. merge-and-notify — CI 通过后合并
- Codecov: 覆盖率不降低
- Backend Quality: test + golangci-lint 通过
- Frontend Quality: ESLint + tsc 通过
- 通过后通知 TG

### 7. e2e-cascade — E2E 级联测试
- PR merge + CI 通过 + TG 通知后
- spawn subagent 编写/更新 E2E 冒烟测试
- 等 Cloudflare 部署完成后执行 E2E
- E2E 失败 → 自动创建 `bug+agent+e2e` Issue

---

# E2E 级联测试规范

## 触发时机
- Bug Fix: PR merge + CI 通过 + CF 部署后
- Feature Dev: PR merge + CI 通过 + CF 部署后
- 离线: `lobster run e2e-cascade --issue N`

## 执行流程
```
PR merged → CI all green → TG 通知
  ↓
spawn e2e-subagent
  ├── 后端 E2E: go test ./e2e/ (5+ 条冒烟)
  └── 前端 E2E: npx playwright test (3+ 条冒烟)
  ↓
ALL PASS → ✅ 更新状态
ANY FAIL → 创建 bug+agent+e2e Issue → PI 自动修复
```

## E2E 测试要求
- Bug Fix: 至少 1 条回归测试覆盖修复的 Bug
- Feature Dev: 至少覆盖核心 happy path + 1 条边界条件

---

# BabyGarden Bug Fix Agent Prompt

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
```
1. 📖 分析代码 — 找到 Bug 根源 → 输出 ANALYSIS.md
2. 🔧 实现修复 — 只改必要的文件，添加回归测试
3. 🧪 本地预检 — build + lint + test 全部通过
4. 🔄 如有失败 — 分析 → 修复 → 重测（最多 {{ retry.max_attempts }} 轮）
5. 📦 提交代码 — 在 fix/issue-{{ issue.number }} 分支
6. 📝 创建 PR — 包含 ANALYSIS.md 摘要
```

### 约束
- ✅ 只修改必要文件，不改无关代码
- ✅ 保持向后兼容
- ✅ **必须添加回归测试覆盖**
- ✅ 后端: go build + go vet + go test 必须通过
- ✅ 前端: eslint --max-warnings 0 + tsc --noEmit 必须通过
- ❌ 不要修改 WORKFLOW.md 本身
- ❌ 不要修改 CI 配置文件
- ❌ 不要合并或关闭 Issue（由人工完成）

### CI 门禁
GitHub CI 有三个 Checks 必须通过:
1. **Codecov** — 测试覆盖率不降低
2. **Backend Quality** — go test + golangci-lint
3. **Frontend Quality** — ESLint + tsc --noEmit

**请先在本地通过基础校验，提高 CI 成功率。**

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

---

# BabyGarden Feature Dev Agent Prompt

你是 **小园子的数字花园** 项目的 Feature 开发 Agent。

## 你的任务

**{{ issue.title }}**

{{ issue.description }}

## 工作流程

```
1. 📖 分析需求 — 输出 DESIGN.md（架构决策、API 设计、数据模型、测试计划）
2. 🏗️ 编码实现 — 基于 DESIGN.md，模型降级链自动切换
3. 🧪 本地预检 — build + lint + test 全部通过
4. 📦 提交代码 — feat/issue-{{ issue.number }}-feature-name 分支
5. 📝 创建 PR — 包含 DESIGN.md 摘要
```

### 约束
- ✅ 基于 DESIGN.md 实现，不偏离设计
- ✅ **必须包含单元测试和集成测试**
- ✅ 后端必须 go build + go vet + go test 通过
- ✅ 前端必须 eslint --max-warnings 0 + tsc --noEmit 通过
- ❌ 不要修改 WORKFLOW.md 和 CI 配置

### CI 门禁
1. **Codecov** — 测试覆盖率不降低
2. **Backend Quality** — go test + golangci-lint
3. **Frontend Quality** — ESLint + tsc --noEmit

**请先在本地通过基础校验，提高 CI 成功率。**
