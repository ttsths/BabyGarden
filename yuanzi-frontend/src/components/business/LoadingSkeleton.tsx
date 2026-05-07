import React from 'react';

interface LoadingSkeletonProps {
  type?: 'card' | 'list' | 'text' | 'image';
  count?: number;
  className?: string;
}

/**
 * 加载骨架屏组件
 */
export const LoadingSkeleton: React.FC<LoadingSkeletonProps> = ({
  type = 'card',
  count = 3,
  className = '',
}) => {
  const renderSkeleton = (index: number) => {
    switch (type) {
      case 'card':
        return (
          <div key={index} className="p-4 rounded-lg bg-background-secondary animate-pulse">
            <div className="h-4 bg-neutral-border rounded w-3/4 mb-3" />
            <div className="h-3 bg-neutral-border rounded w-1/2 mb-2" />
            <div className="h-3 bg-neutral-border rounded w-2/3" />
          </div>
        );

      case 'list':
        return (
          <div key={index} className="flex items-center gap-4 p-4 animate-pulse">
            <div className="w-12 h-12 rounded-full bg-neutral-border" />
            <div className="flex-1">
              <div className="h-4 bg-neutral-border rounded w-3/4 mb-2" />
              <div className="h-3 bg-neutral-border rounded w-1/2" />
            </div>
          </div>
        );

      case 'text':
        return (
          <div key={index} className="space-y-2 animate-pulse">
            <div className="h-4 bg-neutral-border rounded w-full" />
            <div className="h-4 bg-neutral-border rounded w-5/6" />
            <div className="h-4 bg-neutral-border rounded w-4/6" />
          </div>
        );

      case 'image':
        return (
          <div
            key={index}
            className="bg-neutral-border rounded-lg animate-pulse"
            style={{ aspectRatio: '1' }}
          />
        );

      default:
        return null;
    }
  };

  return (
    <div className={`space-y-4 ${className}`}>
      {Array.from({ length: count }).map((_, index) => renderSkeleton(index))}
    </div>
  );
};

export default LoadingSkeleton;
