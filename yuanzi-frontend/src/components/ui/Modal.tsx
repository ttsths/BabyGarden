import React, { useEffect } from 'react';
import { Card } from './Card';
import { Button } from './Button';

export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
  footer?: React.ReactNode;
  showClose?: boolean;
}

/**
 * 模态框组件
 */
export const Modal: React.FC<ModalProps> = ({
  isOpen,
  onClose,
  title,
  children,
  footer,
  showClose = true,
}) => {
  // 阻止背景滚动
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }

    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  // ESC 关闭
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEsc);
    return () => document.removeEventListener('keydown', handleEsc);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
    >
      {/* 背景遮罩 */}
      <div
        className="absolute inset-0 bg-black/50 animate-fade-in"
        onClick={onClose}
      />

      {/* 模态框内容 */}
      <Card
        padding="large"
        className="relative w-full max-w-sm animate-slide-up bg-background-secondary"
      >
        {/* 关闭按钮 */}
        {showClose && (
          <button
            onClick={onClose}
            className="absolute top-4 right-4 text-neutral-text-secondary hover:text-neutral-text-primary"
            aria-label="关闭"
          >
            ✕
          </button>
        )}

        {/* 标题 */}
        {title && (
          <h2 className="text-xl font-semibold text-neutral-text-primary mb-4 pr-8">
            {title}
          </h2>
        )}

        {/* 内容 */}
        <div className="text-neutral-text-primary">
          {children}
        </div>

        {/* 底部操作区 */}
        {footer && (
          <div className="mt-6 flex gap-3">
            {footer}
          </div>
        )}
      </Card>
    </div>
  );
};

export default Modal;
