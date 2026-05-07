import React from 'react';
import styles from './StatsCard.module.css';

interface StatsCardProps {
  completed: number;
  total: number;
  percentage: number;
  breakdown?: Array<{
    category: string;
    completed: number;
    total: number;
    percentage: number;
    unit?: string;
  }>;
}

/**
 * 今日完成度统计卡片
 * 展示整体进度和各项细分进度
 */
export const StatsCard: React.FC<StatsCardProps> = ({
  completed,
  total,
  percentage,
  breakdown = [],
}) => {
  const circumference = 2 * Math.PI * 88;
  const strokeDashoffset = circumference * (1 - percentage / 100);

  return (
    <div className={styles.statsCard}>
      {/* 环形进度图 */}
      <div className={styles.progressContainer}>
        <div className={styles.progressRing}>
          <svg width="200" height="200" className="transform -rotate-90">
            {/* 背景圆环 */}
            <circle
              cx="100"
              cy="100"
              r="88"
              fill="none"
              stroke="#F0E6DE"
              strokeWidth="12"
            />
            {/* 进度圆环 */}
            <circle
              cx="100"
              cy="100"
              r="88"
              fill="none"
              stroke="var(--color-brand)"
              strokeWidth="12"
              strokeLinecap="round"
              strokeDasharray={circumference}
              strokeDashoffset={strokeDashoffset}
              className="transition-all duration-500 ease-in-out"
            />
          </svg>
          {/* 中心文字 */}
          <div className={styles.progressText}>
            <span className={styles.percentage}>{percentage}%</span>
            <span className={styles.label}>已完成</span>
          </div>
        </div>
      </div>

      {/* 进度说明 */}
      <div className="text-center mb-xl">
        <h3 className="text-h2 font-medium text-neutral-text-primary mb-sm">
          今日记录完整度
        </h3>
        <p className="text-body text-neutral-text-secondary">
          已完成 {completed} 项 · 共 {total} 项
        </p>
      </div>

      {/* 细项进度 */}
      {breakdown.length > 0 && (
        <div className={styles.breakdown}>
          {breakdown.map((item, index) => (
            <div key={index} className={styles.breakdownItem}>
              <span className={styles.breakdownLabel}>
                {item.category}
              </span>
              <div className={styles.breakdownProgress}>
                <div
                  className={styles.breakdownBar}
                  style={{ width: `${item.percentage}%` }}
                />
              </div>
              <span className={styles.breakdownValue}>
                {item.percentage}%
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default StatsCard;
