import { create } from 'zustand';
import { persist } from 'zustand/middleware';

type ThemeMode = 'light' | 'dark' | 'elderly';

interface ThemeState {
  mode: ThemeMode;
  setMode: (mode: ThemeMode) => void;
  toggleElderlyMode: () => void;
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set, get) => ({
      mode: 'light',

      setMode: (mode) => {
        set({ mode });
        applyTheme(mode);
      },

      toggleElderlyMode: () => {
        const newMode = get().mode === 'elderly' ? 'light' : 'elderly';
        set({ mode: newMode });
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
