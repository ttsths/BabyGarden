// 用户相关
export interface User {
  id: string;
  phone: string;
  username?: string;
  nickname?: string;
  avatar_url?: string;
  status?: number;
  is_admin?: number;
}

// 宝宝信息
export interface Baby {
  id: string;
  family_id: string;
  name: string;
  avatar_url?: string;
  birthday: string;
  gender: number | 'male' | 'female';
  created_at?: string;
}

// 记录类型
export type RecordType = 'feeding' | 'sleep' | 'diaper' | 'excretion' | 'temperature' | 'growth';

// 记录
export interface Record {
  id: string;
  family_id?: string;
  baby_id: string;
  type: RecordType;
  started_at: string;
  ended_at?: string;
  content: globalThis.Record<string, unknown>;
  note?: string;
  created_at?: string;
  created_by?: string;
}

export interface FeedingData {
  type: 'breast' | 'formula';
  side?: 'left' | 'right';
  amount?: number; // ml
  duration?: number; // minutes
}

export interface SleepData {
  startTime: string;
  endTime?: string;
  duration?: number; // minutes
}

export interface DiaperData {
  type: 'pee' | 'poop' | 'mixed';
  note?: string;
}

export type RecordData = FeedingData | SleepData | DiaperData;

// 照片
export interface Photo {
  id: string;
  url: string;
  thumb_url: string;
  width: number;
  height: number;
  taken_at: string;
  description?: string;
  like_count: number;
  comment_count: number;
  liked_by_me: boolean;
}

// 家庭
export interface Family {
  id: string;
  name: string;
  invite_code: string;
  is_paid: boolean;
  storage_limit: number;
  storage_used: number;
}

export interface FamilyMember {
  user_id: string;
  nickname: string;
  avatar_url?: string;
  role: 'admin' | 'member' | 'elder';
  elder_mode: boolean;
}

// API 响应
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
  };
}

// 分页
export interface PageParams {
  page: number;
  pageSize: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  pageSize: number;
}
