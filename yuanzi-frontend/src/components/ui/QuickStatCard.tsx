import React from 'react';
import styles from './QuickStatCard.module.css';

interface QuickStatCardProps {
  icon: string;
  label: string;
  value: string;
  bgColor?: string;
  iconBg?: string;
  onClick?: () => void;
}

/**
 * 快捷统计卡片
 * 展示单项统计数据（喂奶、睡眠、排泄、体温等）
 */
export const QuickStatCard: React.FC<QuickStatCardProps> = ({
  icon,
  label,
  value,
  bgColor = 'bg-background-tertiary',
  iconBg = 'bg-brand-light/20',
  onClick,
}) => {
  return (
    <button
      onClick={onClick}
      className={`${styles.quickStatCard} ${bgColor}`}
    >
      {/* 图标区域 */}
      <div className={`${styles.iconWrapper} ${iconBg}`}>
        {icon}
      </div>
      
      {/* 数据区域 */}
      <span className={styles.value}>{value}</span>
      <span className={styles.label}>{label}</span>
    </button>
  );
};

export default QuickStatCard;
