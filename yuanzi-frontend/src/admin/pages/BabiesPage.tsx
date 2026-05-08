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
  Form,
  Modal,
  Input,
  Select,
  DatePicker,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { DeleteOutlined, PlusOutlined, EditOutlined } from '@ant-design/icons';
import {
  getBabies,
  getBabyDetail,
  deleteBaby,
  createBaby,
  updateBaby,
} from '@/admin/api/adminApi';
import type { AdminBaby, CreateBabyRequest, UpdateBabyRequest } from '@/admin/types/admin';
import dayjs from 'dayjs';

export function BabiesPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [detailId, setDetailId] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingBaby, setEditingBaby] = useState<AdminBaby | null>(null);
  const [form] = Form.useForm();

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

  const createMutation = useMutation({
    mutationFn: (data: CreateBabyRequest) => createBaby(data),
    onSuccess: () => {
      message.success('创建成功');
      setIsModalOpen(false);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['admin', 'babies'] });
    },
    onError: () => {
      message.error('创建失败');
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateBabyRequest }) => updateBaby(id, data),
    onSuccess: () => {
      message.success('更新成功');
      setIsModalOpen(false);
      setEditingBaby(null);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['admin', 'babies'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'baby', detailId] });
    },
    onError: () => {
      message.error('更新失败');
    },
  });

  const handleOpenCreate = () => {
    setEditingBaby(null);
    form.resetFields();
    setIsModalOpen(true);
  };

  const handleOpenEdit = (record: AdminBaby) => {
    setEditingBaby(record);
    form.setFieldsValue({
      name: record.name,
      gender: record.gender,
      birthday: record.birthday ? dayjs(record.birthday) : null,
      family_id: record.family_id,
    });
    setIsModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const payload = {
        ...values,
        birthday: values.birthday ? values.birthday.format('YYYY-MM-DD') : undefined,
      };
      if (editingBaby) {
        updateMutation.mutate({ id: editingBaby.id, data: payload });
      } else {
        createMutation.mutate(payload);
      }
    } catch {
      // validation failed
    }
  };

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
          <Button type="link" icon={<EditOutlined />} onClick={() => handleOpenEdit(record)}>
            编辑
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
          新增宝宝
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

      <Modal
        title={editingBaby ? '编辑宝宝' : '新增宝宝'}
        open={isModalOpen}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalOpen(false);
          setEditingBaby(null);
          form.resetFields();
        }}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="name"
            label="名字"
            rules={[{ required: true, message: '请输入名字' }]}
          >
            <Input placeholder="请输入名字" />
          </Form.Item>
          <Form.Item
            name="gender"
            label="性别"
            rules={[{ required: true, message: '请选择性别' }]}
          >
            <Select
              placeholder="请选择性别"
              options={[
                { label: '男', value: 'male' },
                { label: '女', value: 'female' },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="birthday"
            label="生日"
            rules={[{ required: true, message: '请选择生日' }]}
          >
            <DatePicker style={{ width: '100%' }} placeholder="请选择生日" />
          </Form.Item>
          <Form.Item
            name="family_id"
            label="家庭ID"
            rules={[{ required: true, message: '请输入家庭ID' }]}
          >
            <Input placeholder="请输入家庭ID" />
          </Form.Item>
          <Form.Item
            name="avatar_url"
            label="头像URL"
          >
            <Input placeholder="请输入头像URL" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
