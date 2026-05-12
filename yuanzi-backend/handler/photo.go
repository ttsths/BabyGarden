package handler

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	"yuanzi-backend/config"
	"yuanzi-backend/logger"
	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/oss"
	"yuanzi-backend/pkg/storage"

	"github.com/gin-gonic/gin"
)

const (
	photoUploadExpireSeconds = 300
	photoThumbWidth          = 300
)

// getStorageProvider returns the configured storage Provider (OSS or R2).
func getStorageProvider() storage.Provider {
	provider, err := storage.NewProviderFromConfig()
	if err != nil {
		// Fallback to OSS if config resolution fails.
		return storage.NewOSSProvider()
	}
	return provider
}

// GetPhotoUploadURL 获取照片上传 URL
// @Summary 获取照片上传 URL
// @Description 获取阿里云 OSS 直传签名 URL
// @Tags 照片
// @Accept json
// @Produce json
// @Security Bearer
// @Param data body PhotoUploadURLRequest true "请求参数"
// @Success 200 {object} model.Response{data=PhotoUploadURLResponse}
// @Router /api/v1/photo/upload-url [post]
func GetPhotoUploadURL(c *gin.Context) {
	var req PhotoUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code: model.ERROR_INVALID,
			Msg:  "请求参数错误",
		})
		return
	}
	if req.Size <= 0 {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "文件大小错误"})
		return
	}

	baby, _, err := loadBabyForRecord(c, req.BabyID)
	if err != nil {
		return
	}

	var family model.Family
	if err := mysql.DB.Where("id = ?", baby.FamilyID).First(&family).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "获取家庭信息失败"})
		return
	}

	used, err := calculateFamilyStorageUsed(family.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "校验存储配额失败"})
		return
	}
	if used+req.Size > family.StorageLimit {
		c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_QUOTA_EXCEEDED, Msg: "存储配额不足"})
		return
	}

	filename := path.Base(req.Filename)
	objectKey := buildPhotoObjectKey(family.ID, baby.ID, filename)

	provider := getStorageProvider()
	sig, err := provider.GetUploadSignature(objectKey, req.Size, photoUploadExpireSeconds, storage.WithContentType(req.ContentType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "生成签名失败"})
		return
	}

	userID := middleware.GetUserIDOrZero(c)
	photo := model.Photo{
		BabyID:      baby.ID,
		FamilyID:    baby.FamilyID,
		OSSKey:      objectKey,
		Size:        req.Size,
		ContentType: req.ContentType,
		UploadedBy:  userID,
		UploadedAt:  time.Now(),
		Status:      model.PhotoStatusPending,
	}
	if err := mysql.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "创建照片记录失败"})
		return
	}

	formData := sig.FormData
	if formData == nil {
		formData = make(map[string]string)
	}
	if req.ContentType != "" {
		formData["Content-Type"] = req.ContentType
	}

	expiresAt := sig.ExpiresAt
	if expiresAt == 0 {
		expiresAt = time.Now().Add(photoUploadExpireSeconds * time.Second).Unix()
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: PhotoUploadURLResponse{
			PhotoID:   photo.ID,
			UploadURL: sig.UploadURL,
			AccessURL: sig.AccessURL,
			ThumbURL:  provider.GetThumbnailURL(objectKey, photoThumbWidth),
			ExpiresAt: expiresAt,
			FormData:  formData,
			Headers:   sig.Headers,
		},
	})
}

// PhotoUploadCallback 照片上传回调（OSS 服务端回调）
// @Summary 照片上传回调
// @Description OSS 上传完成后的回调处理
// @Tags 照片
// @Accept json
// @Produce json
// @Param data body PhotoCallbackRequest true "回调参数"
// @Success 200 {object} model.Response
// @Router /api/v1/photo/callback [post]
func PhotoUploadCallback(c *gin.Context) {
	callbackToken := c.GetHeader("X-OSS-Callback-Token")
	callbackSecret := config.GlobalConfig.OSS.CallbackSecret
	if callbackSecret == "" {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "回调鉴权未配置"})
		return
	}
	if err := oss.VerifyCallbackToken(callbackSecret, callbackToken); err != nil {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "回调鉴权失败"})
		return
	}

	var req PhotoCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code: model.ERROR_INVALID,
			Msg:  "请求参数错误",
		})
		return
	}
	if req.Size <= 0 {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "文件大小错误"})
		return
	}

	if err := confirmUploadedPhoto(req.PhotoID, req.Size); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "处理成功",
	})
}

