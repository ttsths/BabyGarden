import { useEffect } from 'react';
import { Outlet, Navigate } from 'react-router-dom';
import { useAdminAuthStore } from '@/admin/store/useAdminAuthStore';

export function AuthGuard() {
  const { isAuthenticated, checkAuth } = useAdminAuthStore();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  if (!isAuthenticated) {
    return <Navigate to="/admin/login" replace />;
  }

  return <Outlet />;
}
