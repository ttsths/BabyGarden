# Yuanzi 后端开发 - 任务拆解

> **生成时间**: 2026-03-09
> **基于**: PRD-v1.1 + API文档
> **负责人**: Backend Engineer
> **优先级**: P0 (MVP核心功能)

---

## 任务总览

### Phase 1: 基础搭建 + 核心记录 (Week 1-4)
- [x] 项目初始化结构已就绪
- [ ] 用户系统（注册/登录）
- [ ] 宝宝信息管理
- [ ] 喂养记录核心功能
- [ ] 睡眠追踪
- [ ] 排泄记录

### Phase 2: 同步与共享 (Week 5-8)
- [ ] 照片管理、云端存储
- [ ] 家庭共享、邀请机制
- [ ] SSE数据同步
- [ ] 极光推送集成

---

## 详细任务分解

### 模块 1: 用户系统与认证

#### Task 1.1: 数据库表设计
- [ ] 设计 users 表结构
  - id, phone, password_hash, nickname, avatar_url, created_at, updated_at
- [ ] 设计 babies 表结构
  - id, user_id, nickname, avatar_url, birth_date, gender, created_at, updated_at
- [ ] 设计 family_members 表结构
  - id, family_id, user_id, role, created_at, updated_at
- [ ] 编写迁移脚本: `mysql/migrations/001_users_and_babies.sql`

#### Task 1.2: 用户认证 Handler
- [ ] `handler/auth.go` - 用户认证处理器
  - RegisterHandler: 手机号+验证码注册
  - LoginHandler: 手机号+密码登录
  - GetProfileHandler: 获取用户信息
  - UpdateProfileHandler: 更新用户信息
  - SendSMSHandler: 发送短信验证码

#### Task 1.3: 宝宝管理 Handler
- [ ] `handler/baby.go` - 宝宝信息管理
  - CreateBabyHandler: 创建宝宝档案
  - GetBabyHandler: 获取宝宝信息
  - UpdateBabyHandler: 更新宝宝信息
  - ListBabiesHandler: 获取宝宝列表

#### Task 1.4: 认证中间件
- [ ] `middleware/auth.go` - JWT认证中间件
  - 生成JWT Token
  - 验证JWT Token
  - 从Token中提取用户信息

#### Task 1.5: 路由配置
- [ ] 更新 `router/router.go`
  - POST /api/auth/register
  - POST /api/auth/login
  - GET /api/auth/sms
  - GET /api/user/profile
  - PUT /api/user/profile
  - POST /api/baby
  - GET /api/baby/:id
  - PUT /api/baby/:id
  - GET /api/baby

---

### 模块 2: 喂养记录

#### Task 2.1: 数据库表设计
- [ ] 设计 feeding_records 表结构
  - id, baby_id, user_id, time, type, amount, side, note, created_at, updated_at
- [ ] 编写迁移脚本: `mysql/migrations/002_feeding_records.sql`

#### Task 2.2: 喂养记录 Model
- [ ] `model/feeding.go` - 喂养记录数据模型
  - FeedingRecord 结构体
  - CreateFeedingRecord 方法
  - GetFeedingRecord 方法
  - ListFeedingRecords 方法 (按日期范围查询)
  - GetLastFeedingRecord 方法 (获取最近一次喂养)

#### Task 2.3: 喂养记录 Handler
- [ ] `handler/feeding.go` - 喂养记录处理器
  - CreateFeedingHandler: 创建喂养记录
  - GetFeedingHandler: 获取单条记录
  - ListFeedingHandler: 获取记录列表 (支持分页、日期筛选)
  - GetFeedingStatsHandler: 获取喂养统计 (今日次数、总奶量等)

#### Task 2.4: 路由配置
- [ ] 更新 `router/router.go`
  - POST /api/records/feeding
  - GET /api/records/feeding/:id
  - GET /api/records/feeding
  - GET /api/records/feeding/stats

