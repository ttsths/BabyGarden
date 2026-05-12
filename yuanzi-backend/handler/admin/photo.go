package admin

import (
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/storage"

	"github.com/gin-gonic/gin"
)

// GetPhotos returns paginated photo list.
// GET /api/v1/admin/photos
func GetPhotos(c *gin.Context) {
	page, pageSize := paginate(c)
	var total int64
	var photos []model.Photo

	query := mysql.DB.Model(&model.Photo{})
	query.Count(&total)

	if err := query.Order("uploaded_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	type photoItem struct {
		ID           string `json:"id"`
		BabyID       string `json:"baby_id"`
		FamilyID     string `json:"family_id"`
		UploaderID   string `json:"uploader_id"`
		Filename     string `json:"filename"`
		OriginalURL  string `json:"original_url"`
		ThumbnailURL string `json:"thumbnail_url"`
		Size         int64  `json:"size"`
		ContentType  string `json:"content_type"`
		Status       string `json:"status"`
		CreatedAt    string `json:"created_at"`
		UploadedAt   string `json:"uploaded_at"`
	}

	provider, _ := storage.NewProviderFromConfig()

	items := make([]photoItem, len(photos))
	for i, p := range photos {
		// Extract filename from OSSKey (last segment after "/")
		filename := p.OSSKey
		if idx := strings.LastIndex(p.OSSKey, "/"); idx >= 0 {
			filename = p.OSSKey[idx+1:]
		}

		originalURL, thumbnailURL := adminPhotoURLs(provider, p)
		createdAt := p.UploadedAt.Format("2006-01-02 15:04:05")
		items[i] = photoItem{
			ID:           p.ID,
			BabyID:       p.BabyID,
			FamilyID:     p.FamilyID,
			UploaderID:   p.UploadedBy,
			Filename:     filename,
			OriginalURL:  originalURL,
			ThumbnailURL: thumbnailURL,
			Size:         p.Size,
			ContentType:  p.ContentType,
			Status:       string(p.Status),
			CreatedAt:    createdAt,
			UploadedAt:   createdAt,
		}
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: model.ListResponse{
			List: items,
			Pagination: model.Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	})
}

// GetPhotoUploadURL admin 获取照片上传 URL
// POST /api/v1/admin/photos/upload-url
func GetPhotoUploadURL(c *gin.Context) {
	var req struct {
		BabyID      string `json:"baby_id" binding:"required"`
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type" binding:"required"`
		Size        int64  `json:"size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	var baby model.Baby
	if err := mysql.DB.Where("id = ?", req.BabyID).First(&baby).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "宝宝不存在"})
		return
	}

	provider, err := storage.NewProviderFromConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "存储配置错误"})
		return
	}
	filename := path.Base(req.Filename)
	if filename == "" || filename == "." {
		filename = "photo.jpg"
	}
	objectKey := buildPhotoObjectKey(baby.FamilyID, baby.ID, filename)

	sig, err := provider.GetUploadSignature(objectKey, req.Size, 300, storage.WithContentType(req.ContentType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "获取上传签名失败"})
		return
	}

	photo := model.Photo{
		BabyID:      baby.ID,
		FamilyID:    baby.FamilyID,
		OSSKey:      objectKey,
		Size:        req.Size,
		ContentType: req.ContentType,
		UploadedBy:  "admin",
		UploadedAt:  time.Now(),
		Status:      model.PhotoStatusPending,
	}
	if err := mysql.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "创建照片记录失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: gin.H{
			"photo_id":       photo.ID,
			"upload_url":     sig.UploadURL,
			"access_url":     sig.AccessURL,
			"thumbnail_url":  sig.ThumbnailURL,
			"upload_headers": sig.Headers,
		},
	})
}

// PhotoUploadConfirm admin 确认照片上传
// POST /api/v1/admin/photos/confirm
func PhotoUploadConfirm(c *gin.Context) {
	var req struct {
		PhotoID string `json:"photo_id" binding:"required"`
		Size    int64  `json:"size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	var photo model.Photo
	if err := mysql.DB.Where("id = ?", req.PhotoID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "照片不存在"})
		return
	}

	if photo.Status != model.PhotoStatusPending {
		c.JSON(http.StatusConflict, model.Response{Code: model.ERROR, Msg: "照片已确认或状态异常"})
		return
	}

	confirmedSize := req.Size
	if confirmedSize <= 0 {
		confirmedSize = photo.Size
	}

	updates := map[string]interface{}{
		"status": model.PhotoStatusActive,
		"size":   confirmedSize,
	}
	if err := mysql.DB.Model(&photo).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "确认失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "确认成功",
	})
}

func buildPhotoObjectKey(familyID, babyID, filename string) string {
	return familyID + "/" + babyID + "/" + filename
}

func adminPhotoURLs(provider storage.Provider, photo model.Photo) (string, string) {
	if photo.OSSKey == "" {
		return "", ""
	}
	if provider != nil {
		thumbKey := photo.OSSKey
		if photo.ThumbnailKey != "" {
			thumbKey = photo.ThumbnailKey
		}
		return provider.GetURL(photo.OSSKey), provider.GetURL(thumbKey)
	}

	publicURL := strings.TrimRight(firstNonEmpty(
		os.Getenv("R2_PUBLIC_URL"),
		os.Getenv("R2_CUSTOM_DOMAIN"),
		storage.GetWorkerVar("R2_PUBLIC_URL"),
		storage.GetWorkerVar("R2_CUSTOM_DOMAIN"),
	), "/")
	if publicURL == "" {
		return "", ""
	}
	originalURL := publicURL + "/" + photo.OSSKey
	thumbnailKey := photo.OSSKey
	if photo.ThumbnailKey != "" {
		thumbnailKey = photo.ThumbnailKey
	}
	return originalURL, publicURL + "/" + thumbnailKey
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// GetPhoto returns photo detail.
// GET /api/v1/admin/photos/:id
func GetPhoto(c *gin.Context) {
	id := c.Param("id")
	var photo model.Photo
	if err := mysql.DB.Where("id = ?", id).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "照片不存在"})
		return
	}

	provider, _ := storage.NewProviderFromConfig()
	originalURL, thumbnailURL := adminPhotoURLs(provider, photo)
	filename := photo.OSSKey
	if idx := strings.LastIndex(photo.OSSKey, "/"); idx >= 0 {
		filename = photo.OSSKey[idx+1:]
	}
	createdAt := photo.UploadedAt.Format("2006-01-02 15:04:05")

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: gin.H{
			"id":            photo.ID,
			"baby_id":       photo.BabyID,
			"family_id":     photo.FamilyID,
			"oss_key":       photo.OSSKey,
			"filename":      filename,
			"original_url":  originalURL,
			"thumbnail_key": photo.ThumbnailKey,
			"thumbnail_url": thumbnailURL,
			"width":         photo.Width,
			"height":        photo.Height,
			"size":          photo.Size,
			"content_type":  photo.ContentType,
			"uploader_id":   photo.UploadedBy,
			"uploaded_by":   photo.UploadedBy,
			"created_at":    createdAt,
			"uploaded_at":   createdAt,
			"status":        string(photo.Status),
		},
	})
}

// DeletePhoto marks photo as deleted.
// DELETE /api/v1/admin/photos/:id
func DeletePhoto(c *gin.Context) {
	id := c.Param("id")
	if err := mysql.DB.Model(&model.Photo{}).Where("id = ?", id).Update("status", model.PhotoStatusDeleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "操作失败"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "已删除"})
}
