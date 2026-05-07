import React from 'react';

export interface LoadingProps {
  size?: 'small' | 'medium' | 'large';
  text?: string;
  fullScreen?: boolean;
  className?: string;
}

/**
 * Loading 加载组件
 * 用于数据加载状态
 */
export const Loading: React.FC<LoadingProps> = ({
  size = 'medium',
  text,
  fullScreen = false,
  className = '',
}) => {
  const sizeStyles = {
    small: 'w-6 h-6',
    medium: 'w-10 h-10',
    large: 'w-16 h-16',
  };

  const content = (
    <div className={`flex flex-col items-center justify-center ${className}`}>
      {/* 加载动画 */}
      <div
        className={`${sizeStyles[size]} border-4 border-brand-light border-t-brand-primary rounded-full animate-spin`}
      />
      
      {/* 加载文字 */}
      {text && (
        <p className="text-body text-neutral-text-secondary mt-md">
          {text}
        </p>
      )}
    </div>
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 bg-background-primary/80 flex items-center justify-center z-50">
        {content}
      </div>
    );
  }

  return content;
};

export default Loading;
