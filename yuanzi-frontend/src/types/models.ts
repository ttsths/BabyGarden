// 用户相关
export interface User {
  id: string;
  phone: string;
  nickname?: string;
  avatar?: string;
  isPaid: boolean;
  createdAt: string;
}

// 宝宝信息
export interface Baby {
  id: string;
  familyId: string;
  name: string;
  avatar?: string;
  birthday: string;
  gender: 'male' | 'female';
  createdAt: string;
}

// 记录类型
export type RecordType = 'feeding' | 'sleep' | 'diaper' | 'temperature' | 'food' | 'medicine';

// 记录
export interface Record {
  id: string;
  familyId: string;
  babyId: string;
  type: RecordType;
  data: RecordData;
  createdAt: string;
  createdBy: string;
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
  familyId: string;
  babyId: string;
  uploadedBy: string;
  ossKey: string;
  originalName: string;
  thumbnailUrl: string;
  originalUrl: string;
  size: number;
  contentType: string;
  status: 'pending' | 'active' | 'failed';
  createdAt: string;
}

// 家庭
export interface Family {
  id: string;
  name: string;
  members: FamilyMember[];
  isPaid: boolean;
  storageQuota: number;
  storageUsed: number;
}

export interface FamilyMember {
  userId: string;
  role: 'admin' | 'member';
  joinedAt: string;
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
