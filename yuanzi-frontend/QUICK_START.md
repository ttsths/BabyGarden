# 🚀 Yuanzi Baby App - 快速启动指南

## 环境要求

- Node.js >= 22.12.0
- npm >= 10.9.0

## 安装依赖

```bash
npm install
```

## 启动开发服务器

```bash
npm run dev
```

访问 http://localhost:3000/

## 构建生产版本

```bash
npm run build
```

## 预览生产构建

```bash
npm run preview
```

## 代码检查

```bash
# TypeScript 类型检查
npx tsc --noEmit

# ESLint 代码检查
npm run lint
```

## 项目结构

```
yuanzi-frontend/
├── src/
│   ├── components/
│   │   ├── ui/           # 基础 UI 组件（15 个）
│   │   ├── features/     # 功能组件
│   │   └── business/     # 业务组件
│   ├── data/             # Mock 数据
│   ├── styles/           # 样式文件
│   ├── pages/            # 页面组件
│   ├── hooks/            # 自定义 Hooks
│   ├── stores/           # 状态管理（Zustand）
│   ├── services/         # API 服务
│   ├── utils/            # 工具函数
│   ├── constants/        # 常量定义
│   ├── types/            # TypeScript 类型
│   └── assets/           # 静态资源
├── DESIGN_TOKENS.md      # 设计令牌文档
├── TASK_COMPLETE.md      # 任务完成报告
└── package.json
```

## 可用组件

### 基础组件（15 个）

1. **Button** - 按钮（5 种变体，4 种尺寸）
2. **Card** - 卡片（4 种 padding）
3. **Input** - 输入框（3 种尺寸）
4. **ProgressRing** - 环形进度条
5. **StatsCard** - 统计卡片
6. **QuickStatCard** - 快捷统计卡片
7. **TimelineItem** - 时间轴项目
8. **BottomNav** - 底部导航
9. **Toast** - 提示组件
10. **Skeleton** - 骨架屏
11. **Avatar** - 头像
12. **Badge** - 标签
13. **Modal** - 弹窗
14. **EmptyState** - 空状态
15. **Loading** - 加载
16. **Switch** - 开关

### 使用示例

```tsx
import { Button, Card, Input, Toast } from '@/components/ui';

function MyComponent() {
  return (
    <Card padding="medium">
      <Input label="用户名" placeholder="请输入" />
      <Button variant="primary" onClick={handleSubmit}>
        提交
      </Button>
    </Card>
  );
}
```

## 设计令牌

所有设计令牌已定义在 `src/styles/variables.css` 和 `tailwind.config.js` 中。

### 色彩系统

```tsx
// 品牌色
bg-brand-primary    // #FF9A8B
bg-brand-light      // #FFB4A2
bg-brand-dark       // #E8887A

// 背景色
bg-background-primary   // #FFFBF7
bg-background-card      // #FFFFFF
bg-background-secondary // #FFF5ED

// 文字色
text-neutral-text-primary   // #2D2D2D
text-neutral-text-secondary // #6B7280
text-neutral-text-tertiary  // #9CA3AF

// 辅助色
bg-accent-positive  // #A2D5C6
bg-accent-sleep     // #D8BFD8
```

### 字体层级

```tsx
text-display   // 28px/600/1.3
text-h1        // 22px/600/1.4
text-h2        // 18px/500/1.4
text-h3        // 16px/500/1.5
text-body      // 15px/400/1.6
text-caption   // 13px/400/1.5
text-small     // 12px/400/1.4
```

### 间距系统

```tsx
p-xs    // 4px
p-sm    // 8px
p-md    // 12px
p-lg    // 16px
p-xl    // 24px
p-2xl   // 32px
```

## Mock 数据

所有 Mock 数据在 `src/data/mockData.ts` 中：

```tsx
import { 
  babyInfo, 
  todayProgress, 
  todayStats, 
  recentTimeline,
  fullTimeline,
  weeklyStats,
  photoWallData,
  settingsData,
} from '@/data/mockData';
```

## 特殊模式

### 暗色模式

在根元素添加 `dark` 类：

```html
<html class="dark">
```

### 祖辈模式

在根元素添加 `elderly-mode` 类：

```html
<html class="elderly-mode">
```

## 开发规范

1. **使用设计令牌** - 不要使用硬编码的颜色值、间距等
2. **组件化** - 可复用的 UI 应该放在 `components/ui/`
3. **类型安全** - 所有组件必须有 TypeScript 类型定义
4. **响应式** - 所有组件必须支持移动端优先
5. **无障碍** - 遵循 WCAG 2.2 标准

## 常见问题

### Q: 如何添加新组件？

1. 在 `src/components/ui/` 创建组件文件
2. 导出组件和类型
3. 在 `index.ts` 中添加导出
4. 更新 `DESIGN_TOKENS.md` 文档

### Q: 如何修改设计令牌？

1. 修改 `src/styles/variables.css` 中的 CSS 变量
2. 更新 `tailwind.config.js` 中的映射
3. 更新 `DESIGN_TOKENS.md` 文档

### Q: 如何测试组件？

```bash
npm run test
```

## 相关文档

- [DESIGN_TOKENS.md](./DESIGN_TOKENS.md) - 完整设计令牌文档
- [TASK_COMPLETE.md](./TASK_COMPLETE.md) - 任务完成报告
- [README.md](./README.md) - 项目说明

## 需要帮助？

查看项目文档或联系开发团队。

---

*最后更新：2026-03-11*
