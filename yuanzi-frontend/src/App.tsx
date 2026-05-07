import { Routes, Route } from 'react-router-dom';
import { useState } from 'react';
import { HomePage } from './pages/HomePage';
import { AddRecordModal } from './components/AddRecordModal';
import { AdminRouter } from './admin/router';

// 主应用 - 整合所有组件
function App() {
  const [isModalOpen, setIsModalOpen] = useState(false);

  // 示例数据（按实际设计稿）
  const [stats, setStats] = useState({
    feedingCount: 5,
    sleepHours: 12,
    diaperCount: 8,
  });

  const [recentRecords, setRecentRecords] = useState([
    {
      id: '1',
      type: 'feeding' as const,
      title: 'Bottle feeding',
      subtitle: 'Formula 120ml',
      time: '20m ago',
    },
    {
      id: '2',
      type: 'diaper' as const,
      title: 'Diaper changed',
      subtitle: 'Wet · No issues',
      time: '2h ago',
    },
    {
      id: '3',
      type: 'sleep' as const,
      title: 'Nap ended',
      subtitle: 'Duration: 2h 15m',
      time: '3.5h ago',
    },
  ]);

  const handleRecordTypeClick = (type: 'feeding' | 'sleep' | 'diaper') => {
    console.log('Quick action:', type);
    setIsModalOpen(true);
  };

  const handleSaveRecord = (record: any) => {
    console.log('Saving record:', record);
    // 添加记录逻辑
    const newRecord = {
      id: Date.now().toString(),
      type: record.type,
      title: record.type === 'feeding' ? `Bottle feeding (${record.amount}ml)` :
             record.type === 'sleep' ? 'Sleep record' :
             'Diaper changed',
      subtitle: 'Just now',
      time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
    };
    setRecentRecords([newRecord, ...recentRecords]);
    setStats((prev) => ({
      ...prev,
      [record.type === 'feeding' ? 'feedingCount' : record.type === 'sleep' ? 'sleepHours' : 'diaperCount']:
         prev[record.type === 'feeding' ? 'feedingCount' : record.type === 'sleep' ? 'sleepHours' : 'diaperCount'] + 1
    }));
    setIsModalOpen(false);
  };

  return (
    <div className="App">
      <Routes>
        <Route path="/admin/*" element={<AdminRouter />} />
        <Route
          path="*"
          element={
            <>
              <HomePage
        userName="Mommy"
        date="Monday, October 23rd"
        stats={stats}
        recentRecords={recentRecords}
        onRecordTypeClick={handleRecordTypeClick}
        onViewAllRecords={() => console.log('View all records')}
        onRecordClick={(id) => console.log('Record clicked:', id)}
      />

      <AddRecordModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSave={handleSaveRecord}
      />
            </>
          }
        />
      </Routes>
    </div>
  );
}

export default App;
