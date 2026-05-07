import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '@/components/ui/Card';
import { BottomNav } from '@/components/ui/BottomNav';
import { useThemeStore } from '@/stores/useThemeStore';

interface SettingItemProps {
  icon: string;
  title: string;
  subtitle?: string;
  rightElement?: React.ReactNode;
  onClick?: () => void;
}

const SettingItem: React.FC<SettingItemProps> = ({
  icon,
  title,
  subtitle,
  rightElement,
  onClick,
}) => {
  return (
    <button
      onClick={onClick}
      className="w-full flex items-center gap-4 p-4 hover:bg-background-tertiary transition-colors"
    >
      <span className="text-xl">{icon}</span>
      <div className="flex-1 text-left">
        <div className="text-neutral-text-primary font-medium">{title}</div>
        {subtitle && (
          <div className="text-sm text-neutral-text-secondary mt-0.5">{subtitle}</div>
        )}
      </div>
      {rightElement || (
        <span className="text-neutral-text-secondary">›</span>
      )}
    </button>
  );
};

/**
 * 设置页面
 */
export const SettingsPage: React.FC = () => {
  const navigate = useNavigate();
  const { isDarkMode, isElderMode, toggleDarkMode, toggleElderMode } = useThemeStore();

  // 模拟用户数据
  const user = {
    name: '申屠海三',
    email: 'shentuhaisan@example.com',
  };

  const baby = {
    name: '圆子',
    age: '3 个月 15 天',
  };

  return (
    <div className="min-h-screen bg-background-primary pb-20">
      <header className="px-6 pt-12 pb-4 bg-background-secondary sticky top-0 z-10">
        <div className="flex items-center">
          <button
            onClick={() => navigate(-1)}
            className="text-neutral-text-primary mr-4"
          >
            ← 返回
          </button>
          <h1 className="text-xl font-semibold text-neutral-text-primary">
            设置
          </h1>
        </div>
      </header>

      <main className="px-4 py-6 space-y-6">
        {/* 账号 */}
        <section>
          <h3 className="text-sm font-medium text-neutral-text-secondary mb-2 px-1">
            账号
          </h3>
          <Card padding="none">
            <SettingItem
              icon="👤"
              title={user.name}
              subtitle={user.email}
            />
          </Card>
        </section>

        {/* 宝宝管理 */}
        <section>
          <h3 className="text-sm font-medium text-neutral-text-secondary mb-2 px-1">
            宝宝管理
          </h3>
          <Card padding="none">
            <SettingItem
              icon="👶"
              title={baby.name}
              subtitle={baby.age}
            />
            <div className="border-t border-neutral-border" />
            <SettingItem
              icon="➕"
              title="添加宝宝"
            />
          </Card>
        </section>

        {/* 家庭圈 */}
        <section>
          <h3 className="text-sm font-medium text-neutral-text-secondary mb-2 px-1">
            家庭圈
          </h3>
          <Card padding="none">
            <SettingItem
              icon="👨‍👩‍👧"
              title="家庭成员管理"
            />
          </Card>
        </section>

        {/* 通知 */}
        <section>
          <h3 className="text-sm font-medium text-neutral-text-secondary mb-2 px-1">
            通知
          </h3>
          <Card padding="none">
            <SettingItem
              icon="🔔"
              title="喂奶提醒"
              rightElement={
                <button
                  onClick={() => {}}
                  className="w-12 h-6 rounded-full bg-brand-primary relative"
                >
                  <span className="absolute right-1 top-1 w-4 h-4 rounded-full bg-white" />
                </button>
              }
            />
            <div className="border-t border-neutral-border" />
            <SettingItem
              icon="📱"
              title="推送通知"
            />
          </Card>
        </section>

        {/* 外观 */}
        <section>
          <h3 className="text-sm font-medium text-neutral-text-secondary mb-2 px-1">
            外观
          </h3>
          <Card padding="none">
            <SettingItem
              icon="🌙"
              title="深色模式"
              rightElement={
                <button
                  onClick={toggleDarkMode}
                  className={`w-12 h-6 rounded-full relative transition-colors ${
                    isDarkMode ? 'bg-brand-primary' : 'bg-neutral-border'
                  }`}
                >
                  <span
                    className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${
                      isDarkMode ? 'right-1' : 'left-1'
                    }`}
                  />
                </button>
              }
            />
            <div className="border-t border-neutral-border" />
            <SettingItem
              icon="👵"
              title="祖辈模式"
              subtitle="字体放大，高对比度"
              rightElement={
                <button
                  onClick={toggleElderMode}
                  className={`w-12 h-6 rounded-full relative transition-colors ${
                    isElderMode ? 'bg-brand-primary' : 'bg-neutral-border'
                  }`}
                >
                  <span
                    className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-transform ${
                      isElderMode ? 'right-1' : 'left-1'
                    }`}
                  />
                </button>
              }
            />
          </Card>
        </section>

        {/* 数据 */}
        <section>
          <h3 className="text-sm font-medium text-neutral-text-secondary mb-2 px-1">
            数据
          </h3>
          <Card padding="none">
            <SettingItem
              icon="💾"
              title="数据备份"
            />
            <div className="border-t border-neutral-border" />
            <SettingItem
              icon="📤"
              title="导出数据"
            />
          </Card>
        </section>

        {/* 关于 */}
        <section>
          <h3 className="text-sm font-medium text-neutral-text-secondary mb-2 px-1">
            关于
          </h3>
          <Card padding="none">
            <SettingItem
              icon="ℹ️"
              title="关于我们"
            />
            <div className="border-t border-neutral-border" />
            <SettingItem
              icon="📜"
              title="用户协议"
            />
            <div className="border-t border-neutral-border" />
            <SettingItem
              icon="🔒"
              title="隐私政策"
            />
          </Card>
        </section>

        {/* 版本信息 */}
        <div className="text-center text-sm text-neutral-text-secondary py-4">
          版本 1.0.0 (Build 20260306)
        </div>
      </main>

      <BottomNav />
    </div>
  );
};

export default SettingsPage;
