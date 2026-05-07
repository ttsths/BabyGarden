import React from 'react';

export interface SwitchProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
  size?: 'small' | 'medium';
  className?: string;
}

/**
 * Switch 开关组件
 * 用于设置项的开关控制
 */
export const Switch: React.FC<SwitchProps> = ({
  checked,
  onChange,
  disabled = false,
  size = 'medium',
  className = '',
}) => {
  const sizeStyles = {
    small: 'w-10 h-6',
    medium: 'w-12 h-7',
  };

  const thumbStyles = {
    small: checked ? 'translate-x-5' : 'translate-x-1',
    medium: checked ? 'translate-x-6' : 'translate-x-1',
  };

  return (
    <button
      role="switch"
      aria-checked={checked}
      disabled={disabled}
      onClick={() => !disabled && onChange(!checked)}
      className={`
        ${sizeStyles[size]}
        ${checked ? 'bg-brand-primary' : 'bg-neutral-border'}
        ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
        relative rounded-full transition-colors duration-200
        ${className}
      `}
    >
      <div
        className={`
          ${thumbStyles[size]}
          w-5 h-5 bg-white rounded-full shadow-md
          transform transition-transform duration-200
        `}
      />
    </button>
  );
};

export default Switch;
