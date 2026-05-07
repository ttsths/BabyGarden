import React from 'react';

export interface AvatarProps {
  src?: string;
  alt?: string;
  size?: 'small' | 'medium' | 'large' | 'xlarge';
  fallback?: string;
  className?: string;
}

const sizeClasses = {
  small: 'w-8 h-8 text-sm',
  medium: 'w-10 h-10 text-base',
  large: 'w-12 h-12 text-lg',
  xlarge: 'w-16 h-16 text-xl',
};

/**
 * 头像组件
 */
export const Avatar: React.FC<AvatarProps> = ({
  src,
  alt = '头像',
  size = 'medium',
  fallback = '👤',
  className = '',
}) => {
  const [hasError, setHasError] = React.useState(false);

  if (hasError || !src) {
    return (
      <div
        className={`${sizeClasses[size]} rounded-full bg-brand-light/20 flex items-center justify-center ${className}`}
        aria-label={alt}
      >
        {fallback}
      </div>
    );
  }

  return (
    <img
      src={src}
      alt={alt}
      className={`${sizeClasses[size]} rounded-full object-cover ${className}`}
      onError={() => setHasError(true)}
    />
  );
};

export default Avatar;
