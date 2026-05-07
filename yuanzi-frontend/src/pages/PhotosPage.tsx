import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { familyMembers, emptyStates } from '@/data/mockData';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { BottomNav } from '@/components/ui/BottomNav';

/**
 * 照片墙页面
 * 展示家庭照片，支持按日期分组
 */
export const PhotosPage: React.FC = () => {
  const navigate = useNavigate();
  const [selectedDate, setSelectedDate] = useState<string | null>(null);

  // 模拟照片数据
  const photoGroups = [
    {
      date: '2026 年 3 月',
      days: [
        {
          date: '3 月 6 日',
          count: 5,
          photos: [
            { id: '1', thumbnail: '📷', time: '10:30' },
            { id: '2', thumbnail: '📷', time: '14:20' },
            { id: '3', thumbnail: '📷', time: '16:45' },
          ],
        },
        {
          date: '3 月 5 日',
          count: 7,
          photos: [
            { id: '4', thumbnail: '📷', time: '09:00' },
            { id: '5', thumbnail: '📷', time: '11:30' },
            { id: '6', thumbnail: '📷', time: '15:15' },
            { id: '7', thumbnail: '📷', time: '18:00' },
          ],
        },
        {
          date: '3 月 4 日',
          count: 3,
          photos: [
            { id: '8', thumbnail: '📷', time: '10:00' },
            { id: '9', thumbnail: '📷', time: '14:30' },
          ],
        },
      ],
    },
  ];

  // 照片详情视图
  if (selectedDate) {
    return (
      <div className="min-h-screen bg-background-primary pb-20">
        <header className="px-6 pt-12 pb-4 bg-background-secondary sticky top-0 z-10">
          <div className="flex items-center justify-between">
            <button
              onClick={() => setSelectedDate(null)}
              className="text-neutral-text-primary"
            >
              ← 返回
            </button>
            <h1 className="text-xl font-semibold text-neutral-text-primary">
              {selectedDate}
            </h1>
            <div className="w-10" />
          </div>
        </header>

        <main className="px-4 py-6">
          <div className="grid grid-cols-3 gap-3">
            {[1, 2, 3, 4, 5, 6, 7, 8, 9].map((i) => (
              <div
                key={i}
                className="aspect-square bg-neutral-border rounded-lg flex items-center justify-center text-4xl"
              >
                📷
              </div>
            ))}
          </div>
        </main>

        <BottomNav />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background-primary pb-20">
      <header className="px-6 pt-12 pb-4 bg-background-secondary sticky top-0 z-10">
        <div className="flex items-center justify-between">
          <button
            onClick={() => navigate(-1)}
            className="text-neutral-text-primary"
          >
            ← 返回
          </button>
          <h1 className="text-xl font-semibold text-neutral-text-primary">
            家庭空间
          </h1>
          <button className="text-brand-primary">
            👨‍👩‍👧
          </button>
        </div>
      </header>

      <main className="px-4 py-6 space-y-6">
        {/* 家庭成员 */}
        <Card padding="medium">
          <h3 className="text-base font-semibold text-neutral-text-primary mb-4">
            家庭成员
          </h3>
          <div className="flex gap-4">
            {familyMembers.map((member) => (
              <div key={member.id} className="flex flex-col items-center gap-2">
                <div className="w-14 h-14 rounded-full bg-brand-light/20 flex items-center justify-center text-2xl">
                  {member.avatar}
                </div>
                <span className="text-sm text-neutral-text-secondary">
                  {member.name}
                </span>
              </div>
            ))}
            <button className="flex flex-col items-center gap-2">
              <div className="w-14 h-14 rounded-full bg-neutral-border flex items-center justify-center text-2xl">
                ➕
              </div>
              <span className="text-sm text-neutral-text-secondary">
                邀请
              </span>
            </button>
          </div>
        </Card>

        {/* 照片分组 */}
        {photoGroups.map((group, groupIndex) => (
          <div key={groupIndex}>
            <h3 className="text-base font-semibold text-neutral-text-primary mb-3 px-1">
              {group.date}
            </h3>
            <div className="space-y-4">
              {group.days.map((day, dayIndex) => (
                <Card
                  key={dayIndex}
                  padding="medium"
                  onClick={() => setSelectedDate(day.date)}
                >
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-sm font-medium text-neutral-text-primary">
                      {day.date}
                    </span>
                    <span className="text-xs text-neutral-text-secondary">
                      {day.count} 条记录
                    </span>
                  </div>
                  <div className="grid grid-cols-3 gap-2">
                    {day.photos.slice(0, 3).map((photo) => (
                      <div
                        key={photo.id}
                        className="aspect-square bg-neutral-border rounded-lg flex items-center justify-center text-3xl"
                      >
                        {photo.thumbnail}
                      </div>
                    ))}
                  </div>
                </Card>
              ))}
            </div>
          </div>
        ))}

        {/* 空状态 */}
        {photoGroups.length === 0 && (
          <div className="text-center py-12">
            <div className="text-6xl mb-4">{emptyStates.noPhotos.icon}</div>
            <h3 className="text-lg font-medium text-neutral-text-primary mb-2">
              {emptyStates.noPhotos.title}
            </h3>
            <Button variant="primary" size="medium">
              {emptyStates.noPhotos.action}
            </Button>
          </div>
        )}
      </main>

      {/* 拍照按钮 */}
      <div className="fixed bottom-24 right-6">
        <button
          className="w-14 h-14 rounded-full bg-brand-primary text-white shadow-lg flex items-center justify-center text-2xl hover:bg-brand-dark transition-colors"
          aria-label="拍一张"
        >
          📷
        </button>
      </div>

      <BottomNav />
    </div>
  );
};

export default PhotosPage;