#### Task 2.5: 喂奶提醒 (极光推送)
- [ ] 集成极光推送SDK
- [ ] 设计 feeding_reminders 表
  - id, baby_id, last_feeding_time, next_remind_time, status
- [ ] 实现定时任务检查并推送
- [ ] Handler: ScheduleRemindHandler, CancelRemindHandler

---

### 模块 3: 睡眠追踪

#### Task 3.1: 数据库表设计
- [ ] 设计 sleep_records 表结构
  - id, baby_id, user_id, start_time, end_time, duration, note, created_at, updated_at
- [ ] 编写迁移脚本: `mysql/migrations/003_sleep_records.sql`

#### Task 3.2: 睡眠记录 Model
- [ ] `model/sleep.go` - 睡眠记录数据模型
  - SleepRecord 结构体
  - CreateSleepRecord 方法
  - UpdateSleepRecord 方法 (更新结束时间)
  - GetSleepRecord 方法
  - ListSleepRecords 方法

#### Task 3.3: 睡眠记录 Handler
- [ ] `handler/sleep.go` - 睡眠记录处理器
  - StartSleepHandler: 开始睡眠记录
  - EndSleepHandler: 结束睡眠记录
  - GetSleepHandler: 获取单条记录
  - ListSleepHandler: 获取记录列表

#### Task 3.4: 路由配置
- [ ] 更新 `router/router.go`
  - POST /api/records/sleep/start
  - POST /api/records/sleep/end/:id
  - GET /api/records/sleep/:id
  - GET /api/records/sleep

---

### 模块 4: 排泄记录

#### Task 4.1: 数据库表设计
- [ ] 设计 diaper_records 表结构
  - id, baby_id, user_id, time, type (pee/poo/mixed), note, photo_url, created_at, updated_at
- [ ] 编写迁移脚本: `mysql/migrations/004_diaper_records.sql`

#### Task 4.2: 排泄记录 Model
- [ ] `model/diaper.go` - 排泄记录数据模型
  - DiaperRecord 结构体
  - CreateDiaperRecord 方法
  - GetDiaperRecord 方法
  - ListDiaperRecords 方法

#### Task 4.3: 排泄记录 Handler
- [ ] `handler/diaper.go` - 排泄记录处理器
  - CreateDiaperHandler: 创建排泄记录
  - GetDiaperHandler: 获取单条记录
  - ListDiaperHandler: 获取记录列表

#### Task 4.4: 路由配置
- [ ] 更新 `router/router.go`
  - POST /api/records/diaper
  - GET /api/records/diaper/:id
  - GET /api/records/diaper

---

### 模块 5: 照片管理

#### Task 5.1: 数据库表设计
- [ ] 设计 photos 表结构
  - id, baby_id, user_id, url, thumb_url, caption, created_at, updated_at
- [ ] 编写迁移脚本: `mysql/migrations/005_photos.sql`

#### Task 5.2: OSS集成 (阿里云)
- [ ] 配置阿里云OSS
- [ ] 实现文件上传服务 `pkg/oss.go`
  - UploadFile 方法
  - DeleteFile 方法
  - GenerateThumbnail 方法

#### Task 5.3: 照片 Handler
- [ ] `handler/photo.go` - 照片处理器
  - UploadPhotoHandler: 上传照片
  - ListPhotosHandler: 获取照片列表 (支持分页)
  - GetPhotoHandler: 获取单张照片
  - DeletePhotoHandler: 删除照片

#### Task 5.4: 路由配置
- [ ] 更新 `router/router.go`
  - POST /api/photos/upload
  - GET /api/photos
  - GET /api/photos/:id
  - DELETE /api/photos/:id

---

### 模块 6: 家庭共享

#### Task 6.1: 数据库表设计
- [ ] 设计 families 表结构
  - id, name, invite_code, created_at, updated_at
