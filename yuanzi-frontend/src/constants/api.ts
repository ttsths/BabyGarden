// API 端点常量
export const API_BASE =
  import.meta.env.VITE_API_BASE_URL ||
  'https://yuanzi-backend.shentuhaisan.workers.dev/api/v1';

export const ENDPOINTS = {
  // 认证
  AUTH: {
    SEND_CODE: '/auth/send-code',
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    PROFILE: '/user/profile',
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
    CONFIRM: '/photo/confirm',
    DELETE: (id: string) => `/photo/${id}`,
    COMMENTS: (id: string) => `/photo/${id}/comments`,
    LIKE: (id: string) => `/photo/${id}/like`,
    QUOTA: '/photo/quota',
  },
  // AI
  AI: {
    CHAT: '/ai/chat',
    SPEECH: '/ai/speech/recognize',
    QUOTA: '/ai/quota',
    HISTORY: '/ai/history',
  },
  // 家庭
  FAMILY: {
    CREATE: '/family',
    DETAIL: (id: string) => `/family/${id}`,
    MEMBERS: (id: string) => `/family/${id}/members`,
    INVITE: (id: string) => `/family/${id}/invite`,
    JOIN: '/family/join',
    LEAVE: (id: string) => `/family/${id}/leave`,
  },
  // 同步
  SYNC: {
    STREAM: '/sync/stream',
  },
} as const;
