// API 端点常量
export const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api/v1';

export const ENDPOINTS = {
  // 认证
  AUTH: {
    SEND_CODE: '/auth/send-code',
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    PROFILE: '/auth/profile',
  },
  // 宝宝
  BABY: {
    LIST: '/baby',
    DETAIL: (id: string) => `/baby/${id}`,
    CREATE: '/baby',
    UPDATE: (id: string) => `/baby/${id}`,
  },
  // 记录
  RECORD: {
    LIST: '/record',
    DETAIL: (id: string) => `/record/${id}`,
    CREATE: '/record',
    UPDATE: (id: string) => `/record/${id}`,
    DELETE: (id: string) => `/record/${id}`,
    STATS: '/record/stats',
  },
  // 照片
  PHOTO: {
    LIST: '/photo',
    UPLOAD_URL: '/photo/upload-url',
    CONFIRM: (id: string) => `/photo/${id}/confirm`,
    DELETE: (id: string) => `/photo/${id}`,
    QUOTA: '/photo/quota',
  },
  // AI
  AI: {
    CHAT: '/ai/chat',
    SPEECH: '/ai/speech',
    QUOTA: '/ai/quota',
  },
  // 家庭
  FAMILY: {
    DETAIL: '/family',
    MEMBERS: '/family/members',
    INVITE: '/family/invite',
  },
  // 同步
  SYNC: {
    STREAM: '/sync/stream',
  },
} as const;
