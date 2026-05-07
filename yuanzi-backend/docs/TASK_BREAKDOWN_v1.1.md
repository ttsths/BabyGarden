# Yuanzi 后端开发详细任务拆解

> **版本**: v1.1 (基于 PRD-V1.1)
> **日期**: 2026-03-09
> **技术栈**: Go 1.21 + Gin + MySQL 8.0 + Redis 7 + GORM + JWT
> **项目目录**: `/Users/ttsths/Desktop/AGI/AGI/01-Projects/Projects/Active/2026-Q1/Yuanzi/04-Code/yuanzi-backend`

---

## 📊 任务总览

| 阶段 | 模块 | 任务数 | 优先级 | 状态 |
|------|------|--------|--------|------|
| Phase 1 | 基础设施 | 5 | P0 | ✅ 已完成 (TASK-001/002 完成) |
| Phase 2 | 用户认证 | 4 | P0 | ⏳ 待开发 |
| Phase 3 | 家庭管理 | 5 | P0 | ⏳ 待开发 |
| Phase 4 | 宝宝管理 | 5 | P0 | ⏳ 待开发 |
| Phase 5 | 记录管理 | 8 | P0 | ⏳ 待开发 |
| Phase 6 | 照片管理 | 4 | P1 | ⏳ 待开发 |
| Phase 7 | AI 集成 | 3 | P1 | ⏳ 待开发 |
| Phase 8 | 同步推送 | 2 | P1 | ⏳ 待开发 |
| Phase 9 | 统计分析 | 2 | P1 | ⏳ 待开发 |

---

## 🏗️ Phase 1: 基础设施 (P0)

### TASK-001: 数据库连接与 GORM 配置 ✅
**状态**: ✅ 已完成 | **文件**: `mysql/mysql.go` | **测试**: `mysql/mysql_test.go`

**完成内容**:
- [x] 完成 MySQL 8.0 连接配置
- [x] 配置 GORM ORM 框架
- [x] 实现数据库连接池管理
- [x] 编写 TDD 单元测试

**配置**:
```yaml
# config/config.yaml
database:
  max_idle_conn: 20  # 已从 10 调整为 20
  max_open_conn: 100
```

**测试结果**: `go test ./mysql -v` ✅ PASS

**相关文件**:
- `mysql/mysql.go` - 数据库连接
- `mysql/mysql_test.go` - 单元测试

**Self Code Review**:
- ✅ 分层合规：DAO 层实现
- ✅ 并发安全：GORM 内建支持
- ✅ 配置合理：连接池参数符合要求

---

### TASK-002: Redis 连接与封装 ✅
**状态**: ✅ 已完成 | **文件**: `pkg/gredis/redis.go` | **测试**: `pkg/gredis/redis_test.go`

**完成内容**:
- [x] 完成 Redis 7 连接配置
- [x] 封装 Redis 常用操作（Get/Set/Del/Expire）
- [x] 封装 Hash 操作（HSet/HGet/HGetAll/HDel）
- [x] 封装 List 操作（LPush/LPop/LRange/LLen）
- [x] 封装验证码存储方法（SetCode/GetCode/DelCode）
- [x] 封装 Token 黑名单方法（AddTokenToBlacklist/IsTokenBlacklisted/RemoveTokenFromBlacklist）
- [x] 编写 TDD 单元测试（11 个用例，全部通过）

**测试结果**: `go test ./pkg/gredis -v` ✅ PASS

**相关文件**:
- `pkg/gredis/redis.go` - Redis 连接封装
- `pkg/gredis/redis_test.go` - TDD 单元测试

---

### TASK-003: JWT 认证中间件 ⏳
**状态**: ⏳ 待开发 | **文件**: `middleware/jwt.go`

**已完成**:
- ✅ JWT Token 生成与验证
- ✅ JWT 中间件

**待完善**:
- ⏳ Refresh Token 刷新逻辑需完善
- ⏳ Token黑名单机制需补充测试
- ⏳ 需补充单元测试

**相关文件**:
- `middleware/jwt.go` - JWT 中间件实现

---

### TASK-004: 统一响应格式与错误处理 ⏳
**状态**: ⏳ 待开发 | **文件**: `model/base.go`

