import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Form, Input, Button, message, ConfigProvider } from 'antd';
import { useAdminAuthStore } from '@/admin/store/useAdminAuthStore';

export function LoginPage() {
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const login = useAdminAuthStore((state) => state.login);

  const handleSubmit = async (values: { phone: string; password: string }) => {
    setLoading(true);
    try {
      await login(values.phone, values.password);
      message.success('登录成功');
      navigate('/admin/dashboard');
    } catch {
      message.error('登录失败，请检查账号密码');
    } finally {
      setLoading(false);
    }
  };

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: '#FF9A8B',
        },
      }}
    >
      <div
        style={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          background: '#FFFBF7',
        }}
      >
        <Card
          title="圆子管理后台"
          style={{ width: 400, borderRadius: 12 }}
          headStyle={{ textAlign: 'center', fontSize: 20, color: '#FF9A8B' }}
        >
          <Form form={form} layout="vertical" onFinish={handleSubmit}>
            <Form.Item
              label="手机号"
              name="phone"
              rules={[
                { required: true, message: '请输入手机号' },
                { pattern: /^1\d{10}$/, message: '手机号格式错误' },
              ]}
            >
              <Input placeholder="请输入管理员手机号" size="large" />
            </Form.Item>
            <Form.Item
              label="密码"
              name="password"
              rules={[{ required: true, message: '请输入密码' }]}
            >
              <Input.Password placeholder="请输入密码" size="large" />
            </Form.Item>
            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                block
                size="large"
                style={{ background: '#FF9A8B', borderColor: '#FF9A8B' }}
              >
                登录
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </div>
    </ConfigProvider>
  );
}
