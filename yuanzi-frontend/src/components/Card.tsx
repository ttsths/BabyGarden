import React from 'react';
import { cn } from '../utils/cn';

interface CardProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
  padding?: 'none' | 'sm' | 'md' | 'lg';
  hoverable?: boolean;
}

export const Card: React.FC<CardProps> = ({
  children,
  className,
  onClick,
  padding = 'md',
  hoverable = false,
}) => {
  const paddingStyles = {
    none: '',
    sm: 'p-sm',
    md: 'p-md',
    lg: 'p-lg',
  };

  return (
    <div
      className={cn(
        'bg-white border border-border rounded-lg shadow-md',
        paddingStyles[padding],
        hoverable && 'hover:shadow-lg cursor-pointer transition-shadow duration-300',
        onClick && 'cursor-pointer',
        className
      )}
      onClick={onClick}
    >
      {children}
    </div>
  );
};
