import React from 'react';
import { cn } from '../utils/cn';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
  size?: 'small' | 'medium' | 'large' | 'xl';
  loading?: boolean;
  disabled?: boolean;
  block?: boolean;
  icon?: React.ReactNode;
  iconPosition?: 'left' | 'right';
  children: React.ReactNode;
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'medium',
  loading = false,
  disabled = false,
  block = false,
  icon,
  iconPosition = 'left',
  children,
  className,
  ...props
}) => {
  const baseStyles = 'inline-flex items-center justify-center font-medium transition-all duration-150 ease-in-out disabled:opacity-50 disabled:cursor-not-allowed';

  const variantStyles = {
    primary: 'bg-brand hover:bg-brand-dark active:bg-brand-dark text-white border border-brand',
    secondary: 'bg-white hover:bg-bg-tertiary text-brand border border-brand',
    ghost: 'bg-transparent hover:bg-brand/10 text-brand border border-transparent',
    danger: 'bg-error hover:bg-red-600 text-white border border-error',
  };

  const sizeStyles = {
    small: 'h-8 px-3 text-caption',
    medium: 'h-10 px-4 text-body',
    large: 'h-12 px-5 text-h3',
    xl: 'h-14 px-6 text-h2',
  };

  return (
    <button
      className={cn(
        baseStyles,
        variantStyles[variant],
        sizeStyles[size],
        block && 'w-full',
        'rounded-md',
        className
      )}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? (
        <span className="animate-spin mr-2">⟳</span>
      ) : (
        <>
          {icon && iconPosition === 'left' && <span className="mr-2">{icon}</span>}
          {children}
          {icon && iconPosition === 'right' && <span className="ml-2">{icon}</span>}
        </>
      )}
    </button>
  );
};
