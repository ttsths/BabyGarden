import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { bottomNavItems } from '@/data/mockData';
import styles from './BottomNav.module.css';

interface BottomNavProps {
  items?: typeof bottomNavItems;
}

/**
 * 底部导航栏
 * 固定在页面底部，提供主要导航功能
 */
export const BottomNav: React.FC<BottomNavProps> = ({ items = bottomNavItems }) => {
  const location = useLocation();

  return (
    <nav className={styles.bottomNav}>
      {items.map((item) => {
        const isActive = location.pathname === item.path;
        const isRecord = item.id === 'record';
        
        return (
          <Link
            key={item.id}
            to={item.path}
            className={`${styles.navItem} ${isActive ? styles.active : ''} ${isRecord ? styles.record : ''}`}
          >
            <span className={styles.navIcon}>{item.icon}</span>
            <span className={styles.navLabel}>{item.label}</span>
          </Link>
        );
      })}
    </nav>
  );
};

export default BottomNav;
