import { useState, useCallback } from 'react';
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
  Image,
  Upload,
  Modal,
  Select,
  Progress,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { DeleteOutlined, UploadOutlined, DownloadOutlined } from '@ant-design/icons';
import {
  getPhotos,
  getPhotoDetail,
  deletePhoto,
  getPhotoUploadUrl,
  confirmPhotoUpload,
} from '@/admin/api/adminApi';
import type { AdminPhoto } from '@/admin/types/admin';
import dayjs from 'dayjs';
import axios from 'axios';

interface UploadTask {
  id: string;
  file: File;
  status: 'pending' | 'uploading' | 'confirming' | 'done' | 'error';
  progress: number;
  error?: string;
}

export function PhotosPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [detailId, setDetailId] = useState<string | null>(null);
  const [uploadModalOpen, setUploadModalOpen] = useState(false);
  const [uploadFamilyId, setUploadFamilyId] = useState<string>('');
  const [uploadBabyId, setUploadBabyId] = useState<string>('');
  const [uploadTasks, setUploadTasks] = useState<UploadTask[]>([]);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [downloadProgress, setDownloadProgress] = useState<Record<string, number>>({});

  const { data, isLoading, error } = useQuery({
    queryKey: ['admin', 'photos', page, pageSize],
    queryFn: async () => {
      const res = await getPhotos(page, pageSize);
      return res.data.data;
    },
  });

  const { data: detailData, isLoading: detailLoading } = useQuery({
    queryKey: ['admin', 'photo', detailId],
    queryFn: async () => {
      if (!detailId) return null;
      const res = await getPhotoDetail(detailId);
      return res.data.data;
    },
    enabled: !!detailId,
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deletePhoto(id),
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['admin', 'photos'] });
    },
    onError: () => {
      message.error('删除失败');
    },
  });

  const formatSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
    return `${(bytes / 1024 / 1024).toFixed(2)} MB`;
  };

  const updateTask = useCallback((id: string, updates: Partial<UploadTask>) => {
    setUploadTasks((prev) =>
      prev.map((t) => (t.id === id ? { ...t, ...updates } : t))
    );
  }, []);

  const handleUpload = useCallback(
    async (file: File) => {
      const taskId = `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
      setUploadTasks((prev) => [
        ...prev,
        { id: taskId, file, status: 'pending', progress: 0 },
      ]);

      try {
        if (!uploadFamilyId) {
          throw new Error('请选择家庭');
        }

        updateTask(taskId, { status: 'uploading', progress: 10 });

        const urlRes = await getPhotoUploadUrl({
          filename: file.name,
          content_type: file.type || 'image/jpeg',
        });
        const { upload_url, download_url } = urlRes.data.data;

        updateTask(taskId, { progress: 40 });

        await axios.put(upload_url, file, {
          headers: {
            'Content-Type': file.type || 'image/jpeg',
          },
          onUploadProgress: (e) => {
            if (e.total) {
              const percent = Math.round(40 + (e.loaded / e.total) * 40);
              updateTask(taskId, { progress: percent });
            }
          },
        });

        updateTask(taskId, { status: 'confirming', progress: 85 });

        await confirmPhotoUpload({
          filename: file.name,
          url: download_url,
          family_id: uploadFamilyId,
          baby_id: uploadBabyId || undefined,
        });

        updateTask(taskId, { status: 'done', progress: 100 });
        message.success(`${file.name} 上传成功`);
        queryClient.invalidateQueries({ queryKey: ['admin', 'photos'] });
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : '上传失败';
        updateTask(taskId, { status: 'error', error: msg });
        message.error(`${file.name} 上传失败: ${msg}`);
      }
    },
    [uploadFamilyId, uploadBabyId, updateTask, queryClient]
  );

  const handleDownloadSingle = useCallback(async (photo: AdminPhoto) => {
    try {
      setDownloadProgress((prev) => ({ ...prev, [photo.id]: 0 }));
      const res = await axios.get(photo.original_url, {
        responseType: 'blob',
        onDownloadProgress: (e) => {
          if (e.total) {
            const percent = Math.round((e.loaded / e.total) * 100);
            setDownloadProgress((prev) => ({ ...prev, [photo.id]: percent }));
          }
        },
      });
      const blob = new Blob([res.data]);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = photo.filename;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
      setDownloadProgress((prev) => ({ ...prev, [photo.id]: 100 }));
      message.success('下载成功');
    } catch {
      message.error('下载失败');
    }
  }, []);

  const handleBatchDownload = useCallback(async () => {
    const selectedPhotos =
      data?.list.filter((p) => selectedRowKeys.includes(p.id)) || [];
    if (selectedPhotos.length === 0) {
      message.warning('请先选择照片');
      return;
    }
    if (selectedPhotos.length === 1) {
      await handleDownloadSingle(selectedPhotos[0]);
      return;
    }

    message.warning('批量下载功能开发中');
  }, [data, selectedRowKeys, handleDownloadSingle]);

  const columns: ColumnsType<AdminPhoto> = [
    {
      title: '缩略图',
      dataIndex: 'thumbnail_url',
      key: 'thumbnail_url',
      render: (url: string) =>
        url ? (
          <Image
            src={url}
            alt="thumb"
            width={64}
            height={64}
            style={{ objectFit: 'cover', borderRadius: 4 }}
            preview={false}
          />
        ) : (
          <span>无</span>
        ),
    },
    {
      title: '文件名',
      dataIndex: 'filename',
      key: 'filename',
    },
    {
      title: '大小',
      dataIndex: 'size',
      key: 'size',
      render: (size: number) => formatSize(size),
    },
    {
      title: '上传者ID',
      dataIndex: 'uploader_id',
      key: 'uploader_id',
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
          <Button
            type="link"
            icon={<DownloadOutlined />}
            onClick={() => handleDownloadSingle(record)}
            loading={downloadProgress[record.id] > 0 && downloadProgress[record.id] < 100}
          >
            {downloadProgress[record.id] > 0 && downloadProgress[record.id] < 100
              ? `${downloadProgress[record.id]}%`
              : '下载'}
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
      <Space style={{ marginBottom: 16 }}>
        <Button
          type="primary"
          icon={<UploadOutlined />}
          onClick={() => {
            setUploadTasks([]);
            setUploadModalOpen(true);
          }}
        >
          上传照片
        </Button>
        <Button
          icon={<DownloadOutlined />}
          onClick={handleBatchDownload}
          disabled={selectedRowKeys.length === 0}
        >
          批量下载 ({selectedRowKeys.length})
        </Button>
      </Space>

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
        rowSelection={{
          selectedRowKeys,
          onChange: (keys) => setSelectedRowKeys(keys),
        }}
      />

      <Drawer
        title="照片详情"
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
              <div className="text-gray-500 text-sm">文件名</div>
              <div>{detailData.filename}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">图片预览</div>
              {detailData.original_url ? (
                <Image
                  src={detailData.original_url}
                  alt="photo"
                  style={{ maxWidth: '100%', borderRadius: 8 }}
                />
              ) : (
                <span>无</span>
              )}
            </div>
            <div>
              <div className="text-gray-500 text-sm">大小</div>
              <div>{formatSize(detailData.size)}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">Content-Type</div>
              <div>{detailData.content_type}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">家庭ID</div>
              <div>{detailData.family_id}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">宝宝ID</div>
              <div>{detailData.baby_id}</div>
            </div>
            <div>
              <div className="text-gray-500 text-sm">上传者ID</div>
              <div>{detailData.uploader_id}</div>
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
        title="上传照片"
        open={uploadModalOpen}
        onCancel={() => setUploadModalOpen(false)}
        footer={null}
        width={640}
      >
        <Space direction="vertical" style={{ width: '100%' }} size="large">
          <div>
            <div style={{ marginBottom: 8 }}>
              <span style={{ color: '#ff4d4f', marginRight: 4 }}>*</span>
              目标家庭
            </div>
            <Select
              style={{ width: '100%' }}
              placeholder="选择家庭"
              value={uploadFamilyId || undefined}
              onChange={(v) => {
                setUploadFamilyId(v);
                setUploadBabyId('');
              }}
              options={[]}
              showSearch
            />
          </div>

          <div>
            <div style={{ marginBottom: 8 }}>目标宝宝（可选）</div>
            <Select
              style={{ width: '100%' }}
              placeholder="选择宝宝"
              value={uploadBabyId || undefined}
              onChange={setUploadBabyId}
              options={[]}
              showSearch
              disabled={!uploadFamilyId}
            />
          </div>

          <Upload.Dragger
            multiple
            showUploadList={false}
            beforeUpload={(file) => {
              handleUpload(file);
              return false;
            }}
            accept="image/*"
            disabled={!uploadFamilyId}
          >
            <p className="ant-upload-drag-icon">
              <UploadOutlined />
            </p>
            <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
            <p className="ant-upload-hint">支持批量上传，仅支持图片文件</p>
          </Upload.Dragger>

          {uploadTasks.length > 0 && (
            <div style={{ maxHeight: 300, overflowY: 'auto' }}>
              {uploadTasks.map((task) => (
                <div key={task.id} style={{ marginBottom: 12 }}>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      marginBottom: 4,
                    }}
                  >
                    <span
                      style={{
                        maxWidth: '70%',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                      }}
                      title={task.file.name}
                    >
                      {task.file.name}
                    </span>
                    <span
                      style={{
                        color:
                          task.status === 'done'
                            ? '#52c41a'
                            : task.status === 'error'
                            ? '#ff4d4f'
                            : '#1890ff',
                      }}
                    >
                      {task.status === 'pending' && '等待中'}
                      {task.status === 'uploading' && '上传中'}
                      {task.status === 'confirming' && '确认中'}
                      {task.status === 'done' && '完成'}
                      {task.status === 'error' && '失败'}
                    </span>
                  </div>
                  <Progress
                    percent={task.progress}
                    status={
                      task.status === 'error'
                        ? 'exception'
                        : task.status === 'done'
                        ? 'success'
                        : 'active'
                    }
                    size="small"
                  />
                  {task.error && (
                    <div style={{ color: '#ff4d4f', fontSize: 12 }}>
                      {task.error}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </Space>
      </Modal>
    </div>
  );
}
