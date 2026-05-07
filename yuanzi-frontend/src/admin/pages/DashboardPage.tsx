import { useQuery } from '@tanstack/react-query';
import { Card, Row, Col, Statistic, Spin, Alert } from 'antd';
import {
  UserOutlined,
  TeamOutlined,
  SmileOutlined,
  PictureOutlined,
  FileTextOutlined,
} from '@ant-design/icons';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { getStatsOverview, getDailyStats } from '@/admin/api/adminApi';

export function DashboardPage() {
  const {
    data: overviewData,
    isLoading: overviewLoading,
    error: overviewError,
  } = useQuery({
    queryKey: ['admin', 'stats', 'overview'],
    queryFn: async () => {
      const res = await getStatsOverview();
      return res.data.data;
    },
  });

  const {
    data: dailyData,
    isLoading: dailyLoading,
    error: dailyError,
  } = useQuery({
    queryKey: ['admin', 'stats', 'daily'],
    queryFn: async () => {
      const res = await getDailyStats(30);
      return res.data.data;
    },
  });

  if (overviewLoading || dailyLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 64 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (overviewError || dailyError) {
    return (
      <Alert
        message="数据加载失败"
        description="请刷新页面重试"
        type="error"
        showIcon
      />
    );
  }

  const stats = overviewData || { users: 0, families: 0, babies: 0, photos: 0, records: 0 };

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={8} xl={4}>
          <Card>
            <Statistic
              title="用户总数"
              value={stats.users}
              prefix={<UserOutlined style={{ color: '#FF9A8B' }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8} xl={4}>
          <Card>
            <Statistic
              title="家庭总数"
              value={stats.families}
              prefix={<TeamOutlined style={{ color: '#FF9A8B' }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8} xl={4}>
          <Card>
            <Statistic
              title="宝宝总数"
              value={stats.babies}
              prefix={<SmileOutlined style={{ color: '#FF9A8B' }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8} xl={4}>
          <Card>
            <Statistic
              title="照片总数"
              value={stats.photos}
              prefix={<PictureOutlined style={{ color: '#FF9A8B' }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8} xl={4}>
          <Card>
            <Statistic
              title="记录总数"
              value={stats.records}
              prefix={<FileTextOutlined style={{ color: '#FF9A8B' }} />}
            />
          </Card>
        </Col>
      </Row>

      <Card title="近30天趋势" style={{ marginTop: 24 }}>
        <ResponsiveContainer width="100%" height={360}>
          <LineChart data={dailyData || []}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="date"
              tickFormatter={(value: string) =>
                value.slice(5)
              }
            />
            <YAxis allowDecimals={false} />
            <Tooltip />
            <Legend />
            <Line
              type="monotone"
              dataKey="new_users"
              name="新增用户"
              stroke="#FF9A8B"
              strokeWidth={2}
              dot={false}
            />
            <Line
              type="monotone"
              dataKey="new_babies"
              name="新增宝宝"
              stroke="#82ca9d"
              strokeWidth={2}
              dot={false}
            />
            <Line
              type="monotone"
              dataKey="new_records"
              name="新增记录"
              stroke="#8884d8"
              strokeWidth={2}
              dot={false}
            />
          </LineChart>
        </ResponsiveContainer>
      </Card>
    </div>
  );
}
