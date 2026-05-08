# 🌸 小园子的数字花园 (BabyGarden)

> 母婴成长记录社区平台 — 记录宝宝每一个珍贵瞬间

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ttsths/BabyGarden)](https://goreportcard.com/report/github.com/ttsths/BabyGarden)
[![Last Commit](https://img.shields.io/github/last-commit/ttsths/BabyGarden)](https://github.com/ttsths/BabyGarden)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Deploy](https://github.com/ttsths/BabyGarden/actions/workflows/deploy-backend.yml/badge.svg)](https://github.com/ttsths/BabyGarden/actions/workflows/deploy-backend.yml)

---

## 📖 简介

小园子的数字花园是一个面向母婴家庭的成长记录平台。家长可以记录宝宝的成长里程碑、上传照片、管理家庭成员，并通过 AI 助手获取育儿建议。

**后端 API 服务** 基于 Go 语言开发，部署于 Cloudflare Containers，采用 MySQL + Redis 数据存储。

---

## ✨ 当前功能

### 🏠 家庭成员管理
- 创建家庭、邀请成员、移除成员
- 家庭角色管理

### 👶 宝宝管理
- 添加/编辑/删除宝宝信息
- 多宝宝支持（一个家庭多个孩子）

### 📸 照片管理
- 阿里云 OSS 直传
- 照片列表与删除
- 上传回调处理

### 📝 成长记录
- 创建/编辑/删除成长记录
- 支持文字、照片关联
- 时间线浏览

### 🤖 AI 助手
- AI 对话（阿里云 DashScope）
- 语音识别（阿里云 NLS）
- 配额管理

### 📊 统计分析
- 每日/每周数据统计
- 成长趋势可视化数据

### 🔐 认证与安全
- 短信验证码登录（阿里云 SMS）
- JWT Token 认证 + 自动刷新
- 设备注册与管理

### 📡 实时同步
- Server-Sent Events (SSE) 实时推送

### 🏥 健康检查
- `GET /ping` — 存活检查
- `GET /health` — 数据库/Redis 连接状态

### 🤖 自主编码流水线 (Symphony-style)

基于 **OpenAI Symphony** 设计哲学，仓库内配置驱动的自动化流水线。

**双流水线架构：**

| 流水线 | 触发标签 | 关键步骤 | 审批门禁 |
|--------|----------|----------|----------|
| **Bug Fix** | `bug` + `agent` | preflight → analyze → fix → test → review → **approve** → create-pr → notify | 每次修复后 |
| **Feature Dev** | `feature` + `agent` | preflight → **design** → implement → review → create-pr → **CI监控** → notify | 设计阶段 |

**核心特性：**
- **审批门禁**: Bug 修复每次需确认，Feature 设计阶段需确认
- **CI 监控**: Feature PR 创建后持续监控 checks，必须全部通过
- **指数退避**: 测试失败自动重试（10s → 40s → 160s）
- **热加载**: WORKFLOW.md 修改自动生效
- **隔离工作区**: 每 Issue 独立目录

**4-Phase 架构：**

| Phase | 特性 | 说明 |
|-------|------|------|
| **1** | 仓库内配置 | `WORKFLOW.md` (YAML + prompt 模板)，随代码演进 |
| **2** | 指数退避 + 状态核对 | 重试 10s→40s→160s，每个 tick 核对 Issue 状态 |
| **3** | 目标驱动 | analyze+fix+test 合并为 PI K2P6 智能循环 |
| **4** | 热加载 + 多轮续接 | 检测配置变化自动重载，同线程多轮运行 |

---

## 🏗️ 技术栈

| 层 | 技术 |
|----|------|
| **语言** | Go 1.24 |
| **Web 框架** | [Gin](https://github.com/gin-gonic/gin) |
| **ORM** | [GORM](https://gorm.io/) |
| **数据库** | MySQL 8.0 |
| **缓存** | Redis 7 |
| **对象存储** | 阿里云 OSS |
| **短信服务** | 阿里云 SMS |
| **AI 服务** | 阿里云 DashScope + NLS |
| **推送** | APNs (Apple Push) |
| **部署** | Cloudflare Containers (Firecracker) |
| **CI/CD** | GitHub Actions |
| **容器** | Docker (multi-stage build) |
| **API 文档** | Swagger |
| **日志** | Zap + Lumberjack |
| **配置** | Viper |

---

## 📁 项目结构

```
yuanzi-backend/
├── main.go                    # 入口文件
├── go.mod / go.sum            # Go 依赖
├── Dockerfile                 # 多阶段构建
├── wrangler.toml              # Cloudflare 配置
├── docker-compose.yml         # 本地开发环境
├── src/
│   └── index.js               # CF Worker (Container DO 代理)
│
├── config/
│   ├── config.go              # Viper 配置加载
│   └── config.yaml            # 配置文件
│
├── router/
│   └── router.go              # 路由注册 (API 版本化)
│
├── handler/                   # 请求处理器
│   ├── ping.go                # 健康探针
│   ├── health.go              # 健康检查
│   ├── auth.go                # 认证
│   ├── user.go                # 用户
│   ├── family.go              # 家庭
│   ├── baby.go                # 宝宝
│   ├── photo.go               # 照片
│   ├── record.go              # 记录
│   ├── sync.go                # SSE 同步
│   ├── device.go              # 设备
│   ├── stats.go               # 统计
│   ├── ai.go                  # AI 对话
│   └── push_dispatch.go       # 推送调度
│
├── model/                     # 数据模型
│   ├── user.go
│   ├── baby.go
│   ├── family.go
│   ├── photo.go
│   ├── record.go
│   ├── push_device.go
│   ├── ai_chat.go
│   └── verification_code.go
│
├── middleware/                 # 中间件
│   ├── jwt.go                 # JWT 认证
│   ├── dbcheck.go             # 数据库连接检查
│   └── logger.go              # 请求日志
│
├── pkg/                       # 公共包
│   ├── jwt.go                 # JWT 工具
│   ├── db.go                  # 数据库工具
│   ├── ai/                    # AI 服务封装
│   ├── gredis/                # Redis 封装
│   ├── oss/                   # 阿里云 OSS
│   ├── push/                  # APNs 推送
│   ├── sms/                   # 短信服务
│   └── utils/                 # 工具函数
│
├── mysql/
│   ├── mysql.go               # MySQL 连接与迁移
│   ├── init/yuanzi.sql        # 初始化 SQL
│   └── migrations/            # 数据库迁移脚本
│
├── logger/
│   └── logger.go              # 日志系统
│
├── docs/                      # 文档
│   ├── docs.go                # Swagger 生成
│   ├── swagger.yaml
│   ├── SPRINT1_COMPLETE.md
│   └── TASK_BREAKDOWN*.md
│
└── .github/workflows/
    └── deploy-backend.yml     # CI/CD 工作流
```

---

## 🚀 本地开发

### 环境要求

- Go ≥ 1.24
- Docker & Docker Compose
- MySQL 8.0
- Redis 7

### 快速开始

```bash
# 1. 克隆仓库
git clone https://github.com/ttsths/BabyGarden.git
cd BabyGarden/yuanzi-backend

# 2. 启动依赖服务
docker-compose up -d

# 3. 安装 Go 依赖
go mod download

# 4. 配置环境变量 (编辑 config/config.yaml)
#    - MySQL 连接
#    - Redis 连接
#    - 阿里云 AK/SK (OSS, SMS, AI)

# 5. 初始化数据库
mysql -u root -p < mysql/init/yuanzi.sql

# 6. 运行
go run main.go

# 服务启动在 http://localhost:8080
```

### 运行测试

```bash
# 全部测试（含数据库测试）
go test ./... -v

# 仅 handler 测试（快速）
go test ./handler/... -v -count=1 -timeout 5m
```

### API 文档

启动后访问 Swagger UI：
```
http://localhost:8080/swagger/index.html
```

---

## 📦 部署

### 架构

```
GitHub Actions → Docker Build → Cloudflare Containers (Firecracker)
                                    │
                  ┌─────────────────┴─────────────────┐
                  ▼                                   ▼
            Worker (index.js)              Container (Go binary)
         Durable Object proxy               Port 8080
```

### 部署流程

```bash
# 1. 构建 Docker 镜像
docker build -t yuanzi-backend .

# 2. 本地验证
docker run -p 8080:8080 yuanzi-backend

# 3. 部署到 Cloudflare
npx wrangler deploy

# 或通过 GitHub Actions 自动部署
# 推送代码到 main 分支自动触发
```

### CI/CD

推送 `main` 分支（修改 `yuanzi-backend/` 目录）自动触发：

1. **Test** — `go test ./... -v`
2. **Deploy** — Wrangler 部署到 Cloudflare Containers

### 环境变量

| 变量 | 说明 |
|------|------|
| `ENV` | 运行环境 (`development` / `production`) |
| MySQL 配置 | `config/config.yaml` 中配置 |
| Redis 配置 | `config/config.yaml` 中配置 |
| 阿里云 AK/SK | 短信、OSS、AI 服务需要 |

---

## 🔗 API 端点

### 公开端点

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/ping` | 存活探针 |
| `GET` | `/health` | 健康检查 (DB + Redis) |
| `GET` | `/swagger/*` | API 文档 |

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/auth/send-code` | 发送验证码 |
| `POST` | `/api/v1/auth/login` | 验证码登录 |
| `POST` | `/api/v1/auth/refresh` | 刷新 Token |
| `POST` | `/api/v1/auth/logout` | 登出 |

### 用户（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/v1/user/profile` | 获取个人信息 |
| `PUT` | `/api/v1/user/profile` | 更新个人信息 |

### 家庭（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/family` | 创建家庭 |
| `GET` | `/api/v1/family/:id` | 获取家庭信息 |
| `POST` | `/api/v1/family/:id/invite` | 邀请成员 |
| `GET` | `/api/v1/family/:id/members` | 成员列表 |
| `DELETE` | `/api/v1/family/:id/members/:userId` | 移除成员 |

### 宝宝（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/baby` | 添加宝宝 |
| `GET` | `/api/v1/baby` | 宝宝列表 |
| `GET` | `/api/v1/baby/:id` | 宝宝详情 |
| `PUT` | `/api/v1/baby/:id` | 更新宝宝 |
| `DELETE` | `/api/v1/baby/:id` | 删除宝宝 |

### 照片（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/photo/upload-url` | 获取上传 URL |
| `GET` | `/api/v1/photo` | 照片列表 |
| `DELETE` | `/api/v1/photo/:id` | 删除照片 |
| `POST` | `/api/v1/photo/callback` | OSS 上传回调 |

### 记录（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/record` | 创建记录 |
| `GET` | `/api/v1/record` | 记录列表 |
| `GET` | `/api/v1/record/:id` | 记录详情 |
| `PUT` | `/api/v1/record/:id` | 更新记录 |
| `DELETE` | `/api/v1/record/:id` | 删除记录 |

### 其他（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/v1/sync/stream` | SSE 实时流 |
| `POST` | `/api/v1/device/register` | 设备注册 |
| `GET` | `/api/v1/stats/daily` | 每日统计 |
| `GET` | `/api/v1/stats/weekly` | 每周统计 |
| `POST` | `/api/v1/ai/chat` | AI 对话 |
| `POST` | `/api/v1/ai/speech/recognize` | 语音识别 |
| `GET` | `/api/v1/ai/quota` | AI 配额 |

---

## 🗺️ TODO / 路线图

### 短期 (Sprint 2)

- [ ] **前端开发** — 部署到 `babygarden.pages.dev`
- [ ] **用户注册流程** — 完善注册引导
- [ ] **照片瀑布流** — 首页照片展示
- [ ] **推送通知** — APNs 推送集成

### 中期 (Sprint 3)

- [ ] **社交功能** — 亲友圈互动
- [ ] **成长报告** — AI 自动生成周报/月报
- [ ] **视频支持** — 短视频录制与上传
- [ ] **数据导出** — 成长记录导出 PDF

### 长期

- [ ] **移动端 App** — React Native / Flutter
- [ ] **国际化** — 多语言支持
- [ ] **小程序** — 微信小程序版本
- [ ] **智能推荐** — 个性化育儿内容推荐

---

## 🤝 贡献指南

### 分支策略

```
main       — 生产分支（已部署）
├── feat/* — 功能分支
├── fix/*  — 修复分支
└── chore/* — 杂项分支
```

### PR 规范

1. PR 标题使用 Conventional Commits: `feat:`, `fix:`, `chore:`
2. 关联对应的 Issue (如 `Closes #19`)
3. CI 必须通过才能合并
4. 需要至少 1 人 Review

### 自主编码流水线 (Symphony-style)

本项目集成 **Symphony 风格** 目标驱动 Bug 修复流水线。

**仓库内策略：** [`WORKFLOW.md`](WORKFLOW.md) — YAML front matter 定义 tracker/agent/retry 配置，Markdown body 是 prompt 模板。

**5 步 Lobster 流水线：**
```
1. Reconcile  — 加载配置 + 核对 Issue 状态 + 创建工作区
2. Implement  — PI K2P6 目标驱动修复 (analyze→fix→test 智能循环)
3. Approve 🔒 — 人工审批代码变更
4. Create PR  — 创建分支、提交、推送、生成 PR
5. Notify     — 汇报结果 + 归档日志 + 延迟清理工作区
```

**触发方式：**

| 触发 | 方式 |
|------|------|
| 自动 | 创建 `bug` + `agent` 标签 Issue，Scanner 每 15min 扫描 |
| 手动 | `pi: implement #N` |

**安全护栏：**
- ✅ 人工审批门禁（Step 3 强制 approve）
- ✅ Per-issue 工作区隔离（`/tmp/bug-workspaces/`）
- ✅ 指数退避重试（3 次失败停止）
- ✅ 关键文件保护（WORKFLOW.md、CI 配置不可修改）
- ✅ Reconciliation 状态核对（Issue 关闭 → 自动 Released）

### 项目看板

GitHub Issues/PRs 自动同步到 Obsidian Kanban：
```
01-Projects/Projects/小园子的数字花园/Kanban.md
```

---

## 📊 项目状态

| 指标 | 状态 |
|------|------|
| **后端上线** | ✅ 生产运行 |
| **CI/CD** | ✅ 自动部署 |
| **健康检查** | ✅ `/health` OK |
| **数据库** | ✅ MySQL 连接正常 |
| **缓存** | ✅ Redis 连接正常 |
| **前端** | 🚧 待开发 |
| **Symphony 流水线** | ✅ Phase 1-4 改造完成 |

---

## 📄 License

MIT © 2025 鸭子公主 👑
