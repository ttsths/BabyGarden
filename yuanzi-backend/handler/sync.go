package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/gredis"

	"github.com/gin-gonic/gin"
)

const (
	syncHeartbeatInterval = 30 * time.Second
)

// SSEStream SSE 实时同步流
// @Summary SSE 实时同步
// @Description 建立服务器发送事件连接，接收家庭数据变更
// @Tags 同步
// @Accept json
// @Produce text/event-stream
// @Security Bearer
// @Success 200 "SSE 连接建立"
// @Router /api/v1/sync/stream [get]
func SSEStream(c *gin.Context) {
	// 设置 SSE 响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未认证"})
		return
	}
	familyID := c.Query("family_id")
	if familyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "家庭ID不能为空"})
		return
	}

	var member model.FamilyMember
	if err := mysql.DB.Where("family_id = ? AND user_id = ?", familyID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, model.Response{Code: model.ERROR_FORBID, Msg: "非家庭成员无法访问"})
		return
	}

	// 发送初始化事件
	initData := map[string]interface{}{
		"family_id": familyID,
		"timestamp": time.Now().Unix(),
	}
	writeSSEEvent(c.Writer, "init", initData)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "仅支持 HTTP/1.1"})
		return
	}
	flusher.Flush()

	ticker := time.NewTicker(syncHeartbeatInterval)
	defer ticker.Stop()

	pubsub := gredis.Subscribe(syncChannel(familyID))
	if pubsub == nil {
		// Redis 未连接，只用心跳保持连接
		for {
			select {
			case <-ticker.C:
				writeSSEEvent(c.Writer, "ping", nil)
				flusher.Flush()
			case <-clientGone:
				return
			}
		}
	}
	defer pubsub.Close()
	messages := pubsub.Channel()

	clientGone := c.Request.Context().Done()

	for {
		select {
		case msg := <-messages:
			if msg == nil {
				continue
			}
			writeSSEEventRaw(c.Writer, "message", msg.Payload)
			flusher.Flush()
		case <-ticker.C:
			writeSSEEvent(c.Writer, "ping", nil)
			flusher.Flush()
		case <-clientGone:
			return
		}
	}
}

// publishSyncEvent 发布家庭同步事件
func publishSyncEvent(familyID, event string, payload interface{}) {
	if familyID == "" {
		return
	}
	data := map[string]interface{}{
		"event":     event,
		"timestamp": time.Now().Unix(),
		"data":      payload,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	_ = gredis.Publish(syncChannel(familyID), string(b))
}

func syncChannel(familyID string) string {
	return fmt.Sprintf("sync:family:%s", familyID)
}

// writeSSEEvent 写入 SSE 事件
func writeSSEEvent(w http.ResponseWriter, event string, data interface{}) {
	if data != nil {
		w.Write([]byte("event: " + event + "\n"))
		w.Write([]byte("data: " + toJSON(data) + "\n\n"))
		return
	}
	w.Write([]byte(": ping\n\n"))
}

func writeSSEEventRaw(w http.ResponseWriter, event string, data string) {
	w.Write([]byte("event: " + event + "\n"))
	w.Write([]byte("data: " + data + "\n\n"))
}

// toJSON 转换为 JSON
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
