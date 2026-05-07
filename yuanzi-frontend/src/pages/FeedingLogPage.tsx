/**
 * 喂养记录列表页
 * 基于 Stitch 设计：yuanzi_feeding_log_chinese
 */

import React from 'react';
import { cn } from '../utils/cn';

interface FeedingLogPageProps {
  onBack: () => void;
  onAddRecord?: () => void;
}

interface FeedingRecord {
  id: string;
  time: string;
  type: string;
  amount?: number;
  duration?: number;
  note?: string;
}

export const FeedingLogPage: React.FC<FeedingLogPageProps> = ({ onBack, onAddRecord }) => {
  // 示例数据
  const [todayStats] = React.useState({
    totalCount: 5,
    averageAmount: 120,
  });

  const [records] = React.useState<FeedingRecord[]>([
    {
      id: '1',
      time: '06:30',
      type: '配方奶',
      amount: 120,
      note: '刚醒',
    },
    {
      id: '2',
      time: '09:15',
      type: '母乳',
      duration: 15,
      note: '左侧',
    },
    {
      id: '3',
      time: '12:00',
      type: '配方奶',
      amount: 150,
      note: '睡醒后',
    },
    {
      id: '4',
      time: '15:30',
      type: '母乳',
      duration: 20,
      note: '双侧',
    },
    {
      id: '5',
      time: '18:45',
      type: '配方奶',
      amount: 130,
      note: '睡前喂',
    },
  ]);

  return (
    <div className="relative mx-auto max-w-md min-h-screen flex flex-col pb-24 bg-background-light dark:bg-background-dark">
      {/* Header */}
      <header className="sticky top-0 z-10 flex items-center justify-between bg-background-light/80 dark:bg-background-dark/80 backdrop-blur-md p-4">
        <button
          onClick={onBack}
          className="flex h-10 w-10 items-center justify-center rounded-full hover:bg-primary/10 transition-colors"
        >
          <span className="material-symbols-outlined text-slate-700 dark:text-slate-300">
            arrow_back_ios_new
          </span>
        </button>
        <h1 className="text-lg font-bold text-slate-800 dark:text-slate-100">
          喂奶记录
        </h1>
        <button className="flex h-10 w-10 items-center justify-center rounded-full hover:bg-primary/10 transition-colors">
          <span className="material-symbols-outlined text-slate-700 dark:text-slate-300">
            more_horiz
          </span>
        </button>
      </header>

      {/* Today's Summary Card */}
      <section className="px-4 py-2">
        <div className="bg-white dark:bg-slate-800/50 rounded-xl p-5 shadow-sm border border-primary/10">
          <div className="flex items-center gap-2 mb-4">
            <span className="material-symbols-outlined text-primary text-xl">insights</span>
            <h2 className="text-base font-bold text-slate-800 dark:text-slate-100">
              今日统计
            </h2>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="flex flex-col gap-1 p-3 rounded-lg bg-primary/5 dark:bg-primary/10 border border-primary/5">
              <p className="text-xs text-slate-500 dark:text-slate-400 font-medium">
                总计
              </p>
              <p className="text-2xl font-bold text-primary">
                {todayStats.totalCount} <span className="text-sm font-medium">次</span>
              </p>
            </div>
            <div className="flex flex-col gap-1 p-3 rounded-lg bg-primary/5 dark:bg-primary/10 border border-primary/5">
              <p className="text-xs text-slate-500 dark:text-slate-400 font-medium">
                平均
              </p>
              <div className="flex items-baseline gap-1">
                <p className="text-2xl font-bold text-primary">
                  {todayStats.averageAmount}
                </p>
                <span className="text-sm font-medium text-primary">ml</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Time-based Log List */}
      <section className="px-4 mt-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-sm font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
            喂奶明细
          </h3>
          <span className="text-xs text-primary font-medium">今天</span>
        </div>
        <div className="space-y-3">
          {records.map((record) => (
            <FeedingLogItem
              key={record.id}
              record={record}
            />
          ))}
        </div>
      </section>

      {/* Floating Action Button */}
      <button
        onClick={onAddRecord}
        className="fixed bottom-24 right-6 flex h-16 w-16 items-center justify-center rounded-full bg-primary text-white shadow-lg shadow-primary/30 hover:scale-105 active:scale-95 transition-all z-20"
      >
        <span className="material-symbols-outlined text-3xl">add</span>
      </button>

      {/* Bottom Navigation */}
      <nav className="fixed bottom-0 left-0 right-0 max-w-md mx-auto flex h-20 items-center justify-around border-t border-primary/10 bg-white/95 dark:bg-slate-900/95 backdrop-blur-md px-4 pb-4">
        <NavItem icon="home" label="首页" active={false} />
        <NavItem icon="history_edu" label="记录" active />
        <NavItem icon="bar_chart" label="统计" active={false} />
        <NavItem icon="person" label="我的" active={false} />
      </nav>
    </div>
  );
};

// 喂养记录项组件
const FeedingLogItem: React.FC<{ record: FeedingRecord }> = ({ record }) => {
  return (
    <div className="flex items-center gap-4 bg-white dark:bg-slate-800/40 p-4 rounded-xl shadow-sm border border-slate-100 dark:border-slate-700/50">
      <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10 text-primary">
        <span className="material-symbols-outlined fill-1">
          baby_changing_station
        </span>
      </div>
      <div className="flex flex-1 flex-col">
        <p className="text-base font-bold text-slate-800 dark:text-slate-100">
          {record.time}
        </p>
        <p className="text-sm text-slate-500 dark:text-slate-400">
          {record.type}
        </p>
      </div>
      <div className="text-right">
        {record.amount !== undefined ? (
          <p className="text-lg font-bold text-primary">
            {record.amount}<span className="text-xs ml-0.5">ml</span>
          </p>
        ) : (
          <p className="text-lg font-bold text-primary">
            {record.duration}<span className="text-xs ml-0.5">分钟</span>
          </p>
        )}
      </div>
    </div>
  );
};

// 底部导航项
const NavItem: React.FC<{
  icon: string;
  label: string;
  active?: boolean;
}> = ({ icon, label, active = false }) => (
  <a
    href="#"
    className={cn(
      'flex flex-col items-center gap-1',
      active ? 'text-primary' : 'text-slate-400 dark:text-slate-500'
    )}
  >
    <span className={cn('material-symbols-outlined text-2xl', active ? 'fill-1' : '')}>
      {icon}
    </span>
    <p className={cn('text-[10px] font-medium', active ? 'font-semibold' : '')}>
      {label}
    </p>
  </a>
);
