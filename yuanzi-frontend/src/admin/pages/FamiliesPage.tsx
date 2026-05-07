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
  List,
  Tag,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { DeleteOutlined } from '@ant-design/icons';
import {
  getFamilies,
  getFamilyDetail,
  deleteFamily,
} from '@/admin/api/adminApi';
import type { AdminFamily } from '@/admin/types/admin';
import dayjs from 'dayjs';

export function FamiliesPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [detailId, setDetailId] = useState<string | null>(null);

  const { data, isLoading, error } = useQuery({
    queryKey: ['admin', 'families', page, pageSize],
    queryFn: async () => {
      const res = await getFamilies(page, pageSize);
      return res.data.data;
    },
  });

  const { data: detailData, isLoading: detailLoading } = useQuery({
    queryKey: ['admin', 'family', detailId],
    queryFn: async () => {
      if (!detailId) return null;
      const res = await getFamilyDetail(detailId);
      return res.data.data;
    },
    enabled: !!detailId,
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteFamily(id),
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['admin', 'families'] });
    },
    onError: () => {
      message.error('删除失败');
    },
  });

  const columns: ColumnsType<AdminFamily> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '拥有者ID',
      dataIndex: 'owner_id',
      key: 'owner_id',
    },
    {
      title: '存储配额',
      dataIndex: 'storage_quota',
      key: 'storage_quota',
      render: (v: number) => `${(v / 1024 / 1024).toFixed(2)} MB`,
    },
    {
      title: '已用存储',
      dataIndex: 'storage_used',
      key: 'storage_used',
      render: (v: number) => `${(v / 1024 / 1024).toFixed(2)} MB`,
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
        title="家庭详情"
        open={!!detailId}
        onClose={() => setDetailId(null)}
        width={560}
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
              <div className="text-gray-500 text-sm">名称</div>
              <div>{detailData.name}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">拥有者ID</div>
              <div>{detailData.owner_id}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">存储配额</div>
              <div>
                {(detailData.storage_quota / 1024 / 1024).toFixed(2)} MB
              </div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">已用存储</div>
              <div>
                {(detailData.storage_used / 1024 / 1024).toFixed(2)} MB
              </div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">创建时间</div>
              <div>
                {dayjs(detailData.created_at).format('YYYY-MM-DD HH:mm:ss')}
              </div>
            </div>
            <div>
              <div className="text-gray-500 text-sm mb-2">成员列表</div>
              <List
                dataSource={detailData.members || []}
                renderItem={(member) => (
                  <List.Item>
                    <List.Item.Meta
                      avatar={
                        member.avatar_url ? (
                          <img
                            src={member.avatar_url}
                            alt="avatar"
                            className="w-10 h-10 rounded-full object-cover"
                          />
                        ) : null
                      }
                      title={member.nickname || member.user_id}
                      description={
                        <Space>
                          <Tag color={member.role === 'admin' ? 'red' : 'blue'}>
                            {member.role === 'admin' ? '管理员' : '成员'}
                          </Tag>
                          <span>
                            加入时间：
                            {dayjs(member.joined_at).format(
                              'YYYY-MM-DD HH:mm:ss'
                            )}
                          </span>
                        </Space>
                      }
                    />
                  </List.Item>
                )}
              />
            </div>
          </div>
        ) : null}
      </Drawer>
    </div>
  );
}
