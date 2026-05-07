import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card } from '@/components/ui/Card';
import { useAuthStore } from '@/stores/useAuthStore';

/**
 * 登录页面
 */
export const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const { login } = useAuthStore();
  const [phone, setPhone] = useState('');
  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [countdown, setCountdown] = useState(0);

  // 发送验证码
  const handleSendCode = async () => {
    if (!phone) return;
    
    setLoading(true);
    // TODO: 调用 API 发送验证码
    await new Promise(resolve => setTimeout(resolve, 1000));
    setLoading(false);
    
    // 开始倒计时
    setCountdown(60);
    const timer = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          clearInterval(timer);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  };

  // 登录
  const handleLogin = async () => {
    if (!phone || !code) return;
    
    setLoading(true);
    // TODO: 调用 API 登录
    await new Promise(resolve => setTimeout(resolve, 1000));
    login('test-user-id');
    setLoading(false);
    navigate('/');
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-brand-light/20 to-background-primary flex flex-col">
      {/* 顶部装饰 */}
      <div className="flex-1 flex items-center justify-center px-6">
        <div className="w-full max-w-sm">
          {/* Logo */}
          <div className="text-center mb-8">
            <div className="w-20 h-20 mx-auto mb-4 rounded-2xl bg-brand-primary flex items-center justify-center text-4xl shadow-lg">
              👶
            </div>
            <h1 className="text-2xl font-bold text-neutral-text-primary mb-2">
              圆子
            </h1>
            <p className="text-neutral-text-secondary">
              温柔的记录者
            </p>
          </div>

          {/* 登录表单 */}
          <Card padding="large">
            <div className="space-y-4">
              <Input
                label="手机号"
                type="tel"
                placeholder="请输入手机号"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                maxLength={11}
              />

              <div>
                <div className="flex items-center gap-2">
                  <Input
                    label="验证码"
                    type="text"
                    placeholder="请输入验证码"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    maxLength={6}
                    className="flex-1"
                  />
                  <button
                    onClick={handleSendCode}
                    disabled={countdown > 0 || !phone}
                    className={`mt-6 px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap ${
                      countdown > 0
                        ? 'bg-neutral-border text-neutral-text-secondary'
                        : 'bg-brand-primary text-white'
                    }`}
                  >
                    {countdown > 0 ? `${countdown}秒后重试` : '获取验证码'}
                  </button>
                </div>
              </div>

              <Button
                variant="primary"
                size="large"
                fullWidth
                loading={loading}
                disabled={!phone || !code || loading}
                onClick={handleLogin}
                className="mt-6"
              >
                登录
              </Button>
            </div>
          </Card>

          {/* 协议 */}
          <p className="text-center text-xs text-neutral-text-secondary mt-6">
            登录即代表您同意{' '}
            <button className="text-brand-primary">用户协议</button>
            {' '}和{' '}
            <button className="text-brand-primary">隐私政策</button>
          </p>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
