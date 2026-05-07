// LocalStorage Key
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'yuanzi_auth_token',
  AUTH_REFRESH: 'yuanzi_auth_refresh',
  THEME: 'yuanzi_theme',
  BABY_ID: 'yuanzi_current_baby',
  SETTINGS: 'yuanzi_settings',
} as const;

// 路由
export const ROUTES = {
  // 公开路由
  LOGIN: '/login',
  REGISTER: '/register',
  
  // 认证路由
  HOME: '/',
  RECORD: '/record',
  RECORD_DETAIL: '/record/:id',
  TIMELINE: '/timeline',
  PHOTOS: '/photos',
  PHOTO_DETAIL: '/photos/:id',
  STATS: '/stats',
  AI_CHAT: '/ai',
  SETTINGS: '/settings',
  BABY_SETUP: '/baby/setup',
  
  // 祖辈模式路由
  ELDERLY_HOME: '/elderly',
  ELDERLY_RECORD: '/elderly/record',
  ELDERLY_PHOTOS: '/elderly/photos',
} as const;

// 记录类型
export const RECORD_TYPES = {
  FEEDING: 'feeding',
  SLEEP: 'sleep',
  DIAPER: 'diaper',
  TEMPERATURE: 'temperature',
  FOOD: 'food',
  MEDICINE: 'medicine',
} as const;

// 图标映射
export const RECORD_TYPE_ICONS = {
  [RECORD_TYPES.FEEDING]: '🍼',
  [RECORD_TYPES.SLEEP]: '💤',
  [RECORD_TYPES.DIAPER]: '💩',
  [RECORD_TYPES.TEMPERATURE]: '🌡️',
  [RECORD_TYPES.FOOD]: '🍚',
  [RECORD_TYPES.MEDICINE]: '💊',
} as const;
