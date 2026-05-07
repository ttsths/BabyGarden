import React from 'react';
import { babyInfo, todayProgress, todayStats, recentTimeline } from '@/data/mockData';
import { StatsCard } from './StatsCard';
import { QuickStatCard } from './QuickStatCard';
import { TimelineItem } from './TimelineItem';
import { BottomNav } from './BottomNav';

interface HomePageProps {
  // 预留自定义配置接口
  className?: string;
}

/**
 * 园子宝宝 App - 首页
 * 展示今日概览、统计数据和最近记录
 */
export const HomePage: React.FC<HomePageProps> = ({ className = '' }) => {
  return (
    <div className={`min-h-screen bg-background-primary pb-20 ${className}`}>
      {/* 顶部宝宝信息区域 */}
      <header className="bg-gradient-to-r from-brand-light to-brand-primary px-lg pt-12 pb-xl">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-md">
            <div className="w-12 h-12 rounded-full bg-white/30 flex items-center justify-center text-2xl">
              👶
            </div>
            <div>
              <h1 className="text-white text-h1 font-semibold">{babyInfo.name}</h1>
              <p className="text-white/80 text-caption">{babyInfo.age}</p>
            </div>
          </div>
          <button className="w-10 h-10 rounded-full bg-white/20 flex items-center justify-center text-white hover:bg-white/30 transition-colors">
            🔔
          </button>
        </div>
      </header>

      {/* 主要内容区域 */}
      <main className="px-lg -mt-4 space-y-lg">
        {/* 今日完成度统计卡片 */}
        <StatsCard
          completed={todayProgress.completed}
          total={todayProgress.total}
          percentage={todayProgress.percentage}
          breakdown={todayProgress.breakdown}
        />

        {/* 今日概览快捷统计 */}
        <section>
          <h2 className="text-h2 font-medium text-neutral-text-primary mb-md px-sm">
            今日概览
          </h2>
          <div className="grid grid-cols-2 gap-md">
            {todayStats.map((stat) => (
              <QuickStatCard
                key={stat.id}
                icon={stat.icon}
                label={stat.label}
                value={stat.value}
                bgColor={stat.bgColor}
                iconBg={stat.iconBg}
              />
            ))}
          </div>
        </section>

        {/* 最近记录时间轴 */}
        <section>
          <div className="flex items-center justify-between mb-md px-sm">
            <h2 className="text-h2 font-medium text-neutral-text-primary">
              最近记录
            </h2>
            <button className="text-body text-brand-primary hover:text-brand-dark transition-colors">
              查看全部
            </button>
          </div>
          <div className="space-y-md">
            {recentTimeline.map((record) => (
              <TimelineItem
                key={record.id}
                icon={record.icon}
                title={record.title}
                description={record.description}
                time={record.time}
                color={record.color}
              />
            ))}
          </div>
        </section>
      </main>

      {/* 底部导航栏 */}
      <BottomNav />
    </div>
  );
};

export default HomePage;