**已完成**:
- ✅ 统一响应格式 `{code, msg, data}`
- ✅ 错误码常量定义
- ✅ 分页结构定义

**待完善**:
- ⏳ 需补充单元测试
- ⏳ 需补充自定义业务错误码

**相关文件**:
- `model/base.go` - 统一响应格式

---

## 🔐 Phase 2: 用户认证 (P0)

### TASK-005: 数据库表设计
**状态**: ⏳ 待开发 | **优先级**: P0

**表结构**:
```sql
-- users 表
CREATE TABLE users (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    phone           VARCHAR(11) UNIQUE NOT NULL COMMENT '手机号',
    nickname        VARCHAR(50) COMMENT '昵称',
    avatar_url      VARCHAR(500) COMMENT '头像 URL',
    status          TINYINT DEFAULT 1 COMMENT '状态: 0-禁用 1-正常',
    last_login_at   DATETIME COMMENT '最后登录时间',
    last_login_ip   VARCHAR(45) COMMENT '最后登录IP',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      DATETIME COMMENT '软删除',
    INDEX idx_phone (phone),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- verification_codes 表
CREATE TABLE verification_codes (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    phone           VARCHAR(11) NOT NULL COMMENT '手机号',
    code            VARCHAR(6) NOT NULL COMMENT '验证码',
    type            VARCHAR(20) NOT NULL COMMENT '类型: login/reset/bind',
    expires_at      DATETIME NOT NULL COMMENT '过期时间',
    used_at         DATETIME COMMENT '使用时间',
    ip_address      VARCHAR(45) COMMENT 'IP地址',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_phone_type (phone, type),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='验证码表';
```

**交付文件**:
- `mysql/migrations/001_users_and_verification_codes.sql`

---

### TASK-005: 验证码发送（阿里云短信）⏳
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-002, TASK-005

**功能描述**:
- 集成阿里云短信服务
- 实现验证码生成逻辑（6位随机数字）
- 实现发送限流（1分钟1次，1小时5次）
- Redis 存储验证码（5分钟过期）

**API**:
```
POST /api/v1/auth/send-code
Request: { "phone": "13800138000", "type": "login" }
Response: { "code": 200, "data": { "expires_in": 300 }, "msg": "ok" }
```

**交付文件**:
- `pkg/sms/aliyun.go` - 阿里云短信服务封装
- `handler/auth.go` - `SendVerificationCode` 函数实现
- `test/sms_test.go` - 单元测试

---

### TASK-006: 验证码登录
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-005

**功能描述**:
- 验证码校验逻辑
- 用户自动注册（如果不存在）
- 生成 Access Token + Refresh Token
- 更新登录信息（last_login_at, last_login_ip）

**API**:
```
POST /api/v1/auth/login
Request: { "phone": "13800138000", "code": "123456" }
Response: {
  "code": 200,
  "data": {
    "user": { "id", "phone", "nickname", "avatar_url" },
    "access_token": "...",
    "refresh_token": "...",
    "expires_in": 7200,
    "is_new_user": true
  },
  "msg": "ok"
}
```

**交付文件**:
- `handler/auth.go` - `Login` 函数实现
- `model/user.go` - 用户相关数据库操作
- `test/auth_test.go` - 集成测试

---

### TASK-007: Token 刷新
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-006

**功能描述**:
- Refresh Token 验证
- 生成新的 Access Token
- 实现 Token 黑名单（旧 Refresh Token 失效）

**API**:
```
POST /api/v1/auth/refresh
Request: { "refresh_token": "..." }
Response: { "code": 200, "data": { "access_token": "...", "refresh_token": "...", "expires_in": 7200 }, "msg": "ok" }
```

**交付文件**:
- `handler/auth.go` - `RefreshToken` 函数实现
- `test/auth_test.go` - Token 刷新测试

---

### TASK-008: 获取/更新用户信息
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-006

**API**:
```
GET /api/v1/user/profile
Response: { "code": 200, "data": { "id", "phone", "nickname", "avatar_url", "elder_mode" }, "msg": "ok" }

PUT /api/v1/user/profile
Request: { "nickname": "新昵称", "avatar_url": "...", "elder_mode": false }
Response: { "code": 200, "data": { "id", "nickname", "avatar_url", "elder_mode" }, "msg": "ok" }
```

