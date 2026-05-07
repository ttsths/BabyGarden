import { forwardRef } from 'react';
import type { InputHTMLAttributes } from 'react';
import styles from './Input.module.css';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  type?: 'text' | 'number' | 'password' | 'tel' | 'email' | 'date';
  label?: string;
  error?: string;
  fullWidth?: boolean;
  inputSize?: 'small' | 'medium' | 'large';
}

export const Input = forwardRef<HTMLInputElement, InputProps>(({
  type = 'text',
  label,
  error,
  fullWidth = false,
  inputSize = 'medium',
  className = '',
  ...props
}, ref) => {
  const wrapperClasses = [
    styles.wrapper,
    fullWidth ? styles.wrapperFullWidth : '',
    className
  ].filter(Boolean).join(' ');

  const inputClasses = [
    styles.input,
    styles['input-' + inputSize],
    error ? styles.inputError : ''
  ].filter(Boolean).join(' ');

  return (
    <div className={wrapperClasses}>
      {label && (
        <label className={styles.label}>{label}</label>
      )}
      <input
        ref={ref}
        type={type}
        className={inputClasses}
        {...props}
      />
      {error && (
        <span className={styles.error}>{error}</span>
      )}
    </div>
  );
});

Input.displayName = 'Input';
