import React, { useEffect, useState } from 'react';

export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface ToastProps {
  message: string;
  type?: ToastType;
  duration?: number;
  onClose?: () => void;
}

const toastStyles: Record<ToastType, string> = {
  success: 'bg-success',
  error: 'bg-error',
  warning: 'bg-warning text-neutral-text-primary',
  info: 'bg-info',
};

/**
 * Toast 提示组件
 */
export const Toast: React.FC<ToastProps> = ({
  message,
  type = 'info',
  duration = 2000,
  onClose,
}) => {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false);
      onClose?.();
    }, duration);

    return () => clearTimeout(timer);
  }, [duration, onClose]);

  if (!isVisible) return null;

  return (
    <div
      className={`fixed bottom-24 left-1/2 -translate-x-1/2 px-6 py-3 rounded-full text-white text-sm font-medium shadow-lg animate-slide-up z-50 ${toastStyles[type]}`}
      role="alert"
    >
      {message}
    </div>
  );
};

// Toast 管理器
let toastContainer: HTMLElement | null = null;

// eslint-disable-next-line react-refresh/only-export-components
export const showToast = (message: string, type: ToastType = 'info', duration?: number) => {
  if (!toastContainer) {
    toastContainer = document.createElement('div');
    toastContainer.className = 'fixed bottom-24 left-1/2 -translate-x-1/2 z-50';
    document.body.appendChild(toastContainer);
  }

  const toast = document.createElement('div');
  toast.className = `px-6 py-3 rounded-full text-white text-sm font-medium shadow-lg animate-slide-up ${toastStyles[type]}`;
  toast.textContent = message;
  
  toastContainer.appendChild(toast);

  setTimeout(() => {
    toast.classList.add('animate-fade-out');
    setTimeout(() => {
      toast.remove();
      if (toastContainer && toastContainer.children.length === 0) {
        toastContainer.remove();
        toastContainer = null;
      }
    }, 200);
  }, duration || 2000);
};

export default Toast;