// PhotoUploadConfirm 照片上传确认（R2 / 客户端主动确认）
// @Summary 照片上传确认
// @Description 客户端完成直传后确认上传（用于 R2 等无服务端回调的存储后端）
// @Tags 照片
// @Accept json
// @Produce json
// @Security Bearer
// @Param data body PhotoConfirmRequest true "请求参数"
// @Success 200 {object} model.Response
// @Router /api/v1/photo/confirm [post]
func PhotoUploadConfirm(c *gin.Context) {
	var req PhotoConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}
	if req.PhotoID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "照片ID不能为空"})
		return
	}

	var photo model.Photo
	if err := mysql.DB.Where("id = ?", req.PhotoID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "照片不存在"})
		return
	}

	// 权限校验：仅上传者本人或管理员可确认
	userID := middleware.GetUserIDOrZero(c)
	if photo.UploadedBy != userID {
		var member model.FamilyMember
		if err := mysql.DB.Where("family_id = ? AND user_id = ?", photo.FamilyID, userID).First(&member).Error; err != nil || !member.IsAdmin() {
			c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "无权确认"})
			return
		}
	}

	if photo.Status != model.PhotoStatusPending {
		c.JSON(http.StatusConflict, model.Response{Code: model.ERROR, Msg: "照片已确认或状态异常"})
		return
	}

	// 使用请求中提供的文件大小，若未提供则保持原有值
	confirmedSize := req.Size
	if confirmedSize <= 0 {
		confirmedSize = photo.Size
	}

	if err := confirmUploadedPhoto(req.PhotoID, confirmedSize); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "确认成功",
	})
}

// ListPhotos 获取照片列表
// @Summary 获取照片列表
// @Description 分页获取指定宝宝的照片列表
// @Tags 照片
// @Accept json
// @Produce json
// @Security Bearer
// @Param baby_id query string true "宝宝ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} model.Response{data=model.ListResponse{list=[]PhotoResponse}}
// @Router /api/v1/photo [get]
func ListPhotos(c *gin.Context) {
	babyID := c.Query("baby_id")
	if babyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "宝宝ID不能为空"})
		return
	}
	if _, _, err := loadBabyForRecord(c, babyID); err != nil {
		return
	}

	page := parsePage(c.DefaultQuery("page", "1"))
	pageSize := parsePageSize(c.DefaultQuery("page_size", "20"))
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	query := mysql.DB.Model(&model.Photo{}).Where("baby_id = ? AND status = ?", babyID, model.PhotoStatusActive)
	if dateFrom != "" {
		start, err := time.ParseInLocation("2006-01-02", dateFrom, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "日期格式错误"})
			return
		}
		query = query.Where("uploaded_at >= ?", start)
	}
	if dateTo != "" {
		end, err := time.ParseInLocation("2006-01-02", dateTo, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "日期格式错误"})
			return
		}
		end = end.Add(24 * time.Hour)
		query = query.Where("uploaded_at < ?", end)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	var photos []model.Photo
	if err := query.Order("taken_at desc, uploaded_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	provider := getStorageProvider()
	list := make([]PhotoResponse, 0, len(photos))
	for _, photo := range photos {
		list = append(list, PhotoResponse{
			ID:          photo.ID,
			URL:         provider.GetURL(photo.OSSKey),
			ThumbURL:    provider.GetThumbnailURL(photo.OSSKey, photoThumbWidth),
			Width:       derefInt(photo.Width),
			Height:      derefInt(photo.Height),
			TakenAt:     formatTime(photo.TakenAt),
			Description: photo.Description,
		})
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: model.ListResponse{
			List: list,
			Pagination: model.Pagination{
				Page:       page,
				PageSize:   pageSize,
				Total:      total,
				TotalPages: calcTotalPages(total, pageSize),
			},
		},
	})
}

