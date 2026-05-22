# Bug #62 Analysis — R2 上传 Network Error + 缩略图展示

## 问题概述
后台照片上传到 R2 时，前端 PUT 到 presigned URL 报 `Network Error`，详情页无法展示缩略图。

## 根因分析

### 根因 1：后端 presigned URL 未绑定 Content-Type（致命）
**位置**: `yuanzi-backend/pkg/storage/r2_provider.go:88-96`

```go
req, err := p.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
    Bucket: aws.String(p.bucket),
    Key:    aws.String(key),
}, s3.WithPresignExpires(...))
```

- `PresignPutObject` 时**没有设置 `ContentType`**
- 前端上传时显式设置了 `Content-Type: image/jpeg`（`PhotosPage.tsx:173-176`）
- S3/R2 签名验证时，若请求携带了签名时未包含的 header，签名不匹配 → 403 Forbidden
- 浏览器中 403 跨域请求显示为 `Network Error`

### 根因 2：前端 axios 可能携带额外 headers
**位置**: `yuanzi-frontend/src/admin/pages/PhotosPage.tsx:170-180`

```ts
await axios.put(upload_url, file, {
  headers: {
    'Content-Type': file.type || 'image/jpeg',
  },
  ...
});
```

- 使用的是原始 `axios` 实例（非 `adminClient`），无 `Authorization` interceptor ✓
- 但 axios 默认行为可能发送 `X-Requested-With` 等 header，触发 CORS preflight
- 需要显式控制 headers，确保不携带额外 header 到 R2 域名

### 根因 3：缩略图字段缺失 + R2 不支持 OSS 图片处理参数
**位置**: 
- 后端 `yuanzi-backend/handler/admin/photo.go:43-70` — `photoItem` 缺少 `thumbnail_url`
- 后端 `yuanzi-backend/pkg/storage/r2_provider.go:115-120` — `GetThumbnailURL` 使用 `x-oss-process=image/resize,w_%d`（阿里 OSS 语法）

- 后端返回列表时只有 `original_url`，但前端表格期望 `thumbnail_url` → 缩略图列永远显示"无"
- R2 不支持 `x-oss-process` 参数，缩略图应直接用 `original_url` 或另行处理

## 修复方案

### 后端修复（r2_provider.go）
1. `GetUploadSignature` 签名时传入 `ContentType`，与前端保持一致
2. `GetThumbnailURL` 去掉 OSS 专属参数，直接返回原图 URL（R2 暂无图片处理）

### 后端修复（admin/photo.go）
1. `photoItem` 增加 `thumbnail_url` 字段，与 `original_url` 一致
2. 详情接口返回增加 `original_url` 字段（当前只返回了 `oss_key`）

### 前端修复（PhotosPage.tsx）
1. R2 PUT 上传时显式设置 `{ transformRequest: [(data, headers) => { delete headers['X-Requested-With']; return data; }] }` 避免额外 headers
2. 或使用 `fetch` 替代 `axios.put` 到 R2，更干净

## 验收标准
- [ ] 照片上传成功，Network 面板无 CORS/签名错误
- [ ] 列表页缩略图正常展示
- [ ] 详情页图片预览正常
- [ ] 后端编译通过 + 测试通过
- [ ] 前端类型检查通过
