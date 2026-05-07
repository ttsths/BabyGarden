/**
 * 园子宝宝 App - Mock 数据层
 * 用于开发和演示的静态数据
 */

// 类型定义
export interface BabyInfo {
  name: string;
  age: string;
  avatar: string;
  birthday: string;
  gender: 'male' | 'female';
}

export interface TimelineRecord {
  id: string;
  type: string;
  icon: string;
  title: string;
  description: string;
  time: string;
  color: string;
}

export interface TimelineDay {
  date: string;
  count: number;
  records: TimelineRecord[];
}

export interface PhotoData {
  id: string;
  url: string;
  time: string;
  caption: string;
}

export interface PhotoDay {
  date: string;
  count: number;
  photos: PhotoData[];
}

export interface PhotoMonth {
  month: string;
  days: PhotoDay[];
}

export interface StatTrend {
  average: number;
  unit: string;
  trend: number[];
}


// 宝宝基本信息
export const babyInfo = {
  name: '圆子',
  age: '3 个月 15 天',
  avatar: '/assets/baby-avatar.png',
  birthday: '2025-11-24',
  gender: 'female',
};

// 今日完成度数据
export const todayProgress = {
  completed: 8,
  total: 10,
  percentage: 80,
  breakdown: [
    { category: '喂奶', completed: 5, total: 6, percentage: 83 },
    { category: '睡眠', completed: 12, total: 12, percentage: 100, unit: 'h' },
    { category: '排泄', completed: 3, total: 4, percentage: 75 },
  ],
};

// 今日概览统计数据
export const todayStats = [
  {
    id: 'feeding',
    icon: '🍼',
    label: '喂奶',
    value: '5 次',
    bgColor: 'bg-orange-50',
    iconBg: 'bg-orange-100',
  },
  {
    id: 'sleep',
    icon: '💤',
    label: '睡眠',
    value: '12 小时',
    bgColor: 'bg-purple-50',
    iconBg: 'bg-purple-100',
  },
  {
    id: 'diaper',
    icon: '💩',
    label: '排泄',
    value: '3 次',
    bgColor: 'bg-green-50',
    iconBg: 'bg-green-100',
  },
  {
    id: 'temperature',
    icon: '🌡️',
    label: '体温',
    value: '36.5°C',
    bgColor: 'bg-red-50',
    iconBg: 'bg-red-100',
  },
];

// 最近记录时间轴
export const recentTimeline = [
  {
    id: '1',
    type: 'feeding',
    icon: '🍼',
    title: '喂奶',
    description: '120ml',
    time: '10:30',
    color: 'text-brand-primary',
  },
  {
    id: '2',
    type: 'diaper',
    icon: '💩',
    title: '换尿布',
    description: '大便，黄色糊状',
    time: '09:15',
    color: 'text-accent-positive',
  },
  {
    id: '3',
    type: 'wake',
    icon: '☀️',
    title: '起床',
    description: '宝宝自己醒了',
    time: '07:00',
    color: 'text-warning',
  },
];

// 完整时间轴数据（用于时间轴页面）
export const fullTimeline = [
  {
    date: '今天 · 3 月 6 日 周五',
    count: 8,
    records: [
      {
        id: '1',
        type: 'feeding',
        icon: '🍼',
        title: '喂奶',
        description: '母乳 120ml',
        time: '10:30',
        color: 'text-brand-primary',
      },
      {
        id: '2',
        type: 'diaper',
        icon: '💩',
        title: '换尿布',
        description: '大便，黄色糊状',
        time: '09:15',
        color: 'text-accent-positive',
      },
      {
        id: '3',
        type: 'sleep',
        icon: '🌙',
        title: '午睡开始',
        description: '',
        time: '12:00',
        color: 'text-accent-sleep',
      },
      {
        id: '4',
        type: 'sleep',
        icon: '☀️',
        title: '午睡结束',
        description: '睡了 2 小时',
        time: '14:00',
        color: 'text-brand-primary',
      },
      {
        id: '5',
        type: 'feeding',
        icon: '🍼',
        title: '喂奶',
        description: '奶粉 180ml',
        time: '14:30',
        color: 'text-brand-primary',
      },
      {
        id: '6',
        type: 'wake',
        icon: '☀️',
        title: '起床',
        description: '宝宝自己醒了',
        time: '07:00',
        color: 'text-warning',
      },
    ],
  },
  {
    date: '昨天 · 3 月 5 日 周四',
    count: 10,
    records: [
      {
        id: '7',
        type: 'feeding',
        icon: '🍼',
        title: '喂奶',
        description: '母乳 150ml',
        time: '22:00',
        color: 'text-brand-primary',
      },
      {
        id: '8',
        type: 'sleep',
        icon: '🌙',
        title: '夜间睡眠开始',
        description: '',
        time: '20:30',
        color: 'text-accent-sleep',
      },
    ],
  },
];

