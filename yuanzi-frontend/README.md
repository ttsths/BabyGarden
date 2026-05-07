# 圆子 (Yuanzi) - 母婴记录应用前端

> 简洁而智能的育儿记录应用

## 技术栈

| 技术 | 版本 | 说明 |
|-----|------|------|
| React | 19.x | UI 框架 |
| TypeScript | 5.x | 类型系统 |
| Vite | 7.x | 构建工具 |
| Ant Design Mobile | 5.x | UI 组件库 |
| Tailwind CSS | 4.x | 原子化 CSS |
| Zustand | 5.x | 状态管理 |
| React Query | 5.x | 服务端状态管理 |
| React Router | 6.x | 路由管理 |

## 项目结构

```
yuanzi-frontend/
├── public/                 # 静态资源
├── src/
│   ├── components/         # 组件库
│   │   ├── ui/            # 基础 UI 组件
│   │   ├── business/      # 业务组件
│   │   └── features/      # 功能组件
│   ├── pages/             # 页面组件
│   ├── stores/            # 状态管理 (Zustand)
│   ├── hooks/             # 自定义 Hooks
│   ├── services/          # API 服务层
│   ├── utils/             # 工具函数
│   ├── types/             # TypeScript 类型
│   ├── constants/         # 常量定义
│   ├── styles/            # 全局样式
│   └── assets/            # 静态资源
├── .env.example           # 环境变量示例
├── package.json           # 项目配置
├── tailwind.config.js     # Tailwind 配置
├── tsconfig.json          # TypeScript 配置
└── vite.config.ts         # Vite 配置
```

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，填入实际的 API 地址
```

### 3. 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:3000

### 4. 构建生产版本

```bash
npm run build
```

## 开发规范

### 代码风格

- **ESLint**: 代码检查
- **Prettier**: 代码格式化
- **TypeScript**: 严格类型检查

### 提交规范

遵循 Conventional Commits：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type 类型：**
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具链

### 组件开发规范

1. **Props 定义**: 使用 TypeScript 接口，添加 JSDoc 注释
2. **样式隔离**: 使用 CSS Modules
3. **命名规范**: 组件使用 PascalCase，文件与组件同名
4. **导出方式**: 组件目录下创建 index.ts 统一导出

## 核心功能

### 1. 语音输入
- Web Speech API 优先
- 阿里云 NLS 降级方案
- 长按说话 + 波形动画

### 2. 照片管理
- 前端生成缩略图
- OSS 直传
- 懒加载 + 虚拟列表

### 3. 实时同步
- SSE (Server-Sent Events)
- 断线重连 + 消息去重
- 乐观更新

### 4. 祖辈模式
- CSS 变量主题切换
- 字体放大 1.5 倍
- 高对比度配色
- 双击首页快捷切换

## 后端接口

后端采用 **Go + Gin** 框架，API 遵循 RESTful 规范，Swagger 自动生成文档。

参考项目: `/Users/ttsths/go/src/zeus`

## 浏览器支持

- iOS Safari 14+
- Android Chrome 90+
- 微信内置浏览器

## 许可证

MIT
