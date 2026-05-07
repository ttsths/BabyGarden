package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"yuanzi-backend/config"
	"yuanzi-backend/middleware"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/ai"
	"yuanzi-backend/pkg/gredis"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// AIChat AI 问答
// @Summary AI 问答
// @Description 发送问题给 AI 助手获取回答
// @Tags AI
// @Accept json
// @Produce json
// @Security Bearer
// @Param data body AIChatRequest true "问答请求"
// @Success 200 {object} model.Response{data=AIChatResponse}
// @Router /api/v1/ai/chat [post]
func AIChat(c *gin.Context) {
	var req AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code: model.ERROR_INVALID,
			Msg:  "请求参数错误",
		})
		return
	}

	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	used, remaining, err := consumeQuota(userID, quotaTypeChat, chatDailyLimit)
	if err != nil {
		if errors.Is(err, errQuotaExceeded) {
			c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_QUOTA_EXCEEDED, Msg: "AI配额已用完"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "配额检查失败"})
		return
	}
	_ = used

	var baby *model.Baby
	if req.BabyID != nil && *req.BabyID != "" {
		loaded, _, err := loadBabyForRecord(c, *req.BabyID)
		if err != nil {
			_ = rollbackQuota(userID, quotaTypeChat)
			return
		}
		baby = loaded
	}

	messages := buildChatMessages(req, baby)
	resp, err := aiChatFunc(messages)
	if err != nil {
		_ = rollbackQuota(userID, quotaTypeChat)
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "AI服务异常"})
		return
	}

	answer := strings.TrimSpace(resp.Output.Text)
	if answer == "" {
		answer = "AI暂时无法回答，请稍后再试"
	}

	record := model.AIChatRecord{
		UserID:     userID,
		BabyID:     req.BabyID,
		Question:   req.Question,
		Answer:     answer,
		TokensUsed: resp.Usage.TotalTokens,
		Model:      "qwen-turbo",
		CreatedAt:  time.Now(),
	}
	if err := mysql.DB.Create(&record).Error; err != nil {
		_ = rollbackQuota(userID, quotaTypeChat)
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "保存问答记录失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "success",
		Data: AIChatResponse{
			Answer:         answer,
			TokensUsed:     resp.Usage.TotalTokens,
			RemainingQuota: remaining,
		},
	})
}

// SpeechRecognize 语音识别
// @Summary 语音识别
// @Description 将语音转为文字
// @Tags AI
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param audio formData file true "音频文件"
// @Success 200 {object} model.Response{data=SpeechRecognizeResponse}
// @Router /api/v1/ai/speech/recognize [post]
func SpeechRecognize(c *gin.Context) {
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	fileHeader, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "音频文件不能为空"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "读取音频失败"})
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil || len(data) == 0 {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "音频文件错误"})
		return
	}

	_, remaining, err := consumeQuota(userID, quotaTypeSpeech, speechDailyLimit)
	if err != nil {
		if errors.Is(err, errQuotaExceeded) {
			c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_QUOTA_EXCEEDED, Msg: "语音配额已用完"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "配额检查失败"})
		return
	}

	result, err := speechRecognizeFunc(data)
	if err != nil {
		_ = rollbackQuota(userID, quotaTypeSpeech)
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "语音识别失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "success",
		Data: SpeechRecognizeResponse{
			Text:           result.Text,
			Confidence:     result.Confidence,
			RemainingQuota: remaining,
		},
	})
}

// GetAIQuota 获取 AI 配额
// @Summary 获取 AI 配额
// @Description 获取用户当日 AI 使用配额情况
// @Tags AI
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} model.Response{data=AIQuotaResponse}
// @Router /api/v1/ai/quota [get]
func GetAIQuota(c *gin.Context) {
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	speechUsed, _ := getQuotaUsage(userID, quotaTypeSpeech)
	chatUsed, _ := getQuotaUsage(userID, quotaTypeChat)

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "success",
		Data: AIQuotaResponse{
			Speech: QuotaDetail{
				Used:      speechUsed,
				Limit:     speechDailyLimit,
				Remaining: calcRemaining(speechDailyLimit, speechUsed),
			},
			AIChat: QuotaDetail{
				Used:      chatUsed,
				Limit:     chatDailyLimit,
				Remaining: calcRemaining(chatDailyLimit, chatUsed),
			},
		},
	})
}

// 请求响应结构

