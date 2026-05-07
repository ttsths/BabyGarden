import { create } from 'zustand';
import { adminLogin } from '@/admin/api/adminApi';

interface AdminAuthState {
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (phone: string, password: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => boolean;
}

export const useAdminAuthStore = create<AdminAuthState>((set) => ({
  token: localStorage.getItem('admin_token'),
  isAuthenticated: !!localStorage.getItem('admin_token'),
  isLoading: false,
  error: null,

  login: async (phone: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      const res = await adminLogin(phone, password);
      const token = res.data.data.token;
      localStorage.setItem('admin_token', token);
      set({ token, isAuthenticated: true, isLoading: false });
    } catch (err) {
      const msg = err instanceof Error ? err.message : '登录失败';
      set({ error: msg, isLoading: false });
      throw err;
    }
  },

  logout: () => {
    localStorage.removeItem('admin_token');
    set({ token: null, isAuthenticated: false, error: null });
  },

  checkAuth: () => {
    const token = localStorage.getItem('admin_token');
    const isAuth = !!token;
    set({ token, isAuthenticated: isAuth });
    return isAuth;
  },
}));
