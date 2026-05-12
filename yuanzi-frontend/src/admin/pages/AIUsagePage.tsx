import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Col, DatePicker, Row, Select, Space, Statistic, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { BarChartOutlined, ClockCircleOutlined, DatabaseOutlined } from '@ant-design/icons';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import { getAIUsage, getAIUsageSummary } from '@/admin/api/adminApi';
import type { AIUsageLog } from '@/admin/types/admin';

const { RangePicker } = DatePicker;

const providerColors: Record<string, string> = {
  grokai: '#FF9A8B',
  cloudflare_workers_ai: '#F59E0B',
  dashscope: '#52C41A',
  cliproxyapi: '#1677FF',
};

export function AIUsagePage() {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [provider, setProvider] = useState<string>();
  const [status, setStatus] = useState<string>();
  const [range, setRange] = useState<[Dayjs | null, Dayjs | null] | null>(null);

  const startDate = range?.[0]?.format('YYYY-MM-DD');
  const endDate = range?.[1]?.format('YYYY-MM-DD');

  const usageQuery = useQuery({
    queryKey: ['admin', 'ai-usage', page, pageSize, provider, status, startDate, endDate],
    queryFn: async () => {
      const res = await getAIUsage({
        page,
        page_size: pageSize,
        provider,
        status,
        request_type: 'chat',
        start_date: startDate,
        end_date: endDate,
      });
      return res.data.data;
    },
  });

  const summaryQuery = useQuery({
    queryKey: ['admin', 'ai-usage-summary', provider, status],
    queryFn: async () => {
      const res = await getAIUsageSummary({
        period: 'day',
        days: 30,
        provider,
        status,
        request_type: 'chat',
      });
      return res.data.data.items;
    },
  });

  // 使用 useMemo 包装 summary 数据，避免每次渲染创建新数组引用
  const summary = useMemo(() => summaryQuery.data || [], [summaryQuery.data]);

  // 使用 useMemo 包装 summary 数据转换，避免依赖项变化导致无限循环
  const { trendData, providerDistribution, overview } = useMemo(() => {
    const byPeriod = new Map<string, Record<string, number | string>>();
    summary.forEach((item) => {
      const row = byPeriod.get(item.period) || { period: item.period };
      row[item.provider] = item.total_tokens;
      byPeriod.set(item.period, row);
    });

    const totals = new Map<string, number>();
    summary.forEach((item) => {
      totals.set(item.provider, (totals.get(item.provider) || 0) + item.total_tokens);
    });

    const now = dayjs();
    const today = now.format('YYYY-MM-DD');
    const weekStart = now.startOf('week');
    const monthStart = now.startOf('month');

    const overviewData = summary.reduce(
      (acc, item) => {
        const period = dayjs(item.period);
        if (item.period === today) acc.today += item.total_tokens;
        if (!period.isBefore(weekStart, 'day')) acc.week += item.total_tokens;
        if (!period.isBefore(monthStart, 'day')) acc.month += item.total_tokens;
        return acc;
      },
      { today: 0, week: 0, month: 0 }
    );

    return {
      trendData: Array.from(byPeriod.values()),
      providerDistribution: Array.from(totals.entries()).map(([name, value]) => ({ name, value })),
      overview: overviewData,
    };
  }, [summary]);

  const providers = Array.from(new Set(summary.map((item) => item.provider)));
  const tableData = usageQuery.data?.list || [];
  const pagination = usageQuery.data?.pagination;

  const columns: ColumnsType<AIUsageLog> = [
    {
      title: '时间',
      dataIndex: 'created_at',
      width: 170,
      render: (value: string) => dayjs(value).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '用户',
      dataIndex: 'user_id',
      ellipsis: true,
    },
    {
      title: 'Provider',
      dataIndex: 'provider',
      render: (value: string) => <Tag color={providerColors[value] || 'default'}>{value || 'unknown'}</Tag>,
    },
    {
      title: '模型',
      dataIndex: 'model',
      ellipsis: true,
    },
    {
      title: '输入',
      dataIndex: 'input_tokens',
      align: 'right',
    },
    {
      title: '输出',
      dataIndex: 'output_tokens',
      align: 'right',
    },
    {
      title: '缓存命中',
      dataIndex: 'cached_tokens',
      align: 'right',
    },
    {
      title: '总计',
      dataIndex: 'total_tokens',
      align: 'right',
      sorter: (a, b) => a.total_tokens - b.total_tokens,
    },
    {
      title: '状态',
      dataIndex: 'status',
      render: (value: string) => (
        <Tag color={value === 'success' ? 'success' : 'error'}>{value === 'success' ? '成功' : '失败'}</Tag>
      ),
    },
  ];

  if (usageQuery.error || summaryQuery.error) {
    return <Alert message="AI使用统计加载失败" description="请刷新页面重试" type="error" showIcon />;
  }

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col xs={24} md={8}>
          <Card>
            <Statistic title="今日Token" value={overview.today} prefix={<ClockCircleOutlined />} />
          </Card>
        </Col>
        <Col xs={24} md={8}>
          <Card>
            <Statistic title="本周Token" value={overview.week} prefix={<BarChartOutlined />} />
          </Card>
        </Col>
        <Col xs={24} md={8}>
          <Card>
            <Statistic title="本月Token" value={overview.month} prefix={<DatabaseOutlined />} />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={16}>
          <Card title="Token使用趋势">
            <ResponsiveContainer width="100%" height={320}>
              <LineChart data={trendData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="period" />
                <YAxis allowDecimals={false} />
                <Tooltip />
                <Legend />
                {providers.map((name) => (
                  <Line
                    key={name}
                    type="monotone"
                    dataKey={name}
                    stroke={providerColors[name] || '#8C8C8C'}
                    strokeWidth={2}
                    dot={false}
                  />
                ))}
              </LineChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card title="Provider分布">
            <ResponsiveContainer width="100%" height={320}>
              <PieChart>
                <Pie data={providerDistribution} dataKey="value" nameKey="name" innerRadius={64} outerRadius={100} label>
                  {providerDistribution.map((item) => (
                    <Cell key={item.name} fill={providerColors[item.name] || '#8C8C8C'} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>

      <Card
        title="AI对话明细"
        style={{ marginTop: 16 }}
        extra={
          <Space wrap>
            <Select
              allowClear
              placeholder="Provider"
              value={provider}
              onChange={(value) => {
                setProvider(value);
                setPage(1);
              }}
              style={{ width: 210 }}
              options={[
                { label: 'GrokAI', value: 'grokai' },
                { label: 'Cloudflare Workers AI', value: 'cloudflare_workers_ai' },
                { label: 'DashScope', value: 'dashscope' },
                { label: 'CLIProxyAPI', value: 'cliproxyapi' },
              ]}
            />
            <Select
              allowClear
              placeholder="状态"
              value={status}
              onChange={(value) => {
                setStatus(value);
                setPage(1);
              }}
              style={{ width: 120 }}
              options={[
                { label: '成功', value: 'success' },
                { label: '失败', value: 'error' },
              ]}
            />
            <RangePicker
              value={range}
              onChange={(value) => {
                setRange(value);
                setPage(1);
              }}
            />
          </Space>
        }
      >
        <Table
          columns={columns}
          dataSource={tableData}
          rowKey="id"
          loading={usageQuery.isLoading}
          scroll={{ x: 1100 }}
          pagination={{
            current: page,
            pageSize,
            total: pagination?.total || 0,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
            onChange: (p, ps) => {
              setPage(p);
              if (ps) setPageSize(ps);
            },
          }}
        />
      </Card>
    </div>
  );
}
