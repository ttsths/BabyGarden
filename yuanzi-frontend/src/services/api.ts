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
    passwordLogin: (identifier: string, password: string) =>
      unwrap<{ access_token: string; refresh_token: string }>(apiClient.post(ENDPOINTS.AUTH.PASSWORD_LOGIN, { identifier, password })),
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
    getWeeklyStats: (babyId: string, date?: string) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.STATS_WEEKLY, { params: { baby_id: babyId, date } })),
    getMonthlyStats: (babyId: string, date?: string) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.STATS_SUMMARY, { params: { baby_id: babyId, range: 'month', end_date: date } })),
    getRangeStats: (babyId: string, startDate: string, endDate: string) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.STATS_SUMMARY, { params: { baby_id: babyId, range: 'custom', start_date: startDate, end_date: endDate } })),
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
    chat: (question: string, babyIdOrOptions?: string | { baby_id?: string; history?: Array<{ role: string; content: string }> }, history?: Array<{ role: string; content: string }>) => {
      const payload = typeof babyIdOrOptions === 'object'
        ? { question, ...babyIdOrOptions }
        : { question, baby_id: babyIdOrOptions, history };
      return unwrap(apiClient.post(ENDPOINTS.AI.CHAT, payload));
    },
    chatStream: async (
      question: string,
      options: { baby_id?: string; history?: Array<{ role: string; content: string }>; onDelta?: (delta: string) => void }
    ) => {
      const token = useAuthStore.getState().token;
      const response = await fetch(`${API_BASE}${ENDPOINTS.AI.CHAT_STREAM}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify({ question, baby_id: options.baby_id, history: options.history }),
      });
      if (!response.ok || !response.body) {
        throw new Error('AI stream failed');
      }
      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';
      let answer = '';
      let donePayload: Record<string, unknown> | null = null;
      while (true) {
        const { value, done } = await reader.read();
        if (done) break;
        buffer += decoder.decode(value, { stream: true });
        const events = buffer.split('\n\n');
        buffer = events.pop() ?? '';
        for (const eventText of events) {
          const lines = eventText.split('\n');
          const event = lines.find((line) => line.startsWith('event: '))?.slice(7);
          const dataLine = lines.find((line) => line.startsWith('data: '));
          if (!dataLine) continue;
          const data = JSON.parse(dataLine.slice(6)) as Record<string, unknown>;
          if (event === 'delta') {
            const delta = String(data.delta ?? '');
            answer += delta;
            options.onDelta?.(delta);
          } else if (event === 'done') {
            donePayload = data;
          } else if (event === 'error') {
            throw new Error(String(data.message ?? 'AI stream failed'));
          }
        }
      }
      return { answer, ...donePayload };
    },
    getHistory: (params?: Record<string, unknown>) => unwrap(apiClient.get(ENDPOINTS.AI.CHATS, { params })),
    history: (params?: Record<string, unknown>) => unwrap(apiClient.get(ENDPOINTS.AI.HISTORY, { params })),
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
    join: (inviteCode: string, role = 'member') => unwrap(apiClient.post(ENDPOINTS.FAMILY.JOIN, { invite_code: inviteCode, role })),
    leave: (id: string) => unwrap(apiClient.delete(ENDPOINTS.FAMILY.LEAVE(id))),
  },
  stats: {
    daily: (babyId: string, date?: string) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.STATS_DAILY, { params: { baby_id: babyId, date } })),
    weekly: (babyId: string, date?: string) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.STATS_WEEKLY, { params: { baby_id: babyId, date } })),
    monthly: (babyId: string, date?: string) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.STATS_SUMMARY, { params: { baby_id: babyId, range: 'month', end_date: date } })),
    range: (babyId: string, startDate: string, endDate: string) =>
      unwrap(apiClient.get(ENDPOINTS.RECORD.STATS_SUMMARY, { params: { baby_id: babyId, range: 'custom', start_date: startDate, end_date: endDate } })),
  },
};
