import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { weeklyStats } from '@/data/mockData';
import { Card } from '@/components/ui/Card';
import { BottomNav } from '@/components/ui/BottomNav';

/**
 * 统计页面
 * 展示宝宝成长数据和趋势分析
 */
export const StatsPage: React.FC = () => {
  const navigate = useNavigate();
  const [timeRange, setTimeRange] = useState<'day' | 'week' | 'month' | 'year'>('week');
  const [statType, setStatType] = useState<'feeding' | 'sleep' | 'diaper' | 'growth'>('feeding');

  // 模拟统计数据
  const stats = {
    feeding: {
      title: '喂养统计',
      icon: '🍼',
      data: {
        average: 6.2,
        unit: '次/天',
        trend: [5, 7, 6, 8, 5, 6, 7],
        labels: ['一', '二', '三', '四', '五', '六', '日'],
      },
    },
    sleep: {
      title: '睡眠统计',
      icon: '💤',
      data: {
        average: 13.5,
        unit: '小时/天',
        trend: [12, 14, 13, 15, 13, 14, 13],
        labels: ['一', '二', '三', '四', '五', '六', '日'],
      },
    },
    diaper: {
      title: '排泄统计',
      icon: '💩',
      data: {
        average: 4.8,
        unit: '次/天',
        trend: [4, 5, 6, 4, 5, 4, 5],
        labels: ['一', '二', '三', '四', '五', '六', '日'],
      },
    },
    growth: {
      title: '成长统计',
      icon: '📈',
      data: {
        weight: { current: 6.8, unit: 'kg', change: '+0.5' },
        height: { current: 62.5, unit: 'cm', change: '+2.1' },
        head: { current: 40.2, unit: 'cm', change: '+0.8' },
      },
    },
  };

  const renderChart = (data: any) => {
    const maxValue = Math.max(...data.trend);
    
    return (
      <div className="h-64 flex items-end justify-between gap-2 pt-8">
        {data.trend.map((value: number, index: number) => {
          const height = (value / maxValue) * 100;
          return (
            <div key={index} className="flex-1 flex flex-col items-center gap-2">
              <span className="text-xs text-neutral-text-secondary font-medium">
                {value}
              </span>
              <div
                className="w-full bg-gradient-to-t from-brand-light to-brand-primary rounded-t-lg transition-all"
                style={{ height: `${height}%` }}
              />
              <span className="text-xs text-neutral-text-secondary">{data.labels[index]}</span>
            </div>
          );
        })}
      </div>
    );
  };

  const renderGrowthStats = () => {
    const growth = stats.growth.data;
    return (
      <div className="space-y-4">
        <div className="grid grid-cols-3 gap-4">
          <div className="text-center p-4 bg-background-tertiary rounded-lg">
            <div className="text-2xl font-bold text-neutral-text-primary">
              {growth.weight.current}
            </div>
            <div className="text-xs text-neutral-text-secondary mt-1">
              体重 (kg)
            </div>
            <div className="text-xs text-success mt-1">
              {growth.weight.change}
            </div>
          </div>
          <div className="text-center p-4 bg-background-tertiary rounded-lg">
            <div className="text-2xl font-bold text-neutral-text-primary">
              {growth.height.current}
            </div>
            <div className="text-xs text-neutral-text-secondary mt-1">
              身长 (cm)
            </div>
            <div className="text-xs text-success mt-1">
              {growth.height.change}
            </div>
          </div>
          <div className="text-center p-4 bg-background-tertiary rounded-lg">
            <div className="text-2xl font-bold text-neutral-text-primary">
              {growth.head.current}
            </div>
            <div className="text-xs text-neutral-text-secondary mt-1">
              头围 (cm)
            </div>
            <div className="text-xs text-success mt-1">
              {growth.head.change}
            </div>
          </div>
        </div>
      </div>
    );
  };

  const currentStat = stats[statType];
  const currentData = currentStat.data as { average: number; unit: string; trend: number[]; labels: string[] };

  return (
    <div className="min-h-screen bg-background-primary pb-20">
      <header className="px-6 pt-12 pb-4 bg-background-secondary sticky top-0 z-10">
        <div className="flex items-center">
          <button
            onClick={() => navigate(-1)}
            className="text-neutral-text-primary mr-4"
          >
            ← 返回
          </button>
          <h1 className="text-xl font-semibold text-neutral-text-primary">
            统计
          </h1>
        </div>
      </header>

      <main className="px-4 py-6 space-y-6">
        {/* 统计类型选择 */}
        <Card padding="medium">
          <div className="grid grid-cols-4 gap-2">
            {(Object.keys(stats) as Array<keyof typeof stats>).map((type) => (
              <button
                key={type}
                onClick={() => setStatType(type)}
                className={`flex flex-col items-center gap-2 p-3 rounded-lg transition-colors ${
                  statType === type
                    ? 'bg-brand-primary text-white'
                    : 'bg-background-tertiary text-neutral-text-secondary'
                }`}
              >
                <span className="text-2xl">{stats[type].icon}</span>
                <span className="text-xs font-medium">
                  {stats[type].title.split('统计')[0]}
                </span>
              </button>
            ))}
          </div>
        </Card>

        {/* 时间范围选择 */}
        {statType !== 'growth' && (
          <div className="flex gap-2">
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
        )}

        {/* 统计图表/数据 */}
        <Card padding="large">
          <h3 className="text-lg font-semibold text-neutral-text-primary mb-2">
            {currentStat.title}
          </h3>
          {statType === 'growth' ? (
            renderGrowthStats()
          ) : (
            <>
              <div className="mb-6">
                <span className="text-3xl font-bold text-neutral-text-primary">
                  {currentData.average}
                </span>
                <span className="text-neutral-text-secondary ml-2">
                  {currentData.unit}
                </span>
              </div>
              {renderChart(currentData)}
            </>
          )}
        </Card>

        {/* 统计摘要 */}
        <Card padding="medium">
          <h3 className="text-base font-semibold text-neutral-text-primary mb-4">
            本周摘要
          </h3>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-neutral-text-secondary">总记录数</span>
              <span className="text-neutral-text-primary font-medium">
                {weeklyStats.totalRecords} 次
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-neutral-text-secondary">完整记录天数</span>
              <span className="text-neutral-text-primary font-medium">
                {weeklyStats.completeDays} 天
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-neutral-text-secondary">记录完整度</span>
              <span className="text-brand-primary font-medium">
                85%
              </span>
            </div>
          </div>
        </Card>
      </main>

      <BottomNav />
    </div>
  );
};

export default StatsPage;
