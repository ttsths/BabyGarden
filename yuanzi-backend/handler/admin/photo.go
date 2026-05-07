package admin

import (
	"net/http"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"

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
		ID          string `json:"id"`
		BabyID      string `json:"baby_id"`
		FamilyID    string `json:"family_id"`
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
		Status      string `json:"status"`
		UploadedAt  string `json:"uploaded_at"`
	}

	items := make([]photoItem, len(photos))
	for i, p := range photos {
		items[i] = photoItem{
			ID:          p.ID,
			BabyID:      p.BabyID,
			FamilyID:    p.FamilyID,
			Size:        p.Size,
			ContentType: p.ContentType,
			Status:      string(p.Status),
			UploadedAt:  p.UploadedAt.Format("2006-01-02 15:04:05"),
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

// GetPhoto returns photo detail.
// GET /api/v1/admin/photos/:id
func GetPhoto(c *gin.Context) {
	id := c.Param("id")
	var photo model.Photo
	if err := mysql.DB.Where("id = ?", id).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "照片不存在"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "获取成功",
		Data: gin.H{
			"id":            photo.ID,
			"baby_id":       photo.BabyID,
			"family_id":     photo.FamilyID,
			"oss_key":       photo.OSSKey,
			"thumbnail_key": photo.ThumbnailKey,
			"width":         photo.Width,
			"height":        photo.Height,
			"size":          photo.Size,
			"content_type":  photo.ContentType,
			"uploaded_by":   photo.UploadedBy,
			"uploaded_at":   photo.UploadedAt.Format("2006-01-02 15:04:05"),
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
