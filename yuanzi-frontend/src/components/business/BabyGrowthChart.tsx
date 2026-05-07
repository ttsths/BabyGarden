import React from 'react';
import { Card } from '@/components/ui/Card';

interface BabyGrowthChartProps {
  data: {
    weight: { current: number; unit: string; change: string };
    height: { current: number; unit: string; change: string };
    head: { current: number; unit: string; change: string };
  };
  className?: string;
}

/**
 * 宝宝成长曲线组件
 * 展示体重、身长、头围数据
 */
export const BabyGrowthChart: React.FC<BabyGrowthChartProps> = ({
  data,
  className = '',
}) => {
  return (
    <Card padding="large" className={className}>
      <h3 className="text-lg font-semibold text-neutral-text-primary mb-6">
        成长数据
      </h3>
      
      <div className="grid grid-cols-3 gap-4">
        {/* 体重 */}
        <div className="text-center p-4 bg-background-tertiary rounded-lg">
          <div className="text-3xl mb-1">⚖️</div>
          <div className="text-2xl font-bold text-neutral-text-primary">
            {data.weight.current}
          </div>
          <div className="text-xs text-neutral-text-secondary mt-1">
            体重 ({data.weight.unit})
          </div>
          <div className="text-xs text-success mt-1 font-medium">
            {data.weight.change}
          </div>
        </div>

        {/* 身长 */}
        <div className="text-center p-4 bg-background-tertiary rounded-lg">
          <div className="text-3xl mb-1">📏</div>
          <div className="text-2xl font-bold text-neutral-text-primary">
            {data.height.current}
          </div>
          <div className="text-xs text-neutral-text-secondary mt-1">
            身长 ({data.height.unit})
          </div>
          <div className="text-xs text-success mt-1 font-medium">
            {data.height.change}
          </div>
        </div>

        {/* 头围 */}
        <div className="text-center p-4 bg-background-tertiary rounded-lg">
          <div className="text-3xl mb-1">🧢</div>
          <div className="text-2xl font-bold text-neutral-text-primary">
            {data.head.current}
          </div>
          <div className="text-xs text-neutral-text-secondary mt-1">
            头围 ({data.head.unit})
          </div>
          <div className="text-xs text-success mt-1 font-medium">
            {data.head.change}
          </div>
        </div>
      </div>

      {/* 生长曲线提示 */}
      <div className="mt-6 p-4 bg-brand-light/10 rounded-lg">
        <p className="text-sm text-neutral-text-secondary">
          💡 宝宝的生长曲线在正常范围内，继续保持！
        </p>
      </div>
    </Card>
  );
};

export default BabyGrowthChart;
