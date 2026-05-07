import React from 'react';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';

interface EmptyStateProps {
  icon: string;
  title: string;
  description?: string;
  actionText?: string;
  onAction?: () => void;
  className?: string;
}

/**
 * 空状态组件
 * 用于无数据、无网络等场景
 */
export const EmptyState: React.FC<EmptyStateProps> = ({
  icon,
  title,
  description,
  actionText,
  onAction,
  className = '',
}) => {
  return (
    <Card padding="large" className={className}>
      <div className="text-center py-12">
        <div className="text-6xl mb-4">{icon}</div>
        <h3 className="text-lg font-medium text-neutral-text-primary mb-2">
          {title}
        </h3>
        {description && (
          <p className="text-neutral-text-secondary mb-6">
            {description}
          </p>
        )}
        {actionText && onAction && (
          <Button
            variant="primary"
            size="medium"
            onClick={onAction}
          >
            {actionText}
          </Button>
        )}
      </div>
    </Card>
  );
};

export default EmptyState;