**交付文件**:
- `handler/user.go` - `GetUserProfile`, `UpdateUserProfile` 函数实现
- `test/user_test.go` - 用户信息测试

---

## 👨‍👩‍👧 Phase 3: 家庭管理 (P0)

### TASK-009: 数据库表设计
**状态**: ⏳ 待开发 | **优先级**: P0

**表结构**:
```sql
-- families 表
CREATE TABLE families (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(100) NOT NULL COMMENT '家庭名称',
    invite_code     VARCHAR(8) UNIQUE NOT NULL COMMENT '8位邀请码',
    created_by      BIGINT UNSIGNED NOT NULL COMMENT '创建者用户ID',
    is_paid         TINYINT DEFAULT 0 COMMENT '是否付费: 0-否 1-是',
    storage_limit   BIGINT DEFAULT 1073741824 COMMENT '存储限制(字节), 默认1GB',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_invite_code (invite_code),
    INDEX idx_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭表';

-- family_members 表
CREATE TABLE family_members (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    family_id       BIGINT UNSIGNED NOT NULL COMMENT '家庭ID',
    user_id         BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    role            VARCHAR(20) NOT NULL DEFAULT 'member' COMMENT '角色: admin/member/elder',
    elder_mode      TINYINT DEFAULT 0 COMMENT '祖辈模式: 0-否 1-是',
    notifications   JSON COMMENT '通知设置',
    joined_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_family_user (family_id, user_id),
    INDEX idx_family_id (family_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role (family_id, role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭成员表';
```

**交付文件**:
- `mysql/migrations/002_families_and_members.sql`
- `model/family.go` - 家庭模型定义
- `test/family_test.go` - 集成测试

---

### TASK-010: 创建家庭
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-009

**功能描述**:
- 创建家庭信息
- 生成 8 位邀请码（唯一性保证）
- 创建者自动成为 admin 角色
- 创建家庭成员关联

**API**:
```
POST /api/v1/family
Request: { "name": "我的家庭" }
Response: { "code": 200, "data": { "id", "name", "invite_code", "is_paid", "storage_limit" }, "msg": "ok" }
```

**交付文件**:
- `model/family.go` - 家庭模型
- `handler/family.go` - `CreateFamily` 函数
- `test/family_test.go` - 家庭创建测试

---

### TASK-011: 获取家庭信息
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-010

**API**:
```
GET /api/v1/family/:id
Response: {
  "code": 200,
  "data": {
    "id", "name", "invite_code", "is_paid", "storage_limit",
    "members": [{ "user_id", "role", "nickname", "avatar_url" }],
    "babies": [{ "id", "name", "birthday", "gender" }]
  },
  "msg": "ok"
}
```

**交付文件**:
- `handler/family.go` - `GetFamilyDetail` 函数

---

### TASK-012: 邀请/加入家庭成员
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-010

**功能描述**:
- 使用邀请码邀请成员
- 创建家庭成员关联
- 默认角色为 member

**API**:
```
POST /api/v1/family/:id/invite
Request: { "phone": "13800138000" }
Response: { "code": 200, "data": { "success": true }, "msg": "ok" }

POST /api/v1/family/join
Request: { "invite_code": "ABC12345" }
Response: { "code": 200, "data": { "success": true, "family_id": "..." }, "msg": "ok" }
```

**交付文件**:
- `handler/family.go` - `InviteMember`, `JoinFamily` 函数

---

### TASK-013: 家庭成员管理
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-012

**API**:
```
GET /api/v1/family/:id/members
Response: { "code": 200, "data": { "members": [...] }, "msg": "ok" }

DELETE /api/v1/family/:id/members/:userId
Response: { "code": 200, "data": { "success": true }, "msg": "ok" }
```

**交付文件**:
- `handler/family.go` - `GetMembers`, `RemoveMember` 函数

---

## 👶 Phase 4: 宝宝管理 (P0)