type AIChatRequest struct {
	Question string        `json:"question" binding:"required" example:"宝宝三个月应该喝多少奶"`
	BabyID   *string       `json:"baby_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
	History  []ChatHistory `json:"history,omitempty"`
}

type ChatHistory struct {
	Role    string `json:"role" example:"user"`
	Content string `json:"content" example:"之前的问题"`
}

type AIChatResponse struct {
	Answer         string `json:"answer"`
	TokensUsed     int    `json:"tokens_used"`
	RemainingQuota int    `json:"remaining_quota"`
}

type SpeechRecognizeResponse struct {
	Text           string  `json:"text"`
	Confidence     float64 `json:"confidence"`
	RemainingQuota int     `json:"remaining_quota"`
}

type AIQuotaResponse struct {
	Speech QuotaDetail `json:"speech"`
	AIChat QuotaDetail `json:"ai_chat"`
}

type QuotaDetail struct {
	Used      int `json:"used"`
	Limit     int `json:"limit"`
	Remaining int `json:"remaining"`
}

var (
	aiChatFunc          = defaultChatFunc
	speechRecognizeFunc = defaultSpeechFunc
)

const (
	quotaTypeChat    = "chat"
	quotaTypeSpeech  = "speech"
	chatDailyLimit   = 10
	speechDailyLimit = 20
)

var errQuotaExceeded = errors.New("quota exceeded")

func defaultChatFunc(messages []ai.ChatMessage) (*ai.ChatResponse, error) {
	if config.GlobalConfig.AI.DashScopeAPIKey == "" {
		return nil, errors.New("dashscope api key missing")
	}
	client := ai.NewClient()
	return client.Chat(messages)
}

func defaultSpeechFunc(data []byte) (*ai.SpeechResult, error) {
	return ai.RecognizeSpeech(data)
}

func buildChatMessages(req AIChatRequest, baby *model.Baby) []ai.ChatMessage {
	messages := []ai.ChatMessage{{Role: "system", Content: ai.BuildSystemPrompt()}}
	if baby != nil {
		messages = append(messages, ai.ChatMessage{Role: "system", Content: buildBabyContext(baby)})
	}
	if len(req.History) > 0 {
		start := 0
		if len(req.History) > 4 {
			start = len(req.History) - 4
		}
		for _, item := range req.History[start:] {
			if item.Role == "" || item.Content == "" {
				continue
			}
			messages = append(messages, ai.ChatMessage{Role: item.Role, Content: item.Content})
		}
	}
	messages = append(messages, ai.ChatMessage{Role: "user", Content: req.Question})
	return messages
}

func buildBabyContext(baby *model.Baby) string {
	gender := "未知"
	if baby.Gender == 1 {
		gender = "男"
	} else if baby.Gender == 2 {
		gender = "女"
	}
	months := baby.AgeInMonths()
	context := "宝宝信息：姓名=" + baby.Name + "，性别=" + gender + "，月龄=" + strconv.Itoa(months) + "，生日=" + baby.Birthday.Format("2006-01-02")
	if baby.BirthWeight != nil {
		context += "，出生体重=" + strconv.FormatFloat(*baby.BirthWeight, 'f', 2, 64) + "kg"
	}
	if baby.BirthHeight != nil {
		context += "，出生身高=" + strconv.FormatFloat(*baby.BirthHeight, 'f', 1, 64) + "cm"
	}
	if baby.IsPremature == 1 {
		context += "，早产"
	}
	return context
}

func quotaKey(userID, quotaType string) string {
	date := time.Now().In(time.Local).Format("20060102")
	return "ai:quota:" + userID + ":" + quotaType + ":" + date
}

func getQuotaUsage(userID, quotaType string) (int, error) {
	val, err := gredis.Get(quotaKey(userID, quotaType))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func consumeQuota(userID, quotaType string, limit int) (int, int, error) {
	used, err := getQuotaUsage(userID, quotaType)
	if err != nil {
		return 0, 0, err
	}
	if used >= limit {
		return used, 0, errQuotaExceeded
	}
	key := quotaKey(userID, quotaType)
	value, err := gredis.Incr(key)
	if err != nil {
		return used, 0, err
	}
	_ = gredis.Expire(key, int(quotaTTL().Seconds()))
	newUsed := int(value)
	return newUsed, calcRemaining(limit, newUsed), nil
}

func rollbackQuota(userID, quotaType string) error {
	_, err := gredis.Decr(quotaKey(userID, quotaType))
	return err
}

func calcRemaining(limit, used int) int {
	remaining := limit - used
	if remaining < 0 {
		return 0
	}
	return remaining
}

func quotaTTL() time.Duration {
	now := time.Now().In(time.Local)
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return next.Sub(now) + time.Minute
}
