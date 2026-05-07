import type { FC, ReactNode } from 'react';
import styles from './Card.module.css';

export interface CardProps {
  children: ReactNode;
  className?: string;
  onClick?: () => void;
  padding?: 'none' | 'small' | 'medium' | 'large';
}

export const Card: FC<CardProps> = ({
  children,
  className = '',
  onClick,
  padding = 'medium',
}) => {
  return (
    <div
      className={`
        ${styles.card}
        ${styles[`card--${padding}`]}
        ${className}
      `}
      onClick={onClick}
    >
      {children}
    </div>
  );
};
