import axios from 'axios';
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import { ENDPOINTS, API_BASE } from '@/constants/api';
import { useAuthStore } from '@/stores/useAuthStore';

// 创建 Axios 实例
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  async (error: AxiosError) => {
    // 401 未授权，跳转登录
    if (error.response?.status === 401) {
      useAuthStore.getState().logout();
      return Promise.reject(error);
    }

    // 500 服务器错误
    if (error.response?.status === 500) {
      console.error('Server error:', error.message);
    }

    return Promise.reject(error);
  }
);

export default apiClient;

// API 服务方法
export const api = {
  // 认证
  auth: {
    sendCode: (phone: string) => 
      apiClient.post(ENDPOINTS.AUTH.SEND_CODE, { phone }),
    login: (phone: string, code: string) => 
      apiClient.post(ENDPOINTS.AUTH.LOGIN, { phone, code }),
    logout: () => 
      apiClient.post(ENDPOINTS.AUTH.LOGOUT),
    getProfile: () => 
      apiClient.get(ENDPOINTS.AUTH.PROFILE),
  },
  
  // 宝宝
  baby: {
    getList: () => 
      apiClient.get(ENDPOINTS.BABY.LIST),
    getDetail: (id: string) => 
      apiClient.get(ENDPOINTS.BABY.DETAIL(id)),
    create: (data: Record<string, unknown>) => 
      apiClient.post(ENDPOINTS.BABY.CREATE, data),
    update: (id: string, data: Record<string, unknown>) => 
      apiClient.put(ENDPOINTS.BABY.UPDATE(id), data),
  },
  
  // 记录
  record: {
    getList: (babyId: string, params?: Record<string, unknown>) => 
      apiClient.get(ENDPOINTS.RECORD.LIST, { params: { baby_id: babyId, ...params } }),
    getDetail: (id: string) => 
      apiClient.get(ENDPOINTS.RECORD.DETAIL(id)),
    create: (data: Record<string, unknown>) => 
      apiClient.post(ENDPOINTS.RECORD.CREATE, data),
    update: (id: string, data: Record<string, unknown>) => 
      apiClient.put(ENDPOINTS.RECORD.UPDATE(id), data),
    delete: (id: string) => 
      apiClient.delete(ENDPOINTS.RECORD.DELETE(id)),
    getStats: (babyId: string, range: Record<string, string>) => 
      apiClient.get(ENDPOINTS.RECORD.STATS, { params: { babyId, ...range } }),
  },
  
  // 照片
  photo: {
    getList: (babyId: string, params?: Record<string, unknown>) => 
      apiClient.get(ENDPOINTS.PHOTO.LIST, { params: { baby_id: babyId, ...params } }),
    getUploadUrl: (data: Record<string, unknown>) => 
      apiClient.post(ENDPOINTS.PHOTO.UPLOAD_URL, data),
    confirmUpload: (photoId: string, size?: number) =>
      apiClient.post(ENDPOINTS.PHOTO.CONFIRM, { photo_id: photoId, size }),
    delete: (id: string) => 
      apiClient.delete(ENDPOINTS.PHOTO.DELETE(id)),
    getComments: (id: string) =>
      apiClient.get(ENDPOINTS.PHOTO.COMMENTS(id)),
    comment: (id: string, content: string) =>
      apiClient.post(ENDPOINTS.PHOTO.COMMENTS(id), { content }),
    like: (id: string) =>
      apiClient.post(ENDPOINTS.PHOTO.LIKE(id)),
    unlike: (id: string) =>
      apiClient.delete(ENDPOINTS.PHOTO.LIKE(id)),
    getQuota: () => 
      apiClient.get(ENDPOINTS.PHOTO.QUOTA),
  },
  
  // AI
  ai: {
    chat: (question: string, options?: { baby_id?: string; history?: Array<{ role: string; content: string }> }) =>
      apiClient.post(ENDPOINTS.AI.CHAT, { question, ...options }),
    speech: (audioBlob: Blob) => {
      const formData = new FormData();
      formData.append('audio', audioBlob);
      return apiClient.post(ENDPOINTS.AI.SPEECH, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },
    getQuota: () => 
      apiClient.get(ENDPOINTS.AI.QUOTA),
    history: (params?: Record<string, unknown>) =>
      apiClient.get(ENDPOINTS.AI.HISTORY, { params }),
  },
  
  // 家庭
  family: {
    create: (name: string) =>
      apiClient.post(ENDPOINTS.FAMILY.CREATE, { name }),
    getDetail: (familyId: string) =>
      apiClient.get(ENDPOINTS.FAMILY.DETAIL(familyId)),
    getMembers: (familyId: string) =>
      apiClient.get(ENDPOINTS.FAMILY.MEMBERS(familyId)),
    invite: (familyId: string, phone: string, role: 'member' | 'elder') =>
      apiClient.post(ENDPOINTS.FAMILY.INVITE(familyId), { phone, role }),
    join: (inviteCode: string, role: 'member' | 'elder' = 'member') =>
      apiClient.post(ENDPOINTS.FAMILY.JOIN, { invite_code: inviteCode, role }),
    leave: (familyId: string) =>
      apiClient.delete(ENDPOINTS.FAMILY.LEAVE(familyId)),
  },

  // 统计
  stats: {
    daily: (babyId: string, date?: string) =>
      apiClient.get('/stats/daily', { params: { baby_id: babyId, date } }),
    weekly: (babyId: string, date?: string) =>
      apiClient.get('/stats/weekly', { params: { baby_id: babyId, date } }),
    monthly: (babyId: string, date?: string) =>
      apiClient.get('/stats/monthly', { params: { baby_id: babyId, date } }),
    range: (babyId: string, startDate: string, endDate: string) =>
      apiClient.get('/stats/range', { params: { baby_id: babyId, start_date: startDate, end_date: endDate } }),
  },
};
