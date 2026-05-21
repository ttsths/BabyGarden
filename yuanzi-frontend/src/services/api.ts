import axios from 'axios';
import type { AxiosError, AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { API_BASE, ENDPOINTS } from '@/constants/api';
import { useAuthStore } from '@/stores/useAuthStore';

interface BackendResponse<T> {
  code: number;
  msg: string;
  data: T;
}

export interface ListResponse<T> {
  list: T[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
});

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

apiClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

function unwrap<T>(request: Promise<AxiosResponse<BackendResponse<T>>>): Promise<T> {
  return request.then((response) => response.data.data);
}

export default apiClient;

export const api = {
  auth: {
    sendCode: (phone: string) => unwrap<unknown>(apiClient.post(ENDPOINTS.AUTH.SEND_CODE, { phone })),
    login: (phone: string, code: string) =>
      unwrap<{ access_token: string; refresh_token: string }>(apiClient.post(ENDPOINTS.AUTH.LOGIN, { phone, code })),
    logout: () => unwrap<unknown>(apiClient.post(ENDPOINTS.AUTH.LOGOUT)),
    getProfile: () => unwrap(apiClient.get(ENDPOINTS.AUTH.PROFILE)),
  },

  baby: {
    getList: () => unwrap(apiClient.get(ENDPOINTS.BABY.LIST)),
    getDetail: (id: string) => unwrap(apiClient.get(ENDPOINTS.BABY.DETAIL(id))),
    create: (data: Record<string, unknown>) => unwrap(apiClient.post(ENDPOINTS.BABY.CREATE, data)),
    update: (id: string, data: Record<string, unknown>) => unwrap(apiClient.put(ENDPOINTS.BABY.UPDATE(id), data)),
  },

  record: {
    getList: (babyId: string, params?: Record<string, unknown>) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.LIST, { params: { baby_id: babyId, ...params } })),
    getDetail: (id: string) => unwrap(apiClient.get(ENDPOINTS.RECORD.DETAIL(id))),
    create: (data: Record<string, unknown>) => unwrap(apiClient.post(ENDPOINTS.RECORD.CREATE, data)),
    update: (id: string, data: Record<string, unknown>) => unwrap(apiClient.put(ENDPOINTS.RECORD.UPDATE(id), data)),
    delete: (id: string) => unwrap(apiClient.delete(ENDPOINTS.RECORD.DELETE(id))),
    getDailyStats: (babyId: string, date?: string) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.STATS_DAILY, { params: { baby_id: babyId, date } })),
    getSummaryStats: (babyId: string, params: Record<string, string>) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.STATS_SUMMARY, { params: { baby_id: babyId, ...params } })),
  },

  photo: {
    getList: (babyId: string, params?: Record<string, unknown>) =>
      unwrap(apiClient.get(ENDPOINTS.PHOTO.LIST, { params: { baby_id: babyId, ...params } })),
    getUploadUrl: (data: Record<string, unknown>) => unwrap(apiClient.post(ENDPOINTS.PHOTO.UPLOAD_URL, data)),
    confirmUpload: (photoId: string, size?: number) =>
      unwrap(apiClient.post(ENDPOINTS.PHOTO.CONFIRM, { photo_id: photoId, size })),
    delete: (id: string) => unwrap(apiClient.delete(ENDPOINTS.PHOTO.DELETE(id))),
    getComments: (id: string) => unwrap(apiClient.get(ENDPOINTS.PHOTO.COMMENTS(id))),
    comment: (id: string, content: string) => unwrap(apiClient.post(ENDPOINTS.PHOTO.COMMENTS(id), { content })),
    like: (id: string) => unwrap(apiClient.post(ENDPOINTS.PHOTO.LIKE(id))),
    unlike: (id: string) => unwrap(apiClient.delete(ENDPOINTS.PHOTO.LIKE(id))),
  },

  ai: {
    chat: (question: string, babyId?: string, history?: Array<{ role: string; content: string }>) =>
      unwrap(apiClient.post(ENDPOINTS.AI.CHAT, { question, baby_id: babyId, history })),
    getHistory: (params?: Record<string, unknown>) => unwrap(apiClient.get(ENDPOINTS.AI.CHATS, { params })),
    getDetail: (id: string) => unwrap(apiClient.get(ENDPOINTS.AI.CHAT_DETAIL(id))),
    speech: (audioBlob: Blob) => {
      const formData = new FormData();
      formData.append('audio', audioBlob);
      return unwrap(apiClient.post(ENDPOINTS.AI.SPEECH, formData, { headers: { 'Content-Type': 'multipart/form-data' } }));
    },
    getQuota: () => unwrap(apiClient.get(ENDPOINTS.AI.QUOTA)),
  },

  family: {
    getDetail: (id: string) => unwrap(apiClient.get(ENDPOINTS.FAMILY.DETAIL(id))),
    getMembers: (id: string) => unwrap(apiClient.get(ENDPOINTS.FAMILY.MEMBERS(id))),
    invite: (id: string, phone: string, role = 'member') =>
      unwrap(apiClient.post(ENDPOINTS.FAMILY.INVITE(id), { phone, role })),
    join: (inviteCode: string) => unwrap(apiClient.post(ENDPOINTS.FAMILY.JOIN, { invite_code: inviteCode })),
    leave: (id: string) => unwrap(apiClient.post(ENDPOINTS.FAMILY.LEAVE(id))),
  },
};
