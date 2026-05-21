import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '@/types/models';
import { api } from '@/services/api';

interface AuthState {
  // 状态
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  
  // 动作
  login: (token: string, refreshToken: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => Promise<boolean>;
  updateProfile: (profile: Partial<User>) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (token, refreshToken) => {
        set({ token, refreshToken, isLoading: true });
        try {
          // 获取用户信息
          const user = await api.auth.getProfile() as User;
          set({ 
            user, 
            token, 
            refreshToken, 
            isAuthenticated: true,
            isLoading: false 
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      logout: () => {
        set({ 
          user: null, 
          token: null, 
          refreshToken: null, 
          isAuthenticated: false 
        });
      },

      checkAuth: async () => {
        const { token } = get();
        if (!token) {
          return false;
        }
        
        try {
          const user = await api.auth.getProfile() as User;
          set({ user, isAuthenticated: true });
          return true;
        } catch {
          get().logout();
          return false;
        }
      },

      updateProfile: (profile) => {
        set((state) => ({
          user: state.user ? { ...state.user, ...profile } : null,
        }));
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        token: state.token,
        refreshToken: state.refreshToken,
      }),
    }
  )
);

function unwrapResponse<T>(response: unknown): T {
  if (typeof response === 'object' && response !== null && 'data' in response) {
    return (response as { data: T }).data;
  }
  return response as T;
}