### TASK-014: 数据库表设计
**状态**: ⏳ 待开发 | **优先级**: P0

**表结构**:
```sql
-- babies 表
CREATE TABLE babies (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    family_id       BIGINT UNSIGNED NOT NULL COMMENT '家庭ID',
    name            VARCHAR(50) NOT NULL COMMENT '宝宝姓名',
    birthday        DATE NOT NULL COMMENT '出生日期',
    gender          TINYINT NOT NULL COMMENT '性别: 1-男 2-女',
    birth_weight    DECIMAL(5,2) COMMENT '出生体重(kg)',
    birth_height    DECIMAL(4,1) COMMENT '出生身高(cm)',
    avatar_url      VARCHAR(500) COMMENT '头像 URL',
    note            TEXT COMMENT '备注',
    is_premature    TINYINT DEFAULT 0 COMMENT '是否早产: 0-否 1-是',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      DATETIME COMMENT '软删除',
    INDEX idx_family_id (family_id),
    INDEX idx_birthday (birthday)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='宝宝表';
```

**交付文件**:
- `mysql/migrations/003_babies.sql`
- `model/baby.go` - 宝宝模型定义

---

### TASK-015: 宝宝 CRUD
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-014

**API**:
```
POST /api/v1/baby
Request: { "family_id": "xxx", "name": "小宝", "birthday": "2024-01-01", "gender": 1 }
Response: { "code": 200, "data": { "id", "name", "birthday", "gender", "family_id" }, "msg": "ok" }

GET /api/v1/baby?family_id=xxx
Response: { "code": 200, "data": { "list": [...] }, "msg": "ok" }

GET /api/v1/baby/:id
Response: { "code": 200, "data": { "id", "name", "birthday", "gender", "birth_weight", "birth_height", ... }, "msg": "ok" }

PUT /api/v1/baby/:id
Request: { "name": "新名字", "avatar_url": "..." }
Response: { "code": 200, "data": { "id", "name", "avatar_url" }, "msg": "ok" }

DELETE /api/v1/baby/:id
Response: { "code": 200, "data": { "success": true }, "msg": "ok" }
```

**交付文件**:
- `handler/baby.go` - 宝宝处理器
- `model/baby.go` - 宝宝模型

---

## 📝 Phase 5: 记录管理 (P0 - MVP核心)

### TASK-017: 数据库表设计
**状态**: ⏳ 待开发 | **优先级**: P0

**表结构**:
```sql
-- records 表 (核心表)
CREATE TABLE records (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    baby_id         BIGINT UNSIGNED NOT NULL COMMENT '宝宝ID',
    family_id       BIGINT UNSIGNED NOT NULL COMMENT '家庭ID(冗余)',
    type            VARCHAR(20) NOT NULL COMMENT '类型: feeding/sleep/diaper/growth',
    started_at      DATETIME NOT NULL COMMENT '开始时间',
    ended_at        DATETIME COMMENT '结束时间',
    content         JSON NOT NULL COMMENT '类型特定内容',
    note            TEXT COMMENT '备注',
    created_by      BIGINT UNSIGNED NOT NULL COMMENT '创建者用户ID',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      DATETIME COMMENT '软删除',
    INDEX idx_baby_id (baby_id),
    INDEX idx_family_id (family_id),
    INDEX idx_type (type),
    INDEX idx_started_at (started_at),
    INDEX idx_baby_started (baby_id, started_at DESC),
    INDEX idx_not_deleted (baby_id, started_at DESC, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='记录表';
```

**content JSON 结构**:
```json
// 喂养记录 (feeding)
{ "type": "breast", "side": "left", "duration": 15, "amount": 120, "unit": "ml" }

// 睡眠记录 (sleep)
{ "quality": "good", "location": "crib" }

// 排泄记录 (diaper)
{ "type": "wet", "color": "yellow", "consistency": "normal" }

// 成长记录 (growth)
{ "weight": 8.5, "height": 70.0, "head_circumference": 42.0 }
```

**交付文件**:
- `mysql/migrations/004_records.sql`

---

### TASK-018: 喂养记录
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-017