- [ ] 设计 family_invites 表结构
  - id, family_id, inviter_id, invitee_phone, code, status, expires_at, created_at
- [ ] 编写迁移脚本: `mysql/migrations/006_families.sql`

#### Task 6.2: 家庭共享 Model
- [ ] `model/family.go` - 家庭共享数据模型
  - Family 结构体
  - FamilyMember 结构体
  - CreateFamily 方法
  - JoinFamily 方法 (通过邀请码)
  - InviteFamilyMember 方法
  - ListFamilyMembers 方法

#### Task 6.3: 家庭共享 Handler
- [ ] `handler/family.go` - 家庭共享处理器
  - CreateFamilyHandler: 创建家庭
  - GetFamilyHandler: 获取家庭信息
  - InviteMemberHandler: 邀请成员
  - JoinFamilyHandler: 加入家庭 (通过邀请码)
  - ListMembersHandler: 获取成员列表
  - RemoveMemberHandler: 移除成员

#### Task 6.4: 路由配置
- [ ] 更新 `router/router.go`
  - POST /api/family
  - GET /api/family
  - POST /api/family/invite
  - POST /api/family/join
  - GET /api/family/members
  - DELETE /api/family/members/:id

---

### 模块 7: 数据同步 (SSE)

#### Task 7.1: SSE Server
- [ ] `pkg/sse.go` - SSE服务
  - NewSSEServer 方法
  - AddClient 方法
  - RemoveClient 方法
  - BroadcastEvent 方法
  - SendToUser 方法

#### Task 7.2: SSE Event Model
- [ ] 定义事件类型
  - EVENT_RECORD_CREATED: 记录创建
  - EVENT_RECORD_UPDATED: 记录更新
  - EVENT_FAMILY_INVITE: 家庭邀请
  - EVENT_PHOTO_UPLOADED: 照片上传

#### Task 7.3: SSE Handler
- [ ] `handler/sse.go` - SSE处理器
  - SSEHandler: SSE连接端点
  - 当用户建立连接时，调用 AddClient

#### Task 7.4: 事件触发
- [ ] 在创建记录时触发 SSE 事件
- [ ] 在家庭邀请时触发 SSE 事件

#### Task 7.5: 路由配置
- [ ] 更新 `router/router.go`
  - GET /api/sse

---

### 模块 8: 数据统计

#### Task 8.1: 统计 Handler
- [ ] `handler/stats.go` - 数据统计处理器
  - GetDailyStatsHandler: 获取每日统计
    - 喂养次数、总奶量、睡眠时长、排泄次数
  - GetWeeklyStatsHandler: 获取每周统计
  - GetMonthlyStatsHandler: 获取每月统计

#### Task 8.2: 统计 Model
- [ ] `model/stats.go` - 统计数据模型
  - DailyStats 结构体
  - WeeklyStats 结构体
  - CalculateDailyStats 方法
  - CalculateWeeklyStats 方法

#### Task 8.3: 路由配置
- [ ] 更新 `router/router.go`
  - GET /api/stats/daily?date=YYYY-MM-DD
  - GET /api/stats/weekly?week=YYYY-Www
  - GET /api/stats/monthly?month=YYYY-MM

---

### 模块 9: 成长指南

#### Task 9.1: 数据库表设计
- [ ] 设计 milestones 表结构
  - id, age_months, title, description, category, created_at, updated_at
- [ ] 编写迁移脚本: `mysql/migrations/007_milestones.sql`

#### Task 9.2: 成长指南 Handler
- [ ] `handler/milestone.go` - 成长里程碑处理器
  - GetMilestonesHandler: 获取月龄对应的里程碑
  - ListAllMilestonesHandler: 获取所有里程碑

#### Task 9.3: 路由配置
- [ ] 更新 `router/router.go`
  - GET /api/milestones?age=X
  - GET /api/milestones/all

---

### 模块 10: AI问答

