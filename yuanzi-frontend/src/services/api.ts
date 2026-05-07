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
      window.location.href = '/login';
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
    create: (data: any) => 
      apiClient.post(ENDPOINTS.BABY.CREATE, data),
    update: (id: string, data: any) => 
      apiClient.put(ENDPOINTS.BABY.UPDATE(id), data),
  },
  
  // 记录
  record: {
    getList: (babyId: string, params?: any) => 
      apiClient.get(ENDPOINTS.RECORD.LIST, { params: { babyId, ...params } }),
    getDetail: (id: string) => 
      apiClient.get(ENDPOINTS.RECORD.DETAIL(id)),
    create: (data: any) => 
      apiClient.post(ENDPOINTS.RECORD.CREATE, data),
    update: (id: string, data: any) => 
      apiClient.put(ENDPOINTS.RECORD.UPDATE(id), data),
    delete: (id: string) => 
      apiClient.delete(ENDPOINTS.RECORD.DELETE(id)),
    getStats: (babyId: string, range: any) => 
      apiClient.get(ENDPOINTS.RECORD.STATS, { params: { babyId, ...range } }),
  },
  
  // 照片
  photo: {
    getList: (babyId: string, params?: any) => 
      apiClient.get(ENDPOINTS.PHOTO.LIST, { params: { babyId, ...params } }),
    getUploadUrl: (data: any) => 
      apiClient.post(ENDPOINTS.PHOTO.UPLOAD_URL, data),
    confirmUpload: (id: string) => 
      apiClient.post(ENDPOINTS.PHOTO.CONFIRM(id)),
    delete: (id: string) => 
      apiClient.delete(ENDPOINTS.PHOTO.DELETE(id)),
    getQuota: () => 
      apiClient.get(ENDPOINTS.PHOTO.QUOTA),
  },
  
  // AI
  ai: {
    chat: (message: string, context?: any) => 
      apiClient.post(ENDPOINTS.AI.CHAT, { message, context }),
    speech: (audioBlob: Blob) => {
      const formData = new FormData();
      formData.append('audio', audioBlob);
      return apiClient.post(ENDPOINTS.AI.SPEECH, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },
    getQuota: () => 
      apiClient.get(ENDPOINTS.AI.QUOTA),
  },
  
  // 家庭
  family: {
    getDetail: () => 
      apiClient.get(ENDPOINTS.FAMILY.DETAIL),
    getMembers: () => 
      apiClient.get(ENDPOINTS.FAMILY.MEMBERS),
    invite: (phone: string) => 
      apiClient.post(ENDPOINTS.FAMILY.INVITE, { phone }),
  },
};
