import axios from 'axios';
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import type {
  ApiResult,
  AdminLoginResponse,
  OverviewStats,
  DailyStat,
  AdminUser,
  AdminBaby,
  AdminFamily,
  AdminFamilyDetail,
  AdminPhoto,
  AdminRecord,
  PaginatedResponse,
  CreateUserRequest,
  UpdateUserRequest,
  CreateBabyRequest,
  UpdateBabyRequest,
  CreateFamilyRequest,
  UpdateFamilyRequest,
  AddFamilyMemberRequest,
  CreateRecordRequest,
  UpdateRecordRequest,
} from '@/admin/types/admin';

const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api/v1';

const adminClient: AxiosInstance = axios.create({
  baseURL: `${API_BASE}/admin`,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

adminClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('admin_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

adminClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('admin_token');
      window.location.href = '/admin/login';
    }
    return Promise.reject(error);
  }
);

export function adminLogin(phone: string, password: string) {
  return adminClient.post<ApiResult<AdminLoginResponse>>('/login', { phone, password });
}

export function getStatsOverview() {
  return adminClient.get<ApiResult<OverviewStats>>('/stats/overview');
}

export function getDailyStats(days = 30) {
  return adminClient.get<ApiResult<DailyStat[]>>('/stats/daily', { params: { days } });
}

export function getUsers(page = 1, pageSize = 20, keyword?: string) {
  return adminClient.get<ApiResult<PaginatedResponse<AdminUser>>>('/users', {
    params: { page, page_size: pageSize, keyword },
  });
}

export function getUserDetail(id: string) {
  return adminClient.get<ApiResult<AdminUser>>(`/users/${id}`);
}

export function createUser(data: CreateUserRequest) {
  return adminClient.post<ApiResult<AdminUser>>('/users', data);
}

export function updateUser(id: string, data: UpdateUserRequest) {
  return adminClient.put<ApiResult<AdminUser>>(`/users/${id}`, data);
}

export function updateUserStatus(id: string, status: number) {
  return adminClient.put<ApiResult<unknown>>(`/users/${id}/status`, { status });
}

export function deleteUser(id: string) {
  return adminClient.delete<ApiResult<unknown>>(`/users/${id}`);
}

export function getBabies(page = 1, pageSize = 20) {
  return adminClient.get<ApiResult<PaginatedResponse<AdminBaby>>>('/babies', {
    params: { page, page_size: pageSize },
  });
}

export function getBabyDetail(id: string) {
  return adminClient.get<ApiResult<AdminBaby>>(`/babies/${id}`);
}

export function createBaby(data: CreateBabyRequest) {
  return adminClient.post<ApiResult<AdminBaby>>('/babies', data);
}

export function updateBaby(id: string, data: UpdateBabyRequest) {
  return adminClient.put<ApiResult<AdminBaby>>(`/babies/${id}`, data);
}

export function deleteBaby(id: string) {
  return adminClient.delete<ApiResult<unknown>>(`/babies/${id}`);
}

export function getFamilies(page = 1, pageSize = 20) {
  return adminClient.get<ApiResult<PaginatedResponse<AdminFamily>>>('/families', {
    params: { page, page_size: pageSize },
  });
}

export function getFamilyDetail(id: string) {
  return adminClient.get<ApiResult<AdminFamilyDetail>>(`/families/${id}`);
}

export function createFamily(data: CreateFamilyRequest) {
  return adminClient.post<ApiResult<AdminFamily>>('/families', data);
}

export function updateFamily(id: string, data: UpdateFamilyRequest) {
  return adminClient.put<ApiResult<AdminFamily>>(`/families/${id}`, data);
}

export function addFamilyMember(id: string, data: AddFamilyMemberRequest) {
  return adminClient.post<ApiResult<unknown>>(`/families/${id}/members`, data);
}

export function removeFamilyMember(id: string, userId: string) {
  return adminClient.delete<ApiResult<unknown>>(`/families/${id}/members/${userId}`);
}

export function deleteFamily(id: string) {
  return adminClient.delete<ApiResult<unknown>>(`/families/${id}`);
}

export function getPhotos(page = 1, pageSize = 20) {
  return adminClient.get<ApiResult<PaginatedResponse<AdminPhoto>>>('/photos', {
    params: { page, page_size: pageSize },
  });
}

export function getPhotoDetail(id: string) {
  return adminClient.get<ApiResult<AdminPhoto>>(`/photos/${id}`);
}

export function deletePhoto(id: string) {
  return adminClient.delete<ApiResult<unknown>>(`/photos/${id}`);
}

export function getPhotoUploadUrl(data: { baby_id: string; filename: string; content_type: string; size: number }) {
  return adminClient.post<ApiResult<{ upload_url: string; photo_id: string }>>('/photos/upload-url', data);
}

export function confirmPhotoUpload(data: { photo_id: string; size?: number }) {
  return adminClient.post<ApiResult<unknown>>('/photos/confirm', data);
}

export function getRecords(page = 1, pageSize = 20) {
  return adminClient.get<ApiResult<PaginatedResponse<AdminRecord>>>('/records', {
    params: { page, page_size: pageSize },
  });
}

export function getRecordDetail(id: string) {
  return adminClient.get<ApiResult<AdminRecord>>(`/records/${id}`);
}

export function createRecord(data: CreateRecordRequest) {
  return adminClient.post<ApiResult<AdminRecord>>('/records', data);
}

export function updateRecord(id: string, data: UpdateRecordRequest) {
  return adminClient.put<ApiResult<AdminRecord>>(`/records/${id}`, data);
}

export function deleteRecord(id: string) {
  return adminClient.delete<ApiResult<unknown>>(`/records/${id}`);
}
