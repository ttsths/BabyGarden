export interface AdminLoginResponse {
  token: string;
}

export interface OverviewStats {
  users: number;
  families: number;
  babies: number;
  photos: number;
  records: number;
}

export interface DailyStat {
  date: string;
  new_users: number;
  new_babies: number;
  new_records: number;
}

export interface AdminUser {
  id: string;
  phone: string;
  nickname: string;
  avatar_url: string;
  status: number;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

export interface AdminBaby {
  id: string;
  family_id: string;
  name: string;
  avatar_url: string;
  birthday: string;
  gender: string;
  created_at: string;
  updated_at: string;
}

export interface AdminFamily {
  id: string;
  name: string;
  owner_id: string;
  storage_quota: number;
  storage_used: number;
  created_at: string;
  updated_at: string;
}

export interface AdminFamilyMember {
  user_id: string;
  role: string;
  nickname: string;
  avatar_url: string;
  joined_at: string;
}

export interface AdminFamilyDetail extends AdminFamily {
  members: AdminFamilyMember[];
}

export interface AdminPhoto {
  id: string;
  family_id: string;
  baby_id: string;
  uploader_id: string;
  filename: string;
  thumbnail_url: string;
  original_url: string;
  size: number;
  content_type: string;
  created_at: string;
}

export interface AdminRecord {
  id: string;
  family_id: string;
  baby_id: string;
  type: string;
  data: Record<string, unknown>;
  created_by: string;
  created_at: string;
}

export interface PaginatedResponse<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

export interface ApiResult<T> {
  code: number;
  msg: string;
  data: T;
}
