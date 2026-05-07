/**
 * Yuanzi 首页组件
 * 基于 Stitch 设计：yuanzi_baby_app_home
 * 主色：珊瑚粉 #ff998a | 背景：暖白 #fffbf7
 */

import React from 'react';
import { cn } from '../utils/cn';

interface HomePageProps {
  userName: string;
  date: string;
  stats: {
    feedingCount: number;
    sleepHours: number;
    diaperCount: number;
  };
  recentRecords: Array<{
    id: string;
    type: 'feeding' | 'sleep' | 'diaper';
    title: string;
    subtitle: string;
    time: string;
  }>;
  onRecordTypeClick: (type: 'feeding' | 'sleep' | 'diaper') => void;
  onViewAllRecords: () => void;
  onRecordClick: (id: string) => void;
}

export const HomePage: React.FC<HomePageProps> = ({
  userName,
  date,
  stats,
  recentRecords,
  onRecordTypeClick,
  onViewAllRecords,
  onRecordClick,
}) => {
  return (
    <div className="min-h-screen bg-[#FDF8F5] font-display">
      <div className="relative flex h-full min-h-screen w-full max-w-[430px] mx-auto flex-col bg-[#FDF8F5] shadow-2xl overflow-x-hidden">
        
        {/* Header - 头像 + 问候 + 日期 + 日历按钮 */}
        <header className="flex items-center justify-between p-6 pb-2 mt-4">
          <div className="flex items-center gap-3">
            <div className="size-12 rounded-full bg-[#FFE8E0] flex items-center justify-center overflow-hidden border-2 border-[#F4A896]/20">
              <img
                alt="User Profile"
                className="w-full h-full object-cover"
                src="https://lh3.googleusercontent.com/aida-public/AB6AXuCt_GBLYbhErVSbJcUqW-5VA1OBkV4ssC8PzxQLPhl2h9tkwk78zZ-AoUFnsK-P3owNTffpE6K1GI5D3SEVD1XEdDFDYdLHD0geYnn7VssTA3i3z36kli6wKk9rYXC8WhUSZKfXppiB4n0CQsXQpVUBSRkifgCQE4N3Sypx8IXf4S7rmw6t6hSi6PjrNmal3PPIqq8rtNvsQ7u3sihfMtcoRzHQ3iGO3Lf2E5Zd2M-HdeeGte3N5ZW9B-UNXj0vRfbWT6hcQ0JbKGc"
              />
            </div>
            <div>
              <h2 className="text-xl font-bold leading-tight tracking-tight text-[#1A1A1A]">
                Hello, {userName}!
              </h2>
              <p className="text-[#F4A896] text-sm font-medium">{date}</p>
            </div>
          </div>
          <button className="size-10 rounded-full bg-white shadow-sm flex items-center justify-center text-[#666666]">
            <span className="material-symbols-outlined">calendar_today</span>
          </button>
        </header>

        {/* Today's Stats - 2x2 网格 */}
        <section className="px-6 py-4">
          <h3 className="text-lg font-bold mb-4 text-[#1A1A1A]">Today's Stats</h3>
          <div className="grid grid-cols-2 gap-4">
            <StatCard
              icon="child_care"
              label="Feeding"
              value="5x"
            />
            <StatCard
              icon="bedtime"
              label="Sleep"
              value={`${stats.sleepHours}h`}
            />
            <StatCard
              icon="water_drop"
              label="Diapers"
              value={`${stats.diaperCount}x`}
            />
            <StatCard
              icon="emoji_events"
              label="Milestones"
              value="0"
            />
          </div>
        </section>

        {/* Quick Actions - 3 个圆形按钮 */}
        <section className="px-6 py-6 flex justify-center gap-6">
          <QuickActionButton
            icon="restaurant"
            label="Feed"
            onClick={() => onRecordTypeClick('feeding')}
          />
          <QuickActionButton
            icon="nights_stay"
            label="Sleep"
            onClick={() => onRecordTypeClick('sleep')}
          />
          <QuickActionButton
            icon="baby_changing_station"
            label="Diaper"
            onClick={() => onRecordTypeClick('diaper')}
          />
        </section>

        {/* Recent Activities - 列表 */}
        <section className="px-6 py-4 flex-1">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-bold text-[#1A1A1A]">Recent Activities</h3>
            <button
              onClick={onViewAllRecords}
              className="text-[#F4A896] text-sm font-bold"
            >
              See all
            </button>
          </div>
          <div className="space-y-3">
            {recentRecords.map((record) => (
              <ActivityItem
                key={record.id}
                record={record}
                onClick={() => onRecordClick(record.id)}
              />
            ))}
            {recentRecords.length === 0 && (
              <div className="text-center py-8">
                <div className="text-6xl mb-3">👶</div>
                <p className="text-[#999999] text-sm">开始记录吧</p>
              </div>
            )}
          </div>
        </section>

        {/* Bottom Navigation - 固定底部 */}
        <nav className="sticky bottom-0 w-full border-t border-[#F4A896]/10 bg-white/80 backdrop-blur-md px-6 py-3 flex justify-between items-center">
          <NavItem icon="home" label="Home" active />
          <NavItem icon="list_alt" label="Logs" />
          <NavItem icon="trending_up" label="Growth" />
          <NavItem icon="settings" label="Settings" />
        </nav>
      </div>
    </div>
  );
};

// 统计卡片组件
const StatCard: React.FC<{
  icon: string;
  label: string;
  value: string;
}> = ({ icon, label, value }) => {
  return (
    <div className="bg-white p-4 rounded-[16px] shadow-[0_2px_8px_rgba(0,0,0,0.06)] flex flex-col gap-2">
      <div className="size-10 rounded-full bg-[#FFF0ED] flex items-center justify-center text-[#F4A896]">
        <span className="material-symbols-outlined">{icon}</span>
      </div>
      <div>
        <p className="text-2xl font-bold text-[#1A1A1A]">{value}</p>
        <p className="text-[#999999] text-sm">{label}</p>
      </div>
    </div>
  );
};

// 快速操作按钮 - 圆形珊瑚粉按钮
const QuickActionButton: React.FC<{
  icon: string;
  label: string;
  onClick: () => void;
}> = ({ icon, label, onClick }) => {
  return (
    <div className="flex flex-col items-center gap-2">
      <button
        onClick={onClick}
        className="size-[64px] rounded-full bg-[#F4A896] text-white shadow-[0_4px_14px_rgba(244,168,150,0.4)] flex items-center justify-center active:scale-95 transition-transform"
      >
        <span className="material-symbols-outlined text-[28px]">{icon}</span>
      </button>
      <span className="text-sm font-semibold text-[#1A1A1A]">{label}</span>
    </div>
  );
};

// 活动项 - 横向布局
const ActivityItem: React.FC<{
  record: {
    id: string;
    type: 'feeding' | 'sleep' | 'diaper';
    title: string;
    subtitle: string;
    time: string;
  };
  onClick: () => void;
}> = ({ record, onClick }) => {
  const getIconConfig = () => {
    switch (record.type) {
      case 'feeding':
        return { icon: 'restaurant', bg: 'bg-[#FFB59A]' }; // 橙色
      case 'sleep':
        return { icon: 'bedtime', bg: 'bg-[#A8B5E8]' }; // 紫色
      case 'diaper':
        return { icon: 'water_drop', bg: 'bg-[#8ECAFF]' }; // 蓝色
      default:
        return { icon: 'note', bg: 'bg-slate-400' };
    }
  };

  const config = getIconConfig();

  return (
    <div
      onClick={onClick}
      className="flex items-center gap-4 p-4 cursor-pointer"
    >
      <div className={cn('size-11 rounded-full flex items-center justify-center flex-shrink-0', config.bg)}>
        <span className="material-symbols-outlined text-white text-lg">{config.icon}</span>
      </div>
      <div className="flex-1">
        <p className="font-semibold text-base text-[#1A1A1A]">{record.title}</p>
        <p className="text-sm text-[#999999] mt-0.5">{record.subtitle}</p>
      </div>
      <span className="text-sm text-[#999999] flex-shrink-0">{record.time}</span>
    </div>
  );
};

// 底部导航项
const NavItem: React.FC<{
  icon: string;
  label: string;
  active?: boolean;
}> = ({ icon, label, active = false }) => (
  <a href="#" className={cn('flex flex-col items-center gap-1', active ? 'text-[#F4A896]' : 'text-[#999999]')}>
    <span className={cn('material-symbols-outlined', active ? 'text-fill-1' : '')}>{icon}</span>
    <span className={cn('text-[10px]', active ? 'font-bold' : 'font-medium')}>{label}</span>
  </a>
);
