# Bug #54 分析报告 — 照片管理下载保存为 undefined.txt

## 1. Bug 根因分析

### 后端 (yuanzi-backend/handler/admin/photo.go)
`GetPhotos` handler 中的 `photoItem` 结构体缺少两个关键字段：
- **`Filename`** — 照片的原始文件名，前端将其用作 `download` 属性值
- **`OriginalURL`** — 照片的 R2 公开访问 URL，前端通过 axios 请求下载

当前 `photoItem`:
```go
type photoItem struct {
    ID          string `json:"id"`
    BabyID      string `json:"baby_id"`
    FamilyID    string `json:"family_id"`
    Size        int64  `json:"size"`
    ContentType string `json:"content_type"`
    Status      string `json:"status"`
    UploadedAt  string `json:"uploaded_at"`
    // ❌ 缺少: filename / original_url
}
```

### 前端 (yuanzi-frontend/src/admin/pages/PhotosPage.tsx:187)
```tsx
a.download = photo.filename;  // undefined → 浏览器保存为 "undefined.txt"
```
同时 axios 请求 URL `photo.original_url` 也是 undefined。

### 数据流
```
R2 存储 → OSSKey (path/to/file.jpg) → 后端 API → 前端 download
                                    ↑
                    filename/original_url 未从 OSSKey 派生
```

## 2. 影响范围
- **受影响功能**: 管理后台 → 照片管理 → 下载按钮
- **受影响用户**: 所有使用管理后台下载照片的用户
- **严重程度**: High — 核心功能不可用

## 3. 修复方案

### 后端修改 (handler/admin/photo.go)
1. `photoItem` 结构体新增字段:
   - `Filename string json:"filename"` — 从 `OSSKey` 提取文件名
   - `OriginalURL string json:"original_url"` — 生成 R2 公开访问 URL
2. 新增辅助函数 `extractFilename(ossKey string) string` — 取最后一个 `/` 之后的文件名
3. `GetPhotos` 遍历循环中填充 `Filename` 和 `OriginalURL`

### 前端修改 (pages/PhotosPage.tsx)
1. `handleDownloadSingle` — 使用 `photo.filename` 作为下载文件名
2. 兜底: 如果 filename 为空，从 original_url 路径提取文件名
3. 下载 blob 时设置正确的 MIME type (从 `photo.content_type`)

## 4. 测试计划

### 后端测试
- [ ] `TestExtractFilename` — 各种 OSSKey 格式的边界测试
- [ ] `TestGetPhotosFilename` — 验证 GetPhotos 返回 filename 字段
- [ ] `TestGetPhotosOriginalURL` — 验证 GetPhotos 返回 original_url 字段

### 前端测试
- [ ] `handleDownloadSingle` — 有 filename 时正确设置 download 属性
- [ ] `handleDownloadSingle` — filename 为空时从 URL 兜底提取
- [ ] `handleDownloadSingle` — content_type 正确设置 MIME type

## 5. 文件变更清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `yuanzi-backend/handler/admin/photo.go` | 修改 | 新增 filename/original_url 字段 |
| `yuanzi-backend/handler/admin/photo_test.go` | 修改 | 新增测试用例 |
| `yuanzi-frontend/src/admin/pages/PhotosPage.tsx` | 修改 | 修复下载逻辑 + 兜底 |

---
*由 Bug Fix Pipeline Step 2 生成 | 2026-05-09*
