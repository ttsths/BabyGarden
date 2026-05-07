import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Table,
  Input,
  Button,
  Drawer,
  Switch,
  Popconfirm,
  message,
  Tag,
  Space,
  Spin,
  Alert,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { SearchOutlined, DeleteOutlined } from '@ant-design/icons';
import {
  getUsers,
  getUserDetail,
  updateUserStatus,
  deleteUser,
} from '@/admin/api/adminApi';
import type { AdminUser } from '@/admin/types/admin';
import dayjs from 'dayjs';

export function UsersPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [keyword, setKeyword] = useState('');
  const [searchValue, setSearchValue] = useState('');
  const [detailId, setDetailId] = useState<string | null>(null);

  const { data, isLoading, error } = useQuery({
    queryKey: ['admin', 'users', page, pageSize, keyword],
    queryFn: async () => {
      const res = await getUsers(page, pageSize, keyword || undefined);
      return res.data.data;
    },
  });

  const { data: detailData, isLoading: detailLoading } = useQuery({
    queryKey: ['admin', 'user', detailId],
    queryFn: async () => {
      if (!detailId) return null;
      const res = await getUserDetail(detailId);
      return res.data.data;
    },
    enabled: !!detailId,
  });

  const statusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: number }) =>
      updateUserStatus(id, status),
    onSuccess: () => {
      message.success('状态更新成功');
      queryClient.invalidateQueries({ queryKey: ['admin', 'users'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'user', detailId] });
    },
    onError: () => {
      message.error('状态更新失败');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteUser(id),
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['admin', 'users'] });
    },
    onError: () => {
      message.error('删除失败');
    },
  });

  const handleSearch = () => {
    setKeyword(searchValue);
    setPage(1);
  };

  const columns: ColumnsType<AdminUser> = [
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      key: 'nickname',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) =>
        status === 1 ? (
          <Tag color="success">正常</Tag>
        ) : (
          <Tag color="default">禁用</Tag>
        ),
    },
    {
      title: '管理员',
      dataIndex: 'is_admin',
      key: 'is_admin',
      render: (isAdmin: boolean) =>
        isAdmin ? <Tag color="warning">是</Tag> : <Tag>否</Tag>,
    },
    {
      title: '注册时间',
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
          <Switch
            checked={record.status === 1}
            onChange={(checked) =>
              statusMutation.mutate({
                id: record.id,
                status: checked ? 1 : 0,
              })
            }
            loading={statusMutation.isPending}
          />
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
      <div style={{ marginBottom: 16, display: 'flex', gap: 8 }}>
        <Input
          placeholder="搜索手机号或昵称"
          value={searchValue}
          onChange={(e) => setSearchValue(e.target.value)}
          onPressEnter={handleSearch}
          style={{ width: 300 }}
          prefix={<SearchOutlined />}
        />
        <Button type="primary" onClick={handleSearch}>
          搜索
        </Button>
      </div>

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
        title="用户详情"
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
              <div className="text-gray-500 text-sm">手机号</div>
              <div>{detailData.phone}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">昵称</div>
              <div>{detailData.nickname}</div>
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
              <div className="text-gray-500 text-sm">状态</div>
              <div>
                {detailData.status === 1 ? (
                  <Tag color="success">正常</Tag>
                ) : (
                  <Tag color="default">禁用</Tag>
                )}
              </div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">管理员</div>
              <div>{detailData.is_admin ? '是' : '否'}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">注册时间</div>
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
