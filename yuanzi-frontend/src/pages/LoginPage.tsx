import React, { useMemo, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card } from '@/components/ui/Card';
import { useAuthStore } from '@/stores/useAuthStore';
import { api } from '@/services/api';

type LoginClient = 'pc' | 'app' | 'admin';
type LoginMethod = 'password' | 'sms';

interface LoginPageProps {
  mode?: LoginClient;
}

function inferClient(pathname: string, explicit?: LoginClient): LoginClient {
  if (explicit) return explicit;
  if (pathname.startsWith('/app')) return 'app';
  if (pathname.startsWith('/admin')) return 'admin';
  return 'pc';
}

export const LoginPage: React.FC<LoginPageProps> = ({ mode }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuthStore();
  const client = inferClient(location.pathname, mode);
  const [method, setMethod] = useState<LoginMethod>(client === 'app' ? 'sms' : 'password');
  const [identifier, setIdentifier] = useState('mom');
  const [password, setPassword] = useState('yuanzi123');
  const [phone, setPhone] = useState('13800138000');
  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [error, setError] = useState('');

  const copy = useMemo(() => {
    if (client === 'app') return { title: '小园子 APP', subtitle: '家庭成员登录后同步记录' };
    if (client === 'admin') return { title: '圆子管理端', subtitle: '管理员账号密码登录' };
    return { title: 'BabyGarden', subtitle: 'PC 工作台账号密码登录' };
  }, [client]);

  const redirectTo = client === 'app' ? '/app' : '/';

  const handleSendCode = async () => {
    if (!phone || countdown > 0) return;
    setError('');
    setLoading(true);
    try {
      await api.auth.sendCode(phone);
      setCountdown(60);
      const timer = window.setInterval(() => {
        setCountdown((prev) => {
          if (prev <= 1) {
            window.clearInterval(timer);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    } catch (err) {
      setError(err instanceof Error ? err.message : '验证码暂不可用');
    } finally {
      setLoading(false);
    }
  };

  const handleLogin = async () => {
    setError('');
    setLoading(true);
    try {
      const result = method === 'password'
        ? await api.auth.passwordLogin(identifier, password)
        : await api.auth.login(phone, code);
      await login(result.access_token, result.refresh_token);
      navigate(redirectTo, { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败，请检查账号信息');
    } finally {
      setLoading(false);
    }
  };

  const canSubmit = method === 'password'
    ? Boolean(identifier && password)
    : Boolean(phone && code);

  return (
    <div className="min-h-screen bg-gradient-to-b from-brand-light/20 to-background-primary flex flex-col">
      <div className="flex-1 flex items-center justify-center px-6">
        <div className="w-full max-w-sm">
          <div className="text-center mb-8">
            <div className="w-20 h-20 mx-auto mb-4 rounded-2xl bg-brand-primary flex items-center justify-center text-4xl shadow-lg">
              👶
            </div>
            <h1 className="text-2xl font-bold text-neutral-text-primary mb-2">{copy.title}</h1>
            <p className="text-neutral-text-secondary">{copy.subtitle}</p>
          </div>

          <Card padding="large">
            {client === 'app' && (
              <div className="mb-5 grid grid-cols-2 rounded-xl bg-background-tertiary p-1">
                <button
                  type="button"
                  onClick={() => setMethod('sms')}
                  className={`rounded-lg py-2 text-sm font-medium ${method === 'sms' ? 'bg-white text-brand-primary shadow-sm' : 'text-neutral-text-secondary'}`}
                >
                  手机验证码
                </button>
                <button
                  type="button"
                  onClick={() => setMethod('password')}
                  className={`rounded-lg py-2 text-sm font-medium ${method === 'password' ? 'bg-white text-brand-primary shadow-sm' : 'text-neutral-text-secondary'}`}
                >
                  账号密码
                </button>
              </div>
            )}

            <div className="space-y-4">
              {method === 'password' ? (
                <>
                  <Input
                    label="手机号或用户名"
                    placeholder="请输入手机号或用户名"
                    value={identifier}
                    onChange={(event) => setIdentifier(event.target.value)}
                  />
                  <Input
                    label="密码"
                    type="password"
                    placeholder="请输入密码"
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                  />
                </>
              ) : (
                <>
                  <Input
                    label="手机号"
                    type="tel"
                    placeholder="请输入手机号"
                    value={phone}
                    onChange={(event) => setPhone(event.target.value)}
                    maxLength={11}
                  />
                  <div>
                    <label className="block text-sm font-medium text-neutral-text-primary mb-2">验证码</label>
                    <div className="flex items-center gap-2">
                      <input
                        aria-label="验证码"
                        type="text"
                        placeholder="请输入验证码"
                        value={code}
                        onChange={(event) => setCode(event.target.value)}
                        maxLength={6}
                        className="min-w-0 flex-1 rounded-lg border border-neutral-border bg-background-primary px-3 py-2 text-neutral-text-primary"
                      />
                      <button
                        type="button"
                        onClick={handleSendCode}
                        disabled={countdown > 0 || !phone || loading}
                        className={`relative z-10 h-10 rounded-lg px-3 text-sm font-medium whitespace-nowrap ${
                          countdown > 0 || !phone
                            ? 'bg-neutral-border text-neutral-text-secondary'
                            : 'bg-brand-primary text-white shadow-sm'
                        }`}
                      >
                        {countdown > 0 ? `${countdown}秒后重试` : '获取验证码'}
                      </button>
                    </div>
                  </div>
                </>
              )}

              {error && <div className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600">{error}</div>}

              <Button
                variant="primary"
                size="large"
                fullWidth
                loading={loading}
                disabled={!canSubmit || loading}
                onClick={handleLogin}
                className="mt-6"
              >
                登录
              </Button>
            </div>
          </Card>

          <p className="text-center text-xs text-neutral-text-secondary mt-6">
            登录即代表您同意 <button className="text-brand-primary">用户协议</button> 和 <button className="text-brand-primary">隐私政策</button>
          </p>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
