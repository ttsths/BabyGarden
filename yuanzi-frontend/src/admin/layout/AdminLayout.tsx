import { useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import {
  Layout,
  Menu,
  Button,
  theme as antdTheme,
  ConfigProvider,
} from 'antd';
import type { MenuProps } from 'antd';
import {
  DashboardOutlined,
  UserOutlined,
  SmileOutlined,
  TeamOutlined,
  PictureOutlined,
  FileTextOutlined,
  RobotOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { useAdminAuthStore } from '@/admin/store/useAdminAuthStore';

const { Header, Sider, Content } = Layout;

type MenuItem = Required<MenuProps>['items'][number];

const menuItems: MenuItem[] = [
  {
    key: 'dashboard',
    icon: <DashboardOutlined />,
    label: '仪表盘',
  },
  {
    key: 'users',
    icon: <UserOutlined />,
    label: '用户管理',
  },
  {
    key: 'babies',
    icon: <SmileOutlined />,
    label: '宝宝管理',
  },
  {
    key: 'families',
    icon: <TeamOutlined />,
    label: '家庭管理',
  },
  {
    key: 'photos',
    icon: <PictureOutlined />,
    label: '照片管理',
  },
  {
    key: 'records',
    icon: <FileTextOutlined />,
    label: '记录管理',
  },
  {
    key: 'ai-usage',
    icon: <RobotOutlined />,
    label: 'AI统计',
  },
];

export function AdminLayout() {
  const [collapsed, setCollapsed] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const logout = useAdminAuthStore((state) => state.logout);

  const pathKey = location.pathname.split('/').pop() || 'dashboard';

  const handleMenuClick: MenuProps['onClick'] = (e) => {
    navigate(`/admin/${e.key}`);
  };

  const handleLogout = () => {
    logout();
    window.location.href = '/admin/login';
  };

  const {
    token: { colorBgContainer, borderRadiusLG },
  } = antdTheme.useToken();

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: '#FF9A8B',
          borderRadius: 8,
        },
      }}
    >
      <Layout style={{ minHeight: '100vh' }}>
        <Sider
          trigger={null}
          collapsible
          collapsed={collapsed}
          theme="light"
          style={{
            boxShadow: '2px 0 8px rgba(0,0,0,0.06)',
            zIndex: 10,
          }}
        >
          <div
            style={{
              height: 64,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontWeight: 700,
              fontSize: collapsed ? 14 : 18,
              color: '#FF9A8B',
              borderBottom: '1px solid #f0f0f0',
            }}
          >
            {collapsed ? 'YZ' : '圆子管理后台'}
          </div>
          <Menu
            mode="inline"
            selectedKeys={[pathKey]}
            items={menuItems}
            onClick={handleMenuClick}
            style={{ borderRight: 0 }}
          />
        </Sider>
        <Layout>
          <Header
            style={{
              padding: '0 24px',
              background: colorBgContainer,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              boxShadow: '0 2px 8px rgba(0,0,0,0.06)',
            }}
          >
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
            />
            <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
              <span style={{ color: '#666', fontSize: 14 }}>Yuanzi Admin</span>
              <Button
                type="primary"
                danger
                icon={<LogoutOutlined />}
                onClick={handleLogout}
              >
                退出登录
              </Button>
            </div>
          </Header>
          <Content
            style={{
              margin: 24,
              padding: 24,
              background: colorBgContainer,
              borderRadius: borderRadiusLG,
              minHeight: 280,
              overflow: 'auto',
            }}
          >
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  );
}
