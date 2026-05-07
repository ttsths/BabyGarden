import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { recentTimeline, weeklyStats } from '@/data/mockData';
import { Card } from '@/components/ui/Card';
import { TimelineItem } from '@/components/ui/TimelineItem';
import { BottomNav } from '@/components/ui/BottomNav';

/**
 * 时间轴页面
 * 展示所有记录的时间线，支持按日分组
 */
export const TimelinePage: React.FC = () => {
  const navigate = useNavigate();
  const [viewMode, setViewMode] = useState<'timeline' | 'stats'>('timeline');
  const [timeRange, setTimeRange] = useState<'day' | 'week' | 'month' | 'year'>('week');

  // 模拟更多数据
  const timelineData = [
    {
      date: '今天 · 3 月 6 日 周五',
      count: 8,
      records: recentTimeline,
    },
    {
      date: '昨天 · 3 月 5 日 周四',
      count: 10,
      records: [
        {
          id: '4',
          type: 'sleep',
          icon: '💤',
          title: '夜间睡眠',
          description: '睡了 11 小时',
          time: '20:00',
          color: 'text-purple-500',
        },
        {
          id: '5',
          type: 'feeding',
          icon: '🍼',
          title: '喂奶',
          description: '150ml',
          time: '18:30',
          color: 'text-orange-500',
        },
      ],
    },
  ];

  // 统计视图
  const renderStatsView = () => {
    return (
      <div className="space-y-4">
        {/* 时间范围选择 */}
        <div className="flex gap-2 mb-4">
          {(['day', 'week', 'month', 'year'] as const).map((range) => (
            <button
              key={range}
              onClick={() => setTimeRange(range)}
              className={`flex-1 py-2 rounded-lg text-sm font-medium transition-colors ${
                timeRange === range
                  ? 'bg-brand-primary text-white'
                  : 'bg-background-secondary text-neutral-text-secondary'
              }`}
            >
              {range === 'day' && '日'}
              {range === 'week' && '周'}
              {range === 'month' && '月'}
              {range === 'year' && '年'}
            </button>
          ))}
        </div>

        {/* 趋势图表 */}
        <Card padding="medium">
          <h3 className="text-lg font-semibold text-neutral-text-primary mb-4">
            本周喂养趋势
          </h3>
          <div className="h-48 flex items-end justify-between gap-2">
            {weeklyStats.feeding.trend.map((value, index) => {
              const height = (value / 10) * 100;
              const days = ['一', '二', '三', '四', '五', '六', '日'];
              return (
                <div key={index} className="flex-1 flex flex-col items-center gap-2">
                  <div
                    className="w-full bg-gradient-to-t from-brand-light to-brand-primary rounded-t-lg transition-all"
                    style={{ height: `${height}%` }}
                  />
                  <span className="text-xs text-neutral-text-secondary">{days[index]}</span>
                </div>
              );
            })}
          </div>
        </Card>

        {/* 统计数据 */}
        <Card padding="medium">
          <h3 className="text-lg font-semibold text-neutral-text-primary mb-4">
            本周平均
          </h3>
          <div className="grid grid-cols-2 gap-4">
            <div className="p-4 bg-background-tertiary rounded-lg">
              <div className="text-2xl font-bold text-neutral-text-primary">
                {weeklyStats.feeding.average}
              </div>
              <div className="text-sm text-neutral-text-secondary mt-1">
                日均喂奶 ({weeklyStats.feeding.unit})
              </div>
            </div>
            <div className="p-4 bg-background-tertiary rounded-lg">
              <div className="text-2xl font-bold text-neutral-text-primary">
                {weeklyStats.sleep.average}
              </div>
              <div className="text-sm text-neutral-text-secondary mt-1">
                日均睡眠 ({weeklyStats.sleep.unit})
              </div>
            </div>
            <div className="p-4 bg-background-tertiary rounded-lg">
              <div className="text-2xl font-bold text-neutral-text-primary">
                {weeklyStats.totalRecords}
              </div>
              <div className="text-sm text-neutral-text-secondary mt-1">
                总记录数
              </div>
            </div>
            <div className="p-4 bg-background-tertiary rounded-lg">
              <div className="text-2xl font-bold text-neutral-text-primary">
                {weeklyStats.completeDays}
              </div>
              <div className="text-sm text-neutral-text-secondary mt-1">
                完整天数
              </div>
            </div>
          </div>
        </Card>
      </div>
    );
  };

  // 时间轴视图
  const renderTimelineView = () => {
    return (
      <div className="space-y-6">
        {timelineData.map((day, dayIndex) => (
          <div key={dayIndex}>
            <div className="flex items-center justify-between mb-3 px-1">
              <h3 className="text-base font-medium text-neutral-text-primary">
                {day.date}
              </h3>
              <span className="text-sm text-neutral-text-secondary">
                共 {day.count} 条记录
              </span>
            </div>
            <Card padding="medium">
              <div className="space-y-3">
                {day.records.map((record, recordIndex) => (
                  <div
                    key={record.id}
                    className={recordIndex !== day.records.length - 1 ? 'pb-3 border-b border-neutral-border last:border-0' : ''}
                  >
                    <TimelineItem
                      icon={record.icon}
                      title={record.title}
                      description={record.description}
                      time={record.time}
                      color={record.color}
                    />
                  </div>
                ))}
              </div>
            </Card>
          </div>
        ))}
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-background-primary pb-20">
      <header className="px-6 pt-12 pb-4 bg-background-secondary sticky top-0 z-10">
        <div className="flex items-center justify-between">
          <button
            onClick={() => navigate(-1)}
            className="text-neutral-text-primary"
          >
            ← 返回
          </button>
          <h1 className="text-xl font-semibold text-neutral-text-primary">
            时间轴
          </h1>
          <button
            onClick={() => setViewMode(viewMode === 'timeline' ? 'stats' : 'timeline')}
            className="text-brand-primary font-medium"
          >
            {viewMode === 'timeline' ? '📊 统计' : '📅 时间轴'}
          </button>
        </div>
      </header>

      <main className="px-4 py-6">
        {viewMode === 'timeline' ? renderTimelineView() : renderStatsView()}
      </main>

      <BottomNav />
    </div>
  );
};

export default TimelinePage;