// 记录类型选项
export const recordTypes = [
  { id: 'feeding', icon: '🍼', label: '喂奶' },
  { id: 'sleep', icon: '💤', label: '睡觉' },
  { id: 'diaper', icon: '💩', label: '换尿布' },
  { id: 'temperature', icon: '🌡️', label: '体温' },
  { id: 'food', icon: '🍚', label: '辅食' },
  { id: 'medicine', icon: '💊', label: '用药' },
  { id: 'bath', icon: '🛁', label: '洗澡' },
  { id: 'other', icon: '📷', label: '其他' },
];

// 底部导航栏
export const bottomNavItems = [
  { id: 'home', icon: '🏠', label: '首页', path: '/' },
  { id: 'stats', icon: '📊', label: '统计', path: '/stats' },
  { id: 'record', icon: '➕', label: '记录', path: '/record' },
  { id: 'photos', icon: '📸', label: '照片墙', path: '/photos' },
  { id: 'settings', icon: '⚙️', label: '设置', path: '/settings' },
];

// 空状态文案
export const emptyStates = {
  noRecords: {
    icon: '🍼',
    title: '开始记录圆子的第一笔吧 ✨',
    action: '去记录',
  },
  noPhotos: {
    icon: '📸',
    title: '用照片定格美好瞬间',
    action: '拍一张',
  },
  noNetwork: {
    icon: '☁️',
    title: '网络开小差了',
    action: '重试',
  },
  searchNoResult: {
    icon: '🔍',
    title: '没找到相关内容',
    action: '清除搜索',
  },
};

// 家庭成员
export const familyMembers = [
  { id: '1', role: '妈妈', avatar: '👩', name: '妈妈' },
  { id: '2', role: '爸爸', avatar: '👨', name: '爸爸' },
  { id: '3', role: '奶奶', avatar: '👵', name: '奶奶' },
];

// 周统计数据
export const weeklyStats = {
  feeding: {
    average: 6.2,
    unit: '次/天',
    trend: [5, 7, 6, 8, 5, 6, 7],
  },
  sleep: {
    average: 13.5,
    unit: '小时/天',
    trend: [12, 14, 13, 15, 13, 14, 13],
  },
  diaper: {
    average: 4.5,
    unit: '次/天',
    trend: [3, 5, 4, 6, 4, 5, 4],
  },
  totalRecords: 42,
  completeDays: 5,
};

// 照片墙数据
export const photoWallData = [
  {
    month: '2026 年 3 月',
    days: [
      {
        date: '3 月 6 日',
        count: 5,
        photos: [
          { id: '1', url: '/photos/2026-03-06-1.jpg', time: '10:30', caption: '第一次自己坐起来！' },
          { id: '2', url: '/photos/2026-03-06-2.jpg', time: '14:20', caption: '笑得好开心' },
          { id: '3', url: '/photos/2026-03-06-3.jpg', time: '16:45', caption: '洗澡时间' },
          { id: '4', url: '/photos/2026-03-06-4.jpg', time: '18:00', caption: '辅食时间' },
          { id: '5', url: '/photos/2026-03-06-5.jpg', time: '20:30', caption: '睡前故事' },
        ],
      },
      {
        date: '3 月 5 日',
        count: 7,
        photos: [
          { id: '6', url: '/photos/2026-03-05-1.jpg', time: '09:00', caption: '晒太阳' },
          { id: '7', url: '/photos/2026-03-05-2.jpg', time: '11:30', caption: '玩耍时间' },
          { id: '8', url: '/photos/2026-03-05-3.jpg', time: '15:00', caption: '午睡醒来' },
          { id: '9', url: '/photos/2026-03-05-4.jpg', time: '17:30', caption: '和爸爸玩' },
          { id: '10', url: '/photos/2026-03-05-5.jpg', time: '19:00', caption: '洗澡' },
          { id: '11', url: '/photos/2026-03-05-6.jpg', time: '20:00', caption: '喝奶' },
          { id: '12', url: '/photos/2026-03-05-7.jpg', time: '21:00', caption: '睡着了' },
        ],
      },
    ],
  },
];

// 设置项数据
export const settingsData = {
  account: {
    name: '申屠海三',
    email: 'shentuhaisan@example.com',
    avatar: '/avatars/user-1.png',
  },
  babies: [
    { id: '1', name: '圆子', age: '3 个月 15 天', avatar: '/avatars/baby-1.png' },
  ],
  notifications: {
    feedingReminder: true,
    pushNotifications: true,
    dailyReport: true,
  },
  appearance: {
    darkMode: false,
    elderlyMode: false,
  },
};

// Toast 消息类型
export type ToastType = 'success' | 'error' | 'warning' | 'info';

// Toast 配置
export const toastConfig = {
  success: {
    icon: '✅',
    duration: 2000,
    bgColor: 'bg-accent-positive',
  },
  error: {
    icon: '❌',
    duration: 3000,
    bgColor: 'bg-error',
  },
  warning: {
    icon: '⚠️',
    duration: 3000,
    bgColor: 'bg-warning',
  },
  info: {
    icon: 'ℹ️',
    duration: 2000,
    bgColor: 'bg-neutral-text-secondary',
  },
};
