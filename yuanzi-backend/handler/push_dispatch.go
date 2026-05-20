package handler

import (
	"context"
	"fmt"
	"time"

	"yuanzi-backend/config"
	"yuanzi-backend/logger"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/push"
)

const apnsRequestTimeout = 5 * time.Second

func dispatchRecordPush(event string, record *model.Record, content map[string]interface{}) {
	if record == nil {
		return
	}
	cfg := config.GlobalConfig.Push.APNs
	if !push.IsAPNsEnabled(cfg) {
		return
	}

	client, err := push.APNs()
	if err != nil {
		logger.Warn("APNs 客户端初始化失败", logger.Err(err))
		return
	}

	devices, err := loadFamilyIOSDevices(record.FamilyID)
	if err != nil {
		logger.Warn("查询 APNs 设备失败", logger.Err(err), logger.String("family_id", record.FamilyID))
		return
	}
	if len(devices) == 0 {
		return
	}

	payload := buildRecordAPNsPayload(event, record, content)
	for _, device := range devices {
		deviceToken := device.DeviceToken
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), apnsRequestTimeout)
			defer cancel()
			if err := client.Send(ctx, deviceToken, payload); err != nil {
				logger.Warn("APNs 推送失败", logger.Err(err), logger.String("device_token", deviceToken))
			}
		}()
	}
}

func loadFamilyIOSDevices(familyID string) ([]model.PushDevice, error) {
	var devices []model.PushDevice
	err := mysql.DB.Table("push_devices").
		Joins("JOIN family_members ON family_members.user_id = push_devices.user_id").
		Where("family_members.family_id = ? AND push_devices.platform = ? AND push_devices.is_active = 1", familyID, "ios").
		Find(&devices).Error
	return devices, err
}

func buildRecordAPNsPayload(event string, record *model.Record, content map[string]interface{}) map[string]interface{} {
	if content == nil {
		content = recordContentToMap(record.Content)
	}
	title, body := recordPushAlert(event, record)
	recordData := map[string]interface{}{
		"id":         record.ID,
		"baby_id":    record.BabyID,
		"family_id":  record.FamilyID,
		"type":       string(record.Type),
		"started_at": record.StartedAt.Format(time.RFC3339),
		"ended_at":   formatNullTime(record.EndedAt),
		"content":    content,
		"note":       record.Note,
	}
	return map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": map[string]interface{}{
				"title": title,
				"body":  body,
			},
			"sound": "default",
		},
		"event":     event,
		"timestamp": time.Now().Unix(),
		"data": map[string]interface{}{
			"record": recordData,
		},
	}
}

func recordPushAlert(event string, record *model.Record) (string, string) {
	action := "更新"
	switch event {
	case "record_created":
		action = "新增"
	case "record_deleted":
		action = "删除"
	}
	label := recordTypeLabel(record.Type)
	title := fmt.Sprintf("%s记录%s", label, action)
	body := fmt.Sprintf("宝宝记录已%s", action)
	return title, body
}

func recordTypeLabel(recordType model.RecordType) string {
	switch recordType {
	case model.RecordTypeFeeding:
		return "喂养"
	case model.RecordTypeSleep:
		return "睡眠"
	case model.RecordTypeDiaper:
		return "尿布"
	case model.RecordTypeExcretion:
		return "排泄"
	case model.RecordTypeTemperature:
		return "测温"
	case model.RecordTypeGrowth:
		return "成长"
	default:
		return string(recordType)
	}
}
