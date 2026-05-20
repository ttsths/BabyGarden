import { Route, Routes } from 'react-router-dom';
import { AdminRouter } from './admin/router';
import FamilyLivePage from './pages/FamilyLivePage';
import LoginPage from './pages/LoginPage';
import BabySetupPage from './pages/BabySetupPage';

function App() {
  return (
    <Routes>
      <Route path="/admin/*" element={<AdminRouter />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/baby/setup" element={<BabySetupPage />} />
      <Route path="*" element={<FamilyLivePage />} />
    </Routes>
  );
}

export default App;