#### Task 10.1: AI Handler
- [ ] `handler/ai.go` - AI问答处理器
  - AskQuestionHandler: 发送问题到AI
  - VoiceToTextHandler: 语音转文字 (集成阿里云语音识别)

#### Task 10.2: 集成通义千问API
- [ ] 配置通义千问API密钥
- [ ] `pkg/qianwen.go` - 通义千问服务
  - Ask 方法
  - StreamAsk 方法 (流式响应)

#### Task 10.3: 语音识别集成
- [ ] 集成阿里云语音识别SDK
- [ ] `pkg/speech.go` - 语音识别服务
  - RecognizeSpeech 方法

#### Task 10.4: 路由配置
- [ ] 更新 `router/router.go`
  - POST /api/ai/question
  - POST /api/ai/voice

---

## 开发顺序建议

### Sprint 1 (Week 1): 基础设施
1. Task 1.1: 数据库表设计 (users, babies, family_members)
2. Task 1.4: 认证中间件
3. Task 1.5: 路由配置 (auth相关)
4. Task 1.2: 用户认证 Handler

### Sprint 2 (Week 2): 核心记录 - 喂养
5. Task 2.1: 数据库表设计 (feeding_records)
6. Task 2.2: 喂养记录 Model
7. Task 2.3: 喂养记录 Handler
8. Task 2.4: 路由配置
9. Task 2.5: 喂奶提醒

### Sprint 3 (Week 2): 核心记录 - 睡眠 & 排泄
10. Task 3.1: 数据库表设计 (sleep_records)
11. Task 3.2: 睡眠记录 Model
12. Task 3.3: 睡眠记录 Handler
13. Task 3.4: 路由配置
14. Task 4.1-4.4: 排泄记录 (全流程)

### Sprint 4 (Week 3): 宝宝管理
15. Task 1.3: 宝宝管理 Handler
16. Task 1.5: 路由配置 (baby相关)

### Sprint 5 (Week 4): 照片管理
17. Task 5.1: 数据库表设计 (photos)
18. Task 5.2: OSS集成
19. Task 5.3: 照片 Handler
20. Task 5.4: 路由配置

### Sprint 6 (Week 5): 家庭共享
21. Task 6.1: 数据库表设计 (families, family_invites)
22. Task 6.2: 家庭共享 Model
23. Task 6.3: 家庭共享 Handler
24. Task 6.4: 路由配置

### Sprint 7 (Week 6): 数据同步
25. Task 7.1: SSE Server
26. Task 7.2: SSE Event Model
27. Task 7.3: SSE Handler
28. Task 7.4: 事件触发
29. Task 7.5: 路由配置

### Sprint 8 (Week 7): 数据统计 & 成长指南
30. Task 8.1-8.3: 数据统计 (全流程)
31. Task 9.1-9.3: 成长指南 (全流程)

### Sprint 9 (Week 8): AI问答
32. Task 10.2: 集成通义千问API
33. Task 10.3: 语音识别集成
34. Task 10.1: AI Handler
35. Task 10.4: 路由配置

---

## 验收标准

### 代码质量
- [ ] 所有API接口实现完成
- [ ] 数据库迁移脚本完整
- [ ] 错误处理统一
- [ ] 日志记录完善
- [ ] API文档同步更新

### 测试
- [ ] 单元测试覆盖率 > 60%
- [ ] 集成测试通过
- [ ] API接口测试通过

### 性能
- [ ] API响应时间 < 200ms (P95)
- [ ] 数据库查询优化
- [ ] 文件上传支持断点续传

---

## 待确认事项

- [ ] 阿里云OSS配置信息
- [ ] 极光推送AppKey和MasterSecret
- [ ] 通义千问API密钥
- [ ] 阿里云语音识别配置
- [ ] 数据库连接信息 (host, port, user, password, dbname)

---

**下一步**: 使用 Codex 逐个实现上述任务
