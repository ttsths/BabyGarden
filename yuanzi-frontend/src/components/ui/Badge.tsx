import React from 'react';

export interface BadgeProps {
  children: React.ReactNode;
  variant?: 'primary' | 'success' | 'warning' | 'error' | 'neutral';
  size?: 'small' | 'medium';
  dot?: boolean;
  className?: string;
}

const variantStyles = {
  primary: 'bg-brand-primary text-white',
  success: 'bg-success text-white',
  warning: 'bg-warning text-neutral-text-primary',
  error: 'bg-error text-white',
  neutral: 'bg-neutral-border text-neutral-text-secondary',
};

const sizeStyles = {
  small: 'px-2 py-0.5 text-xs',
  medium: 'px-3 py-1 text-sm',
};

/**
 * 徽章组件
 */
export const Badge: React.FC<BadgeProps> = ({
  children,
  variant = 'neutral',
  size = 'small',
  dot = false,
  className = '',
}) => {
  if (dot) {
    return (
      <span className={`inline-flex items-center gap-1 ${className}`}>
        <span className={`w-2 h-2 rounded-full ${variantStyles[variant].split(' ')[0]}`} />
        {children}
      </span>
    );
  }

  return (
    <span
      className={`inline-flex items-center justify-center font-medium rounded-full ${variantStyles[variant]} ${sizeStyles[size]} ${className}`}
    >
      {children}
    </span>
  );
};

export default Badge;
