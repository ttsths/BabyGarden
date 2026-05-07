import { useEffect } from 'react';
import { useThemeStore } from '@/stores/useThemeStore';

/**
 * 主题 Hook
 * 应用暗色模式和祖辈模式
 */
export function useTheme() {
  const { isDarkMode, isElderMode } = useThemeStore();

  useEffect(() => {
    const root = document.documentElement;
    
    // 暗色模式
    if (isDarkMode) {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }

    // 祖辈模式
    if (isElderMode) {
      root.classList.add('elder-mode');
    } else {
      root.classList.remove('elder-mode');
    }
  }, [isDarkMode, isElderMode]);

  return {
    isDarkMode,
    isElderMode,
  };
}
