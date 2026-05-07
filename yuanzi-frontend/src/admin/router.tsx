import { Routes, Route, Navigate } from 'react-router-dom';
import { LoginPage } from '@/admin/pages/LoginPage';
import { AuthGuard } from '@/admin/components/AuthGuard';
import { AdminLayout } from '@/admin/layout/AdminLayout';
import { DashboardPage } from '@/admin/pages/DashboardPage';
import { UsersPage } from '@/admin/pages/UsersPage';
import { BabiesPage } from '@/admin/pages/BabiesPage';
import { FamiliesPage } from '@/admin/pages/FamiliesPage';
import { PhotosPage } from '@/admin/pages/PhotosPage';
import { RecordsPage } from '@/admin/pages/RecordsPage';

export function AdminRouter() {
  return (
    <Routes>
      <Route path="login" element={<LoginPage />} />
      <Route element={<AuthGuard />}>
        <Route element={<AdminLayout />}>
          <Route index element={<Navigate to="dashboard" replace />} />
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="users" element={<UsersPage />} />
          <Route path="babies" element={<BabiesPage />} />
          <Route path="families" element={<FamiliesPage />} />
          <Route path="photos" element={<PhotosPage />} />
          <Route path="records" element={<RecordsPage />} />
        </Route>
      </Route>
    </Routes>
  );
}
