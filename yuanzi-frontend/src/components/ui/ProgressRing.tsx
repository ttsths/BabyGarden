import React from 'react';
import styles from './ProgressRing.module.css';

interface ProgressRingProps {
  percentage: number;
  size?: number;
  strokeWidth?: number;
  children?: React.ReactNode;
}

/**
 * 环形进度条组件
 * 用于展示今日完成度
 */
export const ProgressRing: React.FC<ProgressRingProps> = ({
  percentage,
  size = 200,
  strokeWidth = 12,
  children,
}) => {
  const radius = (size - strokeWidth) / 2;
  const circumference = radius * 2 * Math.PI;
  const offset = circumference - (percentage / 100) * circumference;

  return (
    <div className={styles.progressRing} style={{ width: size, height: size }}>
      <svg width={size} height={size}>
        {/* 背景圆环 */}
        <circle
          className="text-gray-200"
          strokeWidth={strokeWidth}
          stroke="#F0E6DE"
          fill="transparent"
          r={radius}
          cx={size / 2}
          cy={size / 2}
        />
        {/* 进度圆环 */}
        <circle
          strokeWidth={strokeWidth}
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          strokeLinecap="round"
          stroke="var(--color-brand)"
          fill="transparent"
          r={radius}
          cx={size / 2}
          cy={size / 2}
          className="transition-all duration-500 ease-in-out"
        />
      </svg>
      {/* 中心内容 */}
      <div className="absolute inset-0 flex items-center justify-center">
        {children}
      </div>
    </div>
  );
};

export default ProgressRing;
