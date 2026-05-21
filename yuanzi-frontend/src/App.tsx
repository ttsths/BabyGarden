import { Routes, Route } from 'react-router-dom';
import { AdminRouter } from './admin/router';
import { OpenDesignApp } from './pages/OpenDesignApp';
import LoginPage from './pages/LoginPage';
import BabySetupPage from './pages/BabySetupPage';

function App() {
  return (
    <div className="App">
      <Routes>
        <Route path="/admin/*" element={<AdminRouter />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/baby/setup" element={<BabySetupPage />} />
        <Route path="/*" element={<OpenDesignApp />} />
      </Routes>
    </div>
  );
}

export default App;
