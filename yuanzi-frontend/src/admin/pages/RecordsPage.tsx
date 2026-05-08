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
  getRecords,
  getRecordDetail,
  deleteRecord,
  createRecord,
  updateRecord,
  getFamilies,
  getBabies,
} from '@/admin/api/adminApi';
import type { AdminRecord, AdminFamily, AdminBaby, CreateRecordRequest, UpdateRecordRequest } from '@/admin/types/admin';
import dayjs from 'dayjs';

export function RecordsPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [detailId, setDetailId] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingRecord, setEditingRecord] = useState<AdminRecord | null>(null);
  const [form] = Form.useForm();

  const { data: familiesData } = useQuery({
    queryKey: ['admin', 'families', 'all'],
    queryFn: async () => {
      const res = await getFamilies(1, 1000);
      return res.data.data?.list || [];
    },
  });

  const { data: babiesData } = useQuery({
    queryKey: ['admin', 'babies', 'all'],
    queryFn: async () => {
      const res = await getBabies(1, 1000);
      return res.data.data?.list || [];
    },
  });

  const familyMap = Object.fromEntries((familiesData || []).map((f: AdminFamily) => [f.id, f.name]));
  const babyMap = Object.fromEntries((babiesData || []).map((b: AdminBaby) => [b.id, b.name]));

  const familyOptions = (familiesData || []).map((f: AdminFamily) => ({ label: f.name, value: f.id }));
  const babyOptions = (babiesData || []).map((b: AdminBaby) => ({ label: b.name, value: b.id }));

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

  const createMutation = useMutation({
    mutationFn: (data: CreateRecordRequest) => createRecord(data),
    onSuccess: () => {
      message.success('创建成功');
      setIsModalOpen(false);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['admin', 'records'] });
    },
    onError: () => {
      message.error('创建失败');
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateRecordRequest }) => updateRecord(id, data),
    onSuccess: () => {
      message.success('更新成功');
      setIsModalOpen(false);
      setEditingRecord(null);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['admin', 'records'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'record', detailId] });
    },
    onError: () => {
      message.error('更新失败');
    },
  });

  const handleOpenCreate = () => {
    setEditingRecord(null);
    form.resetFields();
    setIsModalOpen(true);
  };

  const handleOpenEdit = (record: AdminRecord) => {
    setEditingRecord(record);
    form.setFieldsValue({
      type: record.type,
      baby_id: record.baby_id,
      family_id: record.family_id,
      note: record.note || '',
      started_at: record.started_at ? dayjs(record.started_at) : null,
      ended_at: record.data?.ended_at ? dayjs(record.data.ended_at as string) : null,
    });
    setIsModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const payload = {
        ...values,
        started_at: values.started_at ? values.started_at.format('YYYY-MM-DDTHH:mm:ssZ') : undefined,
        ended_at: values.ended_at ? values.ended_at.format('YYYY-MM-DDTHH:mm:ssZ') : undefined,
      };
      if (editingRecord) {
        updateMutation.mutate({ id: editingRecord.id, data: payload });
      } else {
        createMutation.mutate(payload);
      }
    } catch {
      // validation failed
    }
  };

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
      title: '宝宝',
      dataIndex: 'baby_id',
      key: 'baby_id',
      width: 120,
      render: (babyId: string) => babyMap[babyId] || babyId,
    },
    {
      title: '家庭',
      dataIndex: 'family_id',
      key: 'family_id',
      width: 120,
      render: (familyId: string) => familyMap[familyId] || familyId,
    },
    {
      title: '创建者',
      dataIndex: 'created_by',
      key: 'created_by',
      render: (createdBy: string) => createdBy,
    },
    {
      title: '结束时间',
      dataIndex: 'data',
      key: 'ended_at',
      render: (data: Record<string, unknown>) => data?.ended_at ? dayjs(data.ended_at as string).format('YYYY-MM-DD HH:mm') : '-',
    },
    {
      title: '开始时间',
      dataIndex: 'started_at',
      render: (value: string) => value ? dayjs(value).format('YYYY-MM-DD HH:mm') : '-',
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
          新增记录
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
              <div className="text-gray-500 text-sm">宝宝</div>
              <div>{babyMap[detailData.baby_id] || detailData.baby_id}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">家庭</div>
              <div>{familyMap[detailData.family_id] || detailData.family_id}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">创建者</div>
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

      <Modal
        title={editingRecord ? '编辑记录' : '新增记录'}
        open={isModalOpen}
        onOk={handleSubmit}
        onCancel={() => {
          setIsModalOpen(false);
          setEditingRecord(null);
          form.resetFields();
        }}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="type"
            label="记录类型"
            rules={[{ required: true, message: '请选择记录类型' }]}
          >
            <Select
              placeholder="请选择记录类型"
              options={[
                { label: '喂奶', value: 'feeding' },
                { label: '睡眠', value: 'sleep' },
                { label: '换尿布', value: 'diaper' },
                { label: '体温', value: 'temperature' },
                { label: '辅食', value: 'food' },
                { label: '用药', value: 'medicine' },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="baby_id"
            label="宝宝"
            rules={[{ required: true, message: '请选择宝宝' }]}
          >
            <Select placeholder="请选择宝宝" options={babyOptions} showSearch />
          </Form.Item>
          <Form.Item
            name="family_id"
            label="家庭"
            rules={[{ required: true, message: '请选择家庭' }]}
          >
            <Select placeholder="请选择家庭" options={familyOptions} showSearch />
          </Form.Item>
          <Form.Item
            name="started_at"
            label="开始时间"
            rules={[{ required: true, message: '请选择开始时间' }]}
          >
            <DatePicker showTime style={{ width: '100%' }} placeholder="请选择开始时间" />
          </Form.Item>
          <Form.Item
            noStyle
            shouldUpdate={(prev, curr) => prev.type !== curr.type}
          >
            {({ getFieldValue }) => {
              const type = getFieldValue('type');
              if (type !== 'sleep') return null;
              return (
                <Form.Item
                  name="ended_at"
                  label="结束时间"
                  rules={[{ required: true, message: '请选择结束时间' }]}
                >
                  <DatePicker showTime style={{ width: '100%' }} placeholder="请选择结束时间" />
                </Form.Item>
              );
            }}
          </Form.Item>
          <Form.Item
            name="note"
            label="备注"
          >
            <Input.TextArea placeholder="请输入备注" rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
