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
  Image,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { DeleteOutlined } from '@ant-design/icons';
import {
  getPhotos,
  getPhotoDetail,
  deletePhoto,
} from '@/admin/api/adminApi';
import type { AdminPhoto } from '@/admin/types/admin';
import dayjs from 'dayjs';

export function PhotosPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [detailId, setDetailId] = useState<string | null>(null);

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
    </div>
  );
}
