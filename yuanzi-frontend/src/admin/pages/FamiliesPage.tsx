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
  Form,
  Modal,
  Input,
  Select,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { DeleteOutlined, PlusOutlined, EditOutlined, UserAddOutlined } from '@ant-design/icons';
import {
  getFamilies,
  getFamilyDetail,
  deleteFamily,
  createFamily,
  updateFamily,
  addFamilyMember,
  getUsers,
} from '@/admin/api/adminApi';
import type { AdminFamily, AdminFamilyMember, AdminUser, CreateFamilyRequest, UpdateFamilyRequest, AddFamilyMemberRequest } from '@/admin/types/admin';
import dayjs from 'dayjs';

export function FamiliesPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [detailId, setDetailId] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isMemberModalOpen, setIsMemberModalOpen] = useState(false);
  const [editingFamily, setEditingFamily] = useState<AdminFamily | null>(null);
  const [currentFamilyId, setCurrentFamilyId] = useState<string | null>(null);
  const [form] = Form.useForm();
  const [memberForm] = Form.useForm();

  const { data: usersData } = useQuery({
    queryKey: ['admin', 'users', 'all'],
    queryFn: async () => {
      const res = await getUsers(1, 1000);
      return res.data.data?.list || [];
    },
  });

  const userOptions = (usersData || []).map((u: AdminUser) => ({ label: `${u.nickname || u.phone} (${u.phone})`, value: u.id }));

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
      message.success('家庭已禁用');
      queryClient.invalidateQueries({ queryKey: ['admin', 'families'] });
    },
    onError: () => {
      message.error('操作失败');
    },
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateFamilyRequest) => createFamily(data),
    onSuccess: () => {
      message.success('创建成功');
      setIsModalOpen(false);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['admin', 'families'] });
    },
    onError: () => {
      message.error('创建失败');
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateFamilyRequest }) => updateFamily(id, data),
    onSuccess: () => {
      message.success('更新成功');
      setIsModalOpen(false);
      setEditingFamily(null);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['admin', 'families'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'family', detailId] });
    },
    onError: () => {
      message.error('更新失败');
    },
  });

  const addMemberMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: AddFamilyMemberRequest }) => addFamilyMember(id, data),
    onSuccess: () => {
      message.success('添加成员成功');
      setIsMemberModalOpen(false);
      memberForm.resetFields();
      queryClient.invalidateQueries({ queryKey: ['admin', 'family', currentFamilyId] });
    },
    onError: () => {
      message.error('添加成员失败');
    },
  });

  const handleOpenCreate = () => {
    setEditingFamily(null);
    form.resetFields();
    setIsModalOpen(true);
  };

  const handleOpenEdit = (record: AdminFamily) => {
    setEditingFamily(record);
    form.setFieldsValue({
      name: record.name,
      invite_code: record.invite_code,
    });
    setIsModalOpen(true);
  };

  const handleOpenAddMember = (familyId: string) => {
    setCurrentFamilyId(familyId);
    memberForm.resetFields();
    setIsMemberModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingFamily) {
        updateMutation.mutate({ id: editingFamily.id, data: values });
      } else {
        createMutation.mutate(values);
      }
    } catch {
      // validation failed
    }
  };

  const handleAddMemberSubmit = async () => {
    try {
      const values = await memberForm.validateFields();
      if (currentFamilyId) {
        addMemberMutation.mutate({ id: currentFamilyId, data: values });
      }
    } catch {
      // validation failed
    }
  };

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
      title: '成员数',
      dataIndex: 'member_count',
      key: 'member_count',
      render: (v: number) => v ?? '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) =>
        status === 1 ? (
          <Tag color="green">启用</Tag>
        ) : (
          <Tag color="red">禁用</Tag>
        ),
    },
    {
      title: '创建时间',
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
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
          </Button>
          <Button type="link" icon={<UserAddOutlined />} onClick={() => handleOpenAddMember(record.id)}>
            添加成员
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
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenCreate}>
          新增家庭
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
                renderItem={(member: AdminFamilyMember) => (
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

      <Modal
        title={editingFamily ? '编辑家庭' : '新增家庭'}
        open={isModalOpen}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalOpen(false);
          setEditingFamily(null);
          form.resetFields();
        }}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="name"
            label="家庭名称"
            rules={[{ required: true, message: '请输入家庭名称' }]}
          >
            <Input placeholder="请输入家庭名称" />
          </Form.Item>
          <Form.Item
            name="invite_code"
            label="邀请码"
            rules={[{ required: true, message: '请输入邀请码' }]}
          >
            <Input placeholder="请输入邀请码" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="添加家庭成员"
        open={isMemberModalOpen}
        onOk={handleAddMemberSubmit}
        onCancel={() => {
          setIsMemberModalOpen(false);
          memberForm.resetFields();
        }}
        confirmLoading={addMemberMutation.isPending}
      >
        <Form form={memberForm} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="user_id"
            label="选择用户"
            rules={[{ required: true, message: '请选择用户' }]}
          >
            <Select
              placeholder="搜索并选择用户"
              options={userOptions}
              showSearch
              optionFilterProp="label"
              loading={!usersData}
            />
          </Form.Item>
          <Form.Item
            name="role"
            label="角色"
            rules={[{ required: true, message: '请选择角色' }]}
          >
            <Select
              placeholder="请选择角色"
              options={[
                { label: '管理员', value: 'admin' },
                { label: '成员', value: 'member' },
                { label: '长辈', value: 'elder' },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
