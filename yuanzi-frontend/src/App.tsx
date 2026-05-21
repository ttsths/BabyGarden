import { useEffect, useState, type ReactNode } from 'react';
import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { AdminRouter } from './admin/router';
import { OpenDesignApp } from './pages/OpenDesignApp';
import LoginPage from './pages/LoginPage';
import BabySetupPage from './pages/BabySetupPage';
import FamilyLivePage from './pages/FamilyLivePage';
import { useAuthStore } from './stores/useAuthStore';

function AuthGate({ children, loginPath }: { children: ReactNode; loginPath: string }) {
  const location = useLocation();
  const { token, isAuthenticated, checkAuth } = useAuthStore();
  const [checking, setChecking] = useState(Boolean(token) && !isAuthenticated);

  useEffect(() => {
    let active = true;
    async function verify() {
      if (!token || isAuthenticated) {
        setChecking(false);
        return;
      }
      await checkAuth().catch(() => false);
      if (active) setChecking(false);
    }
    void verify();
    return () => {
      active = false;
    };
  }, [checkAuth, isAuthenticated, token]);

  if (checking) {
    return <div className="min-h-screen bg-background-primary p-6 text-neutral-text-primary">正在校验登录状态...</div>;
  }
  if (!token && !isAuthenticated) {
    return <Navigate to={loginPath} replace state={{ from: location.pathname }} />;
  }
  return <>{children}</>;
}

function App() {
  return (
    <div className="App">
      <Routes>
        <Route path="/admin/*" element={<AdminRouter />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/app/login" element={<LoginPage mode="app" />} />
        <Route path="/app/*" element={<AuthGate loginPath="/app/login"><FamilyLivePage /></AuthGate>} />
        <Route path="/baby/setup" element={<BabySetupPage />} />
        <Route path="/*" element={<AuthGate loginPath="/login"><OpenDesignApp /></AuthGate>} />
      </Routes>
    </div>
  );
}

export default App;