**API**:
```
POST /api/v1/record
Request: {
  "baby_id": "xxx",
  "type": "feeding",
  "started_at": "2024-03-09T10:00:00Z",
  "content": { "type": "breast", "side": "left", "duration": 15 }
}
Response: { "id", "type", "started_at", "hours_since_last_feed" }
```

**PRD 对应**: US-001 喂养记录
- ✅ 可快速记录喂奶时间
- ✅ 显示距离上次喂奶的间隔

---

### TASK-019: 睡眠记录
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-017

**API**:
```
POST /api/v1/record
Request: {
  "baby_id": "xxx",
  "type": "sleep",
  "started_at": "2024-03-09T20:00:00Z",
  "ended_at": "2024-03-09T22:00:00Z",
  "content": { "quality": "good", "location": "crib" }
}
Response: { "id", "type", "started_at", "ended_at", "duration_hours": 2 }
```

**PRD 对应**: US-002 睡眠追踪
- ✅ 一键记录入睡/起床
- ✅ 自动计算睡眠时长

---

### TASK-020: 排泄记录
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-017

**API**:
```
POST /api/v1/record
Request: {
  "baby_id": "xxx",
  "type": "diaper",
  "started_at": "2024-03-09T14:00:00Z",
  "content": { "type": "wet", "color": "yellow", "consistency": "normal" }
}
```

---

### TASK-021: 记录查询与更新
**状态**: ⏳ 待开发 | **优先级**: P0 | **依赖**: TASK-018

**API**:
```
GET /api/v1/record?baby_id=xxx&type=feeding&date=2024-03-09&page=1&page_size=20
Response: {
  "list": [...],
  "pagination": { "page": 1, "page_size": 20, "total": 10, "total_pages": 1 }
}

GET /api/v1/record/:id
PUT /api/v1/record/:id
DELETE /api/v1/record/:id
```

---

## 📸 Phase 6: 照片管理 (P1)

### TASK-022: 数据库表设计
**状态**: ⏳ 待开发 | **优先级**: P1

**表结构**:
```sql
-- photos 表
CREATE TABLE photos (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    baby_id         BIGINT UNSIGNED NOT NULL COMMENT '宝宝ID',
    family_id       BIGINT UNSIGNED NOT NULL COMMENT '家庭ID',
    oss_key         VARCHAR(500) NOT NULL COMMENT 'OSS对象key',
    thumbnail_key   VARCHAR(500) COMMENT '缩略图key',
    width           INT COMMENT '宽度',
    height          INT COMMENT '高度',
    size            BIGINT NOT NULL COMMENT '文件大小(字节)',
    content_type    VARCHAR(50) DEFAULT 'image/jpeg',
    taken_at        DATETIME COMMENT '照片拍摄时间',
    description     TEXT COMMENT '描述',
    uploaded_by     BIGINT UNSIGNED NOT NULL,
    uploaded_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    status          VARCHAR(20) DEFAULT 'active' COMMENT '状态: pending/active/deleted',
    INDEX idx_baby_id (baby_id),
    INDEX idx_family_id (family_id),
    INDEX idx_taken_at (taken_at DESC),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='照片表';
```

---

### TASK-023: OSS 直传签名
**状态**: ⏳ 待开发 | **优先级**: P1

**API**:
```
POST /api/v1/photo/upload-url
Request: { "baby_id": "xxx", "filename": "photo.jpg", "content_type": "image/jpeg", "size": 1024000 }
Response: {
  "photo_id": "xxx",
  "upload_url": "https://oss-cn-hangzhou.aliyuncs.com/...",
  "form_data": { "OSSAccessKeyId", "signature", "policy" },
  "expires_at": 1709827200000
}
```

---

### TASK-024: 照片管理
**状态**: ⏳ 待开发 | **优先级**: P1

**API**:
```
POST /api/v1/photo/callback  - OSS上传回调
GET /api/v1/photo?baby_id=xxx&page=1&page_size=20
DELETE /api/v1/photo/:id
```

---

## 🤖 Phase 7: AI 集成 (P1)

### TASK-025: AI 问答
**状态**: ⏳ 待开发 | **优先级**: P1

