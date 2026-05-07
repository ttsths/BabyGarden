import React from 'react';
import styles from './TimelineItem.module.css';

interface TimelineItemProps {
  icon: string;
  title: string;
  description?: string;
  time: string;
  color?: string;
  onClick?: () => void;
}

/**
 * 时间轴记录项
 * 展示单条记录的详细信息
 */
export const TimelineItem: React.FC<TimelineItemProps> = ({
  icon,
  title,
  description,
  time,
  color = 'text-brand-primary',
  onClick,
}) => {
  return (
    <button onClick={onClick} className={styles.timelineItem}>
      {/* 图标区域 */}
      <div className={`${styles.iconWrapper} ${color}`}>
        {icon}
      </div>
      
      {/* 内容区域 */}
      <div className={styles.content}>
        <div className={styles.header}>
          <span className={styles.title}>{title}</span>
          <span className={styles.time}>{time}</span>
        </div>
        {description && (
          <p className={styles.description}>{description}</p>
        )}
      </div>
    </button>
  );
};

export default TimelineItem;
