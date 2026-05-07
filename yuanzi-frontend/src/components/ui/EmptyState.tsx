import React from 'react';

export interface EmptyStateProps {
  icon?: string;
  title: string;
  description?: string;
  action?: {
    label: string;
    onClick: () => void;
  };
  className?: string;
}

/**
 * EmptyState 空状态组件
 * 用于无数据、无网络等场景
 */
export const EmptyState: React.FC<EmptyStateProps> = ({
  icon = '📭',
  title,
  description,
  action,
  className = '',
}) => {
  return (
    <div className={`flex flex-col items-center justify-center p-xl text-center ${className}`}>
      {/* 图标 */}
      <div className="text-6xl mb-lg">{icon}</div>
      
      {/* 标题 */}
      <h3 className="text-h2 font-medium text-neutral-text-primary mb-sm">
        {title}
      </h3>
      
      {/* 描述 */}
      {description && (
        <p className="text-body text-neutral-text-secondary mb-lg max-w-xs">
          {description}
        </p>
      )}
      
      {/* 操作按钮 */}
      {action && (
        <button
          onClick={action.onClick}
          className="px-lg py-sm bg-brand-primary text-white rounded-lg font-medium hover:bg-brand-dark transition-colors"
        >
          {action.label}
        </button>
      )}
    </div>
  );
};

export default EmptyState;
