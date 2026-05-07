import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Table,
  Button,
  Drawer,
  Popconfirm,
  message,
  Space,
  Spin,
  Alert,
  Tag,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { DeleteOutlined } from '@ant-design/icons';
import {
  getRecords,
  getRecordDetail,
  deleteRecord,
} from '@/admin/api/adminApi';
import type { AdminRecord } from '@/admin/types/admin';
import dayjs from 'dayjs';

export function RecordsPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [detailId, setDetailId] = useState<string | null>(null);

  const { data, isLoading, error } = useQuery({
    queryKey: ['admin', 'records', page, pageSize],
    queryFn: async () => {
      const res = await getRecords(page, pageSize);
      return res.data.data;
    },
  });

  const { data: detailData, isLoading: detailLoading } = useQuery({
    queryKey: ['admin', 'record', detailId],
    queryFn: async () => {
      if (!detailId) return null;
      const res = await getRecordDetail(detailId);
      return res.data.data;
    },
    enabled: !!detailId,
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteRecord(id),
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['admin', 'records'] });
    },
    onError: () => {
      message.error('删除失败');
    },
  });

  const getTypeColor = (type: string) => {
    const map: Record<string, string> = {
      feeding: 'green',
      sleep: 'blue',
      diaper: 'orange',
      temperature: 'red',
      food: 'purple',
      medicine: 'cyan',
    };
    return map[type] || 'default';
  };

  const getTypeLabel = (type: string) => {
    const map: Record<string, string> = {
      feeding: '喂奶',
      sleep: '睡眠',
      diaper: '换尿布',
      temperature: '体温',
      food: '辅食',
      medicine: '用药',
    };
    return map[type] || type;
  };

  const columns: ColumnsType<AdminRecord> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={getTypeColor(type)}>{getTypeLabel(type)}</Tag>
      ),
    },
    {
      title: '宝宝ID',
      dataIndex: 'baby_id',
      key: 'baby_id',
      width: 100,
    },
    {
      title: '家庭ID',
      dataIndex: 'family_id',
      key: 'family_id',
      width: 100,
    },
    {
      title: '创建者ID',
      dataIndex: 'created_by',
      key: 'created_by',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (value: string) => dayjs(value).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button type="link" onClick={() => setDetailId(record.id)}>
            查看
          </Button>
          <Popconfirm
            title="确认删除"
            description="删除后不可恢复，是否继续？"
            onConfirm={() => deleteMutation.mutate(record.id)}
            okText="确认"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  if (error) {
    return (
      <Alert
        message="数据加载失败"
        description="请刷新页面重试"
        type="error"
        showIcon
      />
    );
  }

  return (
    <div>
      <Table
        columns={columns}
        dataSource={data?.list || []}
        rowKey="id"
        loading={isLoading}
        pagination={{
          current: page,
          pageSize,
          total: data?.total || 0,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
          onChange: (p, ps) => {
            setPage(p);
            if (ps) setPageSize(ps);
          },
        }}
      />

      <Drawer
        title="记录详情"
        open={!!detailId}
        onClose={() => setDetailId(null)}
        width={480}
      >
        {detailLoading ? (
          <Spin />
        ) : detailData ? (
          <div className="space-y-4">
            <div>
              <div className="text-gray-500 text-sm">ID</div>
              <div>{detailData.id}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">类型</div>
              <div>
                <Tag color={getTypeColor(detailData.type)}>
                  {getTypeLabel(detailData.type)}
                </Tag>
              </div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">宝宝ID</div>
              <div>{detailData.baby_id}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">家庭ID</div>
              <div>{detailData.family_id}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">创建者ID</div>
              <div>{detailData.created_by}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">数据</div>
              <pre
                style={{
                  background: '#f5f5f5',
                  padding: 12,
                  borderRadius: 8,
                  overflow: 'auto',
                }}
              >
                {JSON.stringify(detailData.data, null, 2)}
              </pre>
            </div>
            <div>
              <div className="text-gray-500 text-sm">创建时间</div>
              <div>
                {dayjs(detailData.created_at).format('YYYY-MM-DD HH:mm:ss')}
              </div>
            </div>
          </div>
        ) : null}
      </Drawer>
    </div>
  );
}