**API**:
```
POST /api/v1/ai/chat
Request: { "question": "宝宝三个月了，怎么训练抬头？", "baby_id": "xxx" }
Response: { "answer": "...", "tokens_used": 150 }
```

**集成**: 通义千问 API

---

### TASK-026: 语音识别
**状态**: ⏳ 待开发 | **优先级**: P1

**API**:
```
POST /api/v1/ai/speech/recognize
Content-Type: multipart/form-data
Request: { "audio": <binary> }
Response: { "text": "宝宝喝奶了", "confidence": 0.95, "remaining_quota": 19 }
```

**集成**: 阿里云 NLS

---

### TASK-027: AI 配额查询
**状态**: ⏳ 待开发 | **优先级**: P1

**API**:
```
GET /api/v1/ai/quota
Response: {
  "speech": { "used": 5, "limit": 20, "remaining": 15 },
  "chat": { "used": 3, "limit": 10, "remaining": 7 }
}
```

---

## 🔄 Phase 8: 同步推送 (P1)

### TASK-028: SSE 实时同步
**状态**: ⏳ 待开发 | **优先级**: P1

**功能**:
- 实现 SSE 服务端推送
- Redis Pub/Sub 跨实例广播
- 推送新记录事件
- 30s 心跳保活

**API**:
```
GET /api/v1/sync/stream?family_id=xxx
Response: Server-Sent Events stream
```

---

### TASK-029: 极光推送集成
**状态**: ⏳ 待开发 | **优先级**: P1

**功能**:
- 集成极光推送 SDK
- 注册设备 Token
- 推送喂奶提醒（4 小时）

**API**:
```
POST /api/v1/device/register
Request: { "platform": "ios", "device_token": "xxx" }
```

---

## 📊 Phase 9: 数据统计 (P1)

### TASK-030: 日统计
**状态**: ⏳ 待开发 | **优先级**: P1

**API**:
```
GET /api/v1/stats/daily?baby_id=xxx&date=2024-03-09
Response: {
  "feeding": { "count": 8, "total_amount": 600 },
  "sleep": { "count": 3, "total_hours": 12 },
  "diaper": { "count": 5 }
}
```

---

### TASK-031: 周/月统计
**状态**: ⏳ 待开发 | **优先级**: P1

**API**:
```
GET /api/v1/stats/weekly?baby_id=xxx
GET /api/v1/stats/monthly?baby_id=xxx&month=2024-03
```

---

## 📋 开发顺序建议

### Week 1: 基础设施 + 用户认证
1. TASK-005: 数据库表设计 (users, verification_codes)
2. TASK-006: 验证码发送
3. TASK-007: 验证码登录
4. TASK-008: Token 刷新
5. TASK-009: 用户信息

### Week 2: 家庭 + 宝宝管理
6. TASK-010: 数据库表设计 (families, family_members)
7. TASK-011 ~ TASK-014: 家庭管理
8. TASK-015: 数据库表设计 (babies)
9. TASK-016: 宝宝 CRUD

### Week 3: 核心记录 (喂养/睡眠/排泄)
10. TASK-017: 数据库表设计 (records)
11. TASK-018: 喂养记录
12. TASK-019: 睡眠记录
13. TASK-020: 排泄记录
14. TASK-021: 记录查询与更新

### Week 4: 照片 + AI
15. TASK-022 ~ TASK-024: 照片管理
16. TASK-025 ~ TASK-027: AI 集成

### Week 5: 同步推送 + 统计
17. TASK-028: SSE 实时同步
18. TASK-029: 极光推送
19. TASK-030 ~ TASK-031: 数据统计

---

## ✅ 验收标准

### 代码质量
- [ ] 所有 API 接口实现完成
- [ ] 数据库迁移脚本完整
- [ ] 单元测试覆盖率 > 60%
- [ ] 错误处理统一
- [ ] 日志记录完善

### 功能验收
- [ ] 用户可注册/登录
- [ ] 可创建家庭并邀请成员
- [ ] 可创建宝宝档案
- [ ] 可记录喂养/睡眠/排泄
- [ ] 可查看统计图表
- [ ] 数据实时同步

---

**文档更新时间**: 2026-03-09
**下一步**: 使用 codex exec 逐个实现上述任务
