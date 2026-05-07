# 园子 Baby App - 设计令牌

> 基于 Stitch 设计系统 (yuanzi_stitch_design) 更新  
> 主色：珊瑚粉 `#FF9A8B` | 背景：暖白 `#FFFBF7` | 辅助：薄荷绿 `#A2D5C6`

---

## 📋 目录

- [色彩系统](#色彩系统)
- [字体层级](#字体层级)
- [间距系统](#间距系统)
- [圆角规范](#圆角规范)
- [阴影层级](#阴影层级)
- [组件使用指南](#组件使用指南)

---

## 色彩系统

### 品牌色（珊瑚粉）

| 名称 | Tailwind 类 | 色值 | 用途 |
|------|------------|------|------|
| Primary | `bg-primary` | `#FF9A8B` | 主按钮、激活状态、FAB |
| Primary Light | `bg-primary-light` | `#FFB4A2` | 悬停状态 |
| Primary Dark | `bg-primary-dark` | `#E8887A` | 按压状态 |
| Primary Soft | `bg-primary-soft` | `#FFEEEB` | 浅色背景装饰 |

### 辅助色（薄荷绿）

| 名称 | Tailwind 类 | 色值 | 用途 |
|------|------------|------|------|
| Mint | `bg-mint` | `#A2D5C6` | 喂奶图标、正向反馈 |
| Mint Light | `bg-mint-light` | `#E7F3EB` | 薄荷色浅色背景 |
| Mint Text | `text-mint-text` | `#4e9767` | 薄荷色文字 |

### 背景色

| 名称 | Tailwind 类 | 色值 | 用途 |
|------|------------|------|------|
| Background Light | `bg-background-light` | `#FFFBF7` | 页面主背景（暖白） |
| Background Dark | `bg-background-dark` | `#23110F` | 暗色模式背景 |
| Background Card | `bg-background-card` | `#FFFFFF` | 卡片背景（纯白） |
| Neutral Soft | `bg-neutral-soft` | `#F4E8E6` | 柔和中性色背景 |

### 文字色（使用 Slate）

| 名称 | Tailwind 类 | 色值 | 用途 |
|------|------------|------|------|
| Primary | `text-slate-900` | `#0F172A` | 主标题、重要文字 |
| Secondary | `text-slate-600` | `#475569` | 副标题、描述文字 |
| Tertiary | `text-slate-400` | `#94A3B8` | 辅助文字、时间戳 |

### 语义色

| 名称 | Tailwind 类 | 色值 | 用途 |
|------|------------|------|------|
| Success | `bg-success` | `#6BCB77` | 成功状态 |
| Warning | `bg-warning` | `#FFD54F` | 警告状态 |
| Error | `bg-error` | `#FF6B6B` | 错误状态 |
| Info | `bg-info` | `#4FC3F7` | 信息提示 |

### 功能色（图标背景）

| 名称 | Tailwind 类 | 色值 | 用途 |
|------|------------|------|------|
| Blue 100 | `bg-blue-100` | `#DBEAFE` | 睡觉图标背景 |
| Orange 100 | `bg-orange-100` | `#FED7AA` | 换尿布图标背景 |

---

## 字体层级

**字体家族:** `Plus Jakarta Sans`, `Noto Sans SC`, system-ui

| 级别 | Tailwind 类 | 字号 | 字重 | 行高 | 用途 |
|-----|------------|------|------|------|------|
| Display | `text-2xl` | 24px | 600 | 1.3 | 大标题、数据展示 |
| H1 | `text-xl` | 20px | 600 | 1.4 | 页面标题 |
| H2 | `text-lg` | 18px | 600 | 1.4 | 区块标题 |
| H3 | `text-base` | 16px | 500 | 1.5 | 卡片标题 |
| Body | `text-sm` | 14px | 400 | 1.6 | 正文内容 |
| Caption | `text-xs` | 12px | 400 | 1.5 | 辅助文字 |
| Small | `text-[10px]` | 10px | 400/500 | 1.4 | 导航标签 |

---

## 间距系统

基准单位：4px

| 名称 | Tailwind 类 | 数值 | 用途 |
|------|------------|------|------|
| XS | `p-1`, `m-1` | 4px | 最小间距 |
| SM | `p-2`, `m-2` | 8px | 小组件间距 |
| MD | `p-3`, `m-3` | 12px | 卡片内边距 |
| LG | `p-4`, `m-4` | 16px | 页面边距 |
| XL | `p-5`, `m-5` | 20px | 大区块间距 |
| 2XL | `p-6`, `m-6` | 24px | 页面级间距 |
| 3XL | `p-8`, `m-8` | 32px | 特殊大间距 |

---

## 圆角规范

| 名称 | Tailwind 类 | 数值 | 用途 |
|------|------------|------|------|
| Default | `rounded` | 8px (0.5rem) | 默认圆角 |
| SM | `rounded-sm` | 6px (0.375rem) | 小圆角 |
| MD | `rounded-md` | 12px (0.75rem) | 中等圆角 |
| LG | `rounded-lg` | 16px (1rem) | 大圆角 |
| XL | `rounded-xl` | 24px (1.5rem) | 超大圆角 |
| Full | `rounded-full` | 9999px | 圆形按钮、头像 |

---

## 阴影层级

| 名称 | Tailwind 类 | 效果 | 用途 |
|------|------------|------|------|
| SM | `shadow-sm` | 0 1px 2px | 轻微浮起 |
| Card | `shadow-card` | 0 2px 8px | 卡片默认 |
| MD | `shadow-md` | 0 4px 12px | 中等浮起 |
| LG | `shadow-lg` | 0 8px 24px | 明显浮起 |
| XL | `shadow-xl` | 0 12px 32px | 强烈浮起 |
| Primary | `shadow-primary` | 0 4px 14px (珊瑚粉) | FAB、主按钮 |

---

## 组件使用指南

### 快速记录按钮（FAB 风格）

```tsx
<button className="size-16 rounded-full bg-primary text-white shadow-lg shadow-primary/30 active:scale-95 transition-transform">
  <span className="material-symbols-outlined text-3xl">restaurant</span>
</button>
```

### 统计卡片

```tsx
<div className="min-w-[140px] flex-1 bg-white rounded-lg p-4 shadow-sm border border-neutral-soft">
  <div className="size-10 bg-mint/20 rounded-full flex items-center justify-center mb-3">
    <span className="material-symbols-outlined text-mint">baby_changing_station</span>
  </div>
  <p className="text-slate-500 text-sm">喂奶</p>
  <p className="text-xl font-bold">5 次</p>
</div>
```

### 活动列表项

```tsx
<div className="flex items-center gap-4 p-3 bg-background-light rounded-lg">
  <div className="size-10 rounded-full bg-mint/20 flex items-center justify-center text-mint">
    <span className="material-symbols-outlined text-xl">restaurant</span>
  </div>
  <div className="flex-1">
    <p className="font-bold text-slate-800">喂奶 (120ml)</p>
    <p className="text-xs text-slate-500">母乳喂养 • 15 分钟</p>
  </div>
  <span className="text-sm text-slate-400">10:30 AM</span>
</div>
```

### 底部导航

```tsx
<nav className="flex gap-2 border-t border-neutral-soft bg-white px-4 pb-6 pt-2">
  <a className="flex flex-1 flex-col items-center justify-center gap-1 text-primary">
    <span className="material-symbols-outlined fill-icon text-2xl">home</span>
    <p className="text-xs font-semibold">首页</p>
  </a>
  <a className="flex flex-1 flex-col items-center justify-center gap-1 text-slate-400">
    <span className="material-symbols-outlined text-2xl">add_circle</span>
    <p className="text-xs font-medium">记录</p>
  </a>
  {/* ... */}
</nav>
```

### 悬浮操作按钮（FAB）

```tsx
<button className="fixed bottom-24 right-6 flex h-16 w-16 items-center justify-center rounded-full bg-primary text-white shadow-lg shadow-primary/30 hover:scale-105 active:scale-95 transition-all z-20">
  <span className="material-symbols-outlined text-3xl">add</span>
</button>
```

---

## 暗色模式

在根元素添加 `dark` 类：

```html
<html class="dark">
```

暗色模式下：
- 背景：`#23110F`（深棕色）
- 卡片：`bg-slate-800`
- 文字：`text-slate-100`
- 边框：`border-slate-700`

---

## Material Icons 使用

```html
<!-- 空心图标 -->
<span className="material-symbols-outlined">home</span>

<!-- 实心图标 -->
<span className="material-symbols-outlined fill-1">home</span>

<!-- 自定义大小 -->
<span className="material-symbols-outlined text-3xl">add</span>
```

常用图标：
- `home` - 首页
- `add_circle` - 添加记录
- `bar_chart` - 统计
- `settings` - 设置
- `restaurant` - 喂奶
- `bedtime` - 睡觉
- `baby_changing_station` - 换尿布
- `layers` - 尿布
- `insights` - 统计洞察

---

## 更新日志

- **v2.0.0** (2026-03-12): 根据 Stitch 设计系统全面更新 - 主色改为珊瑚粉 `#FF9A8B`，背景改为暖白 `#FFFBF7`，添加薄荷绿辅助色，圆角增大
- **v1.0.0** (2026-03-11): 初始版本
