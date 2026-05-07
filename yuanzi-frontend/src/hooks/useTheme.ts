import { useEffect } from 'react';
import { useThemeStore } from '@/stores/useThemeStore';

/**
 * 主题 Hook
 * 应用暗色模式和祖辈模式
 */
export function useTheme() {
  const isDarkMode = useThemeStore((s) => s.isDarkMode);
  const isElderMode = useThemeStore((s) => s.isElderMode);

  useEffect(() => {
    const root = document.documentElement;
    
    if (isDarkMode) {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }

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
