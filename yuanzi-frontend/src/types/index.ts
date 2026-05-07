// 基础类型定义

export interface Baby {
  id: string;
  name: string;
  birthday: string;
  gender: 1 | 2; // 1-男 2-女
  avatarUrl?: string;
  note?: string;
  age?: string;
}

export interface Record {
  id: string;
  babyId: string;
  type: 'feeding' | 'sleep' | 'diaper' | 'growth';
  startedAt: string;
  endedAt?: string;
  content: RecordContent;
  note?: string;
  createdAt: string;
}

export interface RecordContent {
  // 喂养记录
  type?: 'breast' | 'formula' | 'solid';
  side?: 'left' | 'right' | 'both';
  duration?: number; // 分钟
  amount?: number; // ml 或 g
  unit?: 'ml' | 'g';
  
  // 睡眠记录
  quality?: 'good' | 'normal' | 'poor';
  location?: 'crib' | 'bed' | 'car' | 'stroller';
  
  // 排泄记录
  diaperType?: 'wet' | 'dirty' | 'both';
  color?: 'yellow' | 'green' | 'brown';
  consistency?: 'normal' | 'watery' | 'hard';
  
  // 成长记录
  weight?: number; // kg
  height?: number; // cm
  headCircumference?: number; // cm
}

export interface DailyStats {
  date: string;
  feedingCount: number;
  sleepHours: number;
  diaperCount: number;
  totalRecords: number;
}

export interface RecordType {
  id: string;
  name: string;
  icon: string;
  color: string;
}

export const RECORD_TYPES: RecordType[] = [
  { id: 'feeding', name: '喂奶', icon: '🍼', color: '#FF9A8B' },
  { id: 'sleep', name: '睡觉', icon: '💤', color: '#D8BFD8' },
  { id: 'diaper', name: '换尿布', icon: '💩', color: '#A2D5C6' },
  { id: 'growth', name: '成长', icon: '📏', color: '#FFB4A2' },
];

export interface TabItem {
  key: string;
  label: string;
  icon: string;
}
