import { create } from 'zustand';
import { persist } from 'zustand/middleware';

type ThemeMode = 'light' | 'dark' | 'elderly';

interface ThemeState {
  mode: ThemeMode;
  isDarkMode: boolean;
  isElderMode: boolean;
  setMode: (mode: ThemeMode) => void;
  toggleDarkMode: () => void;
  toggleElderMode: () => void;
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set, get) => ({
      mode: 'light',
      isDarkMode: false,
      isElderMode: false,

      setMode: (mode) => {
        set({
          mode,
          isDarkMode: mode === 'dark',
          isElderMode: mode === 'elderly',
        });
        applyTheme(mode);
      },

      toggleDarkMode: () => {
        const isDark = get().mode === 'dark';
        const newMode = isDark ? 'light' : 'dark';
        set({ mode: newMode, isDarkMode: !isDark, isElderMode: false });
        applyTheme(newMode);
      },

      toggleElderMode: () => {
        const isElder = get().mode === 'elderly';
        const newMode = isElder ? 'light' : 'elderly';
        set({ mode: newMode, isElderMode: !isElder, isDarkMode: false });
        applyTheme(newMode);
      },
    }),
    {
      name: 'theme-storage',
    }
  )
);

// 应用主题到 DOM
function applyTheme(mode: ThemeMode) {
  const root = document.documentElement;
  root.classList.remove('light', 'dark', 'elderly-mode');
  
  if (mode === 'dark') {
    root.classList.add('dark');
  } else if (mode === 'elderly') {
    root.classList.add('elderly-mode');
  }
}

// 初始化主题
applyTheme(useThemeStore.getState().mode);
