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
  getBabies,
  getBabyDetail,
  deleteBaby,
} from '@/admin/api/adminApi';
import type { AdminBaby } from '@/admin/types/admin';
import dayjs from 'dayjs';

export function BabiesPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [detailId, setDetailId] = useState<string | null>(null);

  const { data, isLoading, error } = useQuery({
    queryKey: ['admin', 'babies', page, pageSize],
    queryFn: async () => {
      const res = await getBabies(page, pageSize);
      return res.data.data;
    },
  });

  const { data: detailData, isLoading: detailLoading } = useQuery({
    queryKey: ['admin', 'baby', detailId],
    queryFn: async () => {
      if (!detailId) return null;
      const res = await getBabyDetail(detailId);
      return res.data.data;
    },
    enabled: !!detailId,
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteBaby(id),
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['admin', 'babies'] });
    },
    onError: () => {
      message.error('删除失败');
    },
  });

  const columns: ColumnsType<AdminBaby> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '名字',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '性别',
      dataIndex: 'gender',
      key: 'gender',
      render: (gender: string) =>
        gender === 'male' ? (
          <Tag color="blue">男</Tag>
        ) : gender === 'female' ? (
          <Tag color="pink">女</Tag>
        ) : (
          <Tag>未知</Tag>
        ),
    },
    {
      title: '生日',
      dataIndex: 'birthday',
      key: 'birthday',
    },
    {
      title: '家庭ID',
      dataIndex: 'family_id',
      key: 'family_id',
      width: 100,
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
        title="宝宝详情"
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
              <div className="text-gray-500 text-sm">名字</div>
              <div>{detailData.name}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">头像</div>
              {detailData.avatar_url ? (
                <img
                  src={detailData.avatar_url}
                  alt="avatar"
                  className="w-16 h-16 rounded-full object-cover"
                />
              ) : (
                <span>无</span>
              )}
            </div>
            <div>
              <div className="text-gray-500 text-sm">性别</div>
              <div>
                {detailData.gender === 'male' ? (
                  <Tag color="blue">男</Tag>
                ) : detailData.gender === 'female' ? (
                  <Tag color="pink">女</Tag>
                ) : (
                  <Tag>未知</Tag>
                )}
              </div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">生日</div>
              <div>{detailData.birthday}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">家庭ID</div>
              <div>{detailData.family_id}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">创建时间</div>
              <div>
                {dayjs(detailData.created_at).format('YYYY-MM-DD HH:mm:ss')}
              </div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">更新时间</div>
              <div>
                {dayjs(detailData.updated_at).format('YYYY-MM-DD HH:mm:ss')}
              </div>
            </div>
          </div>
        ) : null}
      </Drawer>
    </div>
  );
}