// DeletePhoto 删除照片
// @Summary 删除照片
// @Description 删除指定照片
// @Tags 照片
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "照片ID"
// @Success 200 {object} model.Response
// @Router /api/v1/photo/{id} [delete]
func DeletePhoto(c *gin.Context) {
	photo, member, err := loadPhotoWithMemberAccess(c)
	if err != nil {
		return
	}
	userID := middleware.GetUserIDOrZero(c)
	if !member.IsAdmin() && photo.UploadedBy != userID {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "无权删除"})
		return
	}

	if err := mysql.DB.Model(&model.Photo{}).Where("id = ?", photo.ID).Update("status", model.PhotoStatusDeleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "删除失败"})
		return
	}

	go deletePhotoFromStorage(photo)

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "删除成功",
	})
}

// 请求响应结构

type PhotoUploadURLRequest struct {
	BabyID      string `json:"baby_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	Filename    string `json:"filename" binding:"required" example:"photo.jpg"`
	ContentType string `json:"content_type" binding:"required" example:"image/jpeg"`
	Size        int64  `json:"size" binding:"required" example:"2048000"`
}

type PhotoUploadURLResponse struct {
	PhotoID   string            `json:"photo_id"`
	UploadURL string            `json:"upload_url"`
	AccessURL string            `json:"access_url"`
	ThumbURL  string            `json:"thumb_url"`
	ExpiresAt int64             `json:"expires_at"`
	FormData  map[string]string `json:"form_data,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

type PhotoCallbackRequest struct {
	PhotoID string `json:"photo_id" binding:"required"`
	Size    int64  `json:"size" binding:"required"`
	ETag    string `json:"etag" binding:"required"`
}

type PhotoConfirmRequest struct {
	PhotoID string `json:"photo_id" binding:"required"`
	Size    int64  `json:"size"`
	ETag    string `json:"etag"`
}

type PhotoResponse struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	ThumbURL    string `json:"thumb_url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	TakenAt     string `json:"taken_at"`
	Description string `json:"description,omitempty"`
}

func loadPhotoWithMemberAccess(c *gin.Context) (*model.Photo, *model.FamilyMember, error) {
	photoID := c.Param("id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "照片ID错误"})
		return nil, nil, errors.New("invalid photo id")
	}

	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return nil, nil, errors.New("unauthorized")
	}

	var photo model.Photo
	if err := mysql.DB.Where("id = ?", photoID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "照片不存在"})
		return nil, nil, err
	}

	var member model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", photo.FamilyID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "无权访问"})
		return nil, nil, err
	}

	return &photo, &member, nil
}

func buildPhotoObjectKey(familyID, babyID, filename string) string {
	if filename == "." || filename == "/" {
		filename = "photo.jpg"
	}
	date := time.Now().Format("20060102")
	return "families/" + familyID + "/babies/" + babyID + "/" + date + "/" + model.NewID() + "_" + filename
}

func formatTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(time.RFC3339)
}

func derefInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

// confirmUploadedPhoto updates photo status from pending to active.
func confirmUploadedPhoto(photoID string, size int64) error {
	var photo model.Photo
	if err := mysql.DB.Where("id = ?", photoID).First(&photo).Error; err != nil {
		return fmt.Errorf("照片不存在")
	}

	updates := map[string]interface{}{
		"size":        size,
		"status":      model.PhotoStatusActive,
		"uploaded_at": time.Now(),
	}
	if err := mysql.DB.Model(&model.Photo{}).Where("id = ?", photo.ID).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新照片失败")
	}

	go submitPhotoThumbnailTask(photo.ID)
	return nil
}

func submitPhotoThumbnailTask(photoID string) {
	if photoID == "" {
		return
	}
	logger.Info("提交缩略图生成任务", logger.String("photo_id", photoID))
}

func deletePhotoFromStorage(photo *model.Photo) {
	if photo == nil || photo.OSSKey == "" {
		return
	}
	provider := getStorageProvider()
	if err := provider.DeleteObject(photo.OSSKey); err != nil {
		logger.Warn("删除存储文件失败", logger.Err(err), logger.String("key", photo.OSSKey))
	}
	if photo.ThumbnailKey != "" && photo.ThumbnailKey != photo.OSSKey {
		if err := provider.DeleteObject(photo.ThumbnailKey); err != nil {
			logger.Warn("删除存储缩略图失败", logger.Err(err), logger.String("key", photo.ThumbnailKey))
		}
	}
}
