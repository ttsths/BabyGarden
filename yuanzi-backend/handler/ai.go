package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

	// 确定首选 Provider 用于动态系统提示词
	firstProvider := firstEnabledProvider()
	messages := buildChatMessages(req, baby, firstProvider)
	requestID := model.NewID() // 生成 request_id 用于跨 Provider 追踪
	resp, err := aiChatFunc(c.Request.Context(), messages)
	if err != nil {
		_ = rollbackQuota(userID, quotaTypeChat)
		logAIUsageAsync(model.AIUsageLog{
			RequestID:    requestID,
			UserID:       userID,
			FamilyID:     familyIDFromBaby(baby),
			RequestType:  "chat",
			Status:       "error",
			ErrorMessage: err.Error(),
		})
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "AI服务异常"})
		return
	}

	answer := strings.TrimSpace(resp.Content)
	if answer == "" {
		answer = "AI暂时无法回答，请稍后再试"
	}
	totalTokens := resp.Usage.TotalTokens
	if totalTokens == 0 {
		totalTokens = resp.Usage.InputTokens + resp.Usage.OutputTokens
	}

	record := model.AIChatRecord{
		UserID:     userID,
		BabyID:     req.BabyID,
		Question:   req.Question,
		Answer:     answer,
		TokensUsed: totalTokens,
		Model:      resp.Model,
		CreatedAt:  time.Now(),
	}
	if err := mysql.DB.Create(&record).Error; err != nil {
		_ = rollbackQuota(userID, quotaTypeChat)
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "保存问答记录失败"})
		return
	}
	logAIUsageAsync(model.AIUsageLog{
		RequestID:    requestID,
		UserID:       userID,
		FamilyID:     familyIDFromBaby(baby),
		Provider:     string(resp.Provider),
		Model:        resp.Model,
		InputTokens:  resp.Usage.InputTokens,
		OutputTokens: resp.Usage.OutputTokens,
		CachedTokens: resp.Usage.CachedTokens,
		TotalTokens:  totalTokens,
		RequestType:  "chat",
		Status:       "success",
	})

	c.JSON(http.StatusOK, model.Response{
		Code: model.SUCCESS,
		Msg:  "success",
		Data: AIChatResponse{
			Answer:         answer,
			TokensUsed:     totalTokens,
			InputTokens:    resp.Usage.InputTokens,
			OutputTokens:   resp.Usage.OutputTokens,
			CachedTokens:   resp.Usage.CachedTokens,
			TotalTokens:    totalTokens,
			Provider:       string(resp.Provider),
			RemainingQuota: remaining,
		},
	})
}

// AIChatStream AI 流式问答。
func AIChatStream(c *gin.Context) {
	var req AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Code: model.ERROR_INVALID, Msg: "请求参数错误"})
		return
	}

	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	_, remaining, err := consumeQuota(userID, quotaTypeChat, chatDailyLimit)
	if err != nil {
		if errors.Is(err, errQuotaExceeded) {
			c.JSON(http.StatusConflict, model.Response{Code: model.ERROR_QUOTA_EXCEEDED, Msg: "AI配额已用完"})
			return
		}
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "配额检查失败"})
		return
	}

	var baby *model.Baby
	if req.BabyID != nil && *req.BabyID != "" {
		loaded, _, err := loadBabyForRecord(c, *req.BabyID)
		if err != nil {
			_ = rollbackQuota(userID, quotaTypeChat)
			return
		}
		baby = loaded
	}

	firstProvider := firstEnabledProvider()
	messages := buildChatMessages(req, baby, firstProvider)
	requestID := model.NewID()
	resp, err := aiChatFunc(c.Request.Context(), messages)
	if err != nil {
		_ = rollbackQuota(userID, quotaTypeChat)
		logAIUsageAsync(model.AIUsageLog{
			RequestID:    requestID,
			UserID:       userID,
			FamilyID:     familyIDFromBaby(baby),
			RequestType:  "chat_stream",
			Status:       "error",
			ErrorMessage: err.Error(),
		})
		streamAIEvent(c, "error", gin.H{"message": "AI服务异常"})
		return
	}

	answer := strings.TrimSpace(resp.Content)
	if answer == "" {
		answer = "AI暂时无法回答，请稍后再试"
	}
	totalTokens := resp.Usage.TotalTokens
	if totalTokens == 0 {
		totalTokens = resp.Usage.InputTokens + resp.Usage.OutputTokens
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	for _, chunk := range splitStreamChunks(answer) {
		streamAIEvent(c, "delta", gin.H{"delta": chunk})
	}

	record := model.AIChatRecord{
		UserID:     userID,
		BabyID:     req.BabyID,
		Question:   req.Question,
		Answer:     answer,
		TokensUsed: totalTokens,
		Model:      resp.Model,
		CreatedAt:  time.Now(),
	}
	if err := mysql.DB.Create(&record).Error; err != nil {
		streamAIEvent(c, "error", gin.H{"message": "保存问答记录失败"})
		return
	}
	logAIUsageAsync(model.AIUsageLog{
		RequestID:    requestID,
		UserID:       userID,
		FamilyID:     familyIDFromBaby(baby),
		Provider:     string(resp.Provider),
		Model:        resp.Model,
		InputTokens:  resp.Usage.InputTokens,
		OutputTokens: resp.Usage.OutputTokens,
		CachedTokens: resp.Usage.CachedTokens,
		TotalTokens:  totalTokens,
		RequestType:  "chat_stream",
		Status:       "success",
	})
	streamAIEvent(c, "done", gin.H{
		"id":              record.ID,
		"answer":          answer,
		"tokens_used":     totalTokens,
		"provider":        string(resp.Provider),
		"remaining_quota": remaining,
	})
}

// ListAIChats 获取当前用户 AI 会话历史。
func ListAIChats(c *gin.Context) {
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	babyID := c.Query("baby_id")
	if babyID != "" {
		if _, _, err := loadBabyForRecord(c, babyID); err != nil {
			return
		}
	}
	page := parsePage(c.DefaultQuery("page", "1"))
	pageSize := parsePageSize(c.DefaultQuery("page_size", "20"))

	query := mysql.DB.Model(&model.AIChatRecord{}).Where("user_id = ?", userID)
	if babyID != "" {
		query = query.Where("baby_id = ?", babyID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}

	var records []model.AIChatRecord
	if err := query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: model.ERROR, Msg: "查询失败"})
		return
	}
	list := make([]AIChatHistoryResponse, 0, len(records))
	for _, item := range records {
		list = append(list, aiChatHistoryResponse(item))
	}

	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: model.ListResponse{
		List: list,
		Pagination: model.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: calcTotalPages(total, pageSize),
		},
	}})
}

func splitStreamChunks(answer string) []string {
	runes := []rune(answer)
	if len(runes) <= 12 {
		return []string{answer}
	}
	chunks := make([]string, 0, len(runes)/12+1)
	for start := 0; start < len(runes); start += 12 {
		end := start + 12
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}

func streamAIEvent(c *gin.Context, event string, payload interface{}) {
	data, _ := json.Marshal(payload)
	if _, err := c.Writer.WriteString("event: " + event + "\n"); err != nil {
		return
	}
	if _, err := c.Writer.WriteString("data: " + string(data) + "\n\n"); err != nil {
		return
	}
	c.Writer.Flush()
}

// GetAIChat 获取当前用户 AI 会话详情。
func GetAIChat(c *gin.Context) {
	userID := middleware.GetUserIDOrZero(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.Response{Code: model.ERROR_NOT_AUTH, Msg: "未登录"})
		return
	}

	var record model.AIChatRecord
	if err := mysql.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, model.Response{Code: model.ERROR_NOT_FUND, Msg: "会话不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: model.SUCCESS, Msg: "获取成功", Data: aiChatHistoryResponse(record)})
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
// @Description 获取用户当日 AI 使用配额情况，含 Provider 链与 Cloudflare 估算额度
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

	cfg := config.GlobalConfig.AI
	speechUsed, _ := getQuotaUsage(userID, quotaTypeSpeech)
	chatUsed, _ := getQuotaUsage(userID, quotaTypeChat)
	tokenStats := getUserTokenStats(userID)

	// Cloudflare Neurons 估算
	cfNeuronsUsed, _ := getCloudflareNeurons()
	cfBudget := cfg.Cloudflare.DailyNeuronBudget
	cfRemaining := cfBudget - cfNeuronsUsed
	if cfRemaining < 0 {
		cfRemaining = 0
	}

	resetAt := nextUTCMidnight()

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
				Tokens:    tokenStats,
			},
			ProviderChain: cfg.ProviderChain,
			Cloudflare: &CloudflareQuotaDetail{
				Model:             cfg.Cloudflare.Model,
				DailyNeuronBudget: cfBudget,
				HardNeuronBudget:  cfg.Cloudflare.HardNeuronBudget,
				NeuronsUsedToday:  cfNeuronsUsed,
				NeuronsRemaining:  cfRemaining,
				ResetAtUTC:        resetAt.Format(time.RFC3339),
				SourceOfTruth:     "Cloudflare Workers AI Dashboard",
			},
		},
	})
}

// CloudflareQuotaDetail Cloudflare Workers AI 额度明细
type CloudflareQuotaDetail struct {
	Model             string `json:"model"`
	DailyNeuronBudget int    `json:"daily_neuron_budget"`
	HardNeuronBudget  int    `json:"hard_neuron_budget"`
	NeuronsUsedToday  int    `json:"estimated_neurons_used_today"`
	NeuronsRemaining  int    `json:"estimated_neurons_remaining"`
	ResetAtUTC        string `json:"reset_at_utc"`
	SourceOfTruth     string `json:"source_of_truth"`
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
	InputTokens    int    `json:"input_tokens"`
	OutputTokens   int    `json:"output_tokens"`
	CachedTokens   int    `json:"cached_tokens"`
	TotalTokens    int    `json:"total_tokens"`
	Provider       string `json:"provider"`
	RemainingQuota int    `json:"remaining_quota"`
}

type SpeechRecognizeResponse struct {
	Text           string  `json:"text"`
	Confidence     float64 `json:"confidence"`
	RemainingQuota int     `json:"remaining_quota"`
}

type AIQuotaResponse struct {
	Speech        QuotaDetail            `json:"speech"`
	AIChat        QuotaDetail            `json:"ai_chat"`
	ProviderChain []string               `json:"provider_chain,omitempty"`
	Cloudflare    *CloudflareQuotaDetail `json:"cloudflare,omitempty"`
}

type QuotaDetail struct {
	Used      int              `json:"used"`
	Limit     int              `json:"limit"`
	Remaining int              `json:"remaining"`
	Tokens    *TokenQuotaStats `json:"tokens,omitempty"`
}

type TokenQuotaStats struct {
	TodayInputTokens  int `json:"today_input_tokens"`
	TodayOutputTokens int `json:"today_output_tokens"`
	TodayCachedTokens int `json:"today_cached_tokens"`
	TodayTotalTokens  int `json:"today_total_tokens"`
	MonthTotalTokens  int `json:"month_total_tokens"`
}

type AIChatHistoryResponse struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	BabyID     *string `json:"baby_id,omitempty"`
	Question   string  `json:"question"`
	Answer     string  `json:"answer"`
	TokensUsed int     `json:"tokens_used"`
	Model      string  `json:"model"`
	CreatedAt  string  `json:"created_at"`
}

var (
	aiChatFunc          = defaultChatFunc
	speechRecognizeFunc = defaultSpeechFunc
)

func aiChatHistoryResponse(record model.AIChatRecord) AIChatHistoryResponse {
	return AIChatHistoryResponse{
		ID:         record.ID,
		UserID:     record.UserID,
		BabyID:     record.BabyID,
		Question:   record.Question,
		Answer:     record.Answer,
		TokensUsed: record.TokensUsed,
		Model:      record.Model,
		CreatedAt:  record.CreatedAt.Format(time.RFC3339),
	}
}

const (
	quotaTypeChat    = "chat"
	quotaTypeSpeech  = "speech"
	chatDailyLimit   = 50
	speechDailyLimit = 20
)

var (
	errQuotaExceeded = errors.New("quota exceeded")
	aiRouterOnce     sync.Once
	aiRouter         *ai.Router
	aiQuotaStore     *redisQuotaStore
)

// initAIRouter 延迟初始化 AI Router（首次调用时创建）
func initAIRouter() {
	aiRouterOnce.Do(func() {
		cfg := config.GlobalConfig.AI
		aiQuotaStore = newRedisQuotaStore()

		var providers []ai.Provider

		// 1. GrokAI（首选）
		if cfg.GrokAI.Enabled || cfg.GrokAI.BaseURL != "" {
			providers = append(providers, ai.NewOpenAICompatProvider(
				ai.ProviderGrokAI,
				cfg.GrokAI.Enabled,
				cfg.GrokAI.BaseURL,
				cfg.GrokAI.APIKey,
				cfg.GrokAI.Model,
				time.Duration(cfg.GrokAI.TimeoutSeconds)*time.Second,
			))
		}

		// 2. Cloudflare Workers AI
		if cfg.Cloudflare.Enabled || cfg.Cloudflare.APIToken != "" {
			providers = append(providers, ai.NewCloudflareWorkersAIProvider(
				cfg.Cloudflare.Enabled,
				cfg.Cloudflare.AccountID,
				cfg.Cloudflare.APIToken,
				cfg.Cloudflare.GatewayID,
				cfg.Cloudflare.UseGateway,
				cfg.Cloudflare.Model,
				time.Duration(cfg.Cloudflare.TimeoutSeconds)*time.Second,
			))
		}

		// 3. DashScope
		dsKey := cfg.DashScope.APIKey
		if dsKey == "" {
			dsKey = cfg.DashScopeAPIKey // fallback 到旧字段
		}
		providers = append(providers, ai.NewDashScopeProvider(
			cfg.DashScope.Enabled,
			dsKey,
			cfg.DashScope.Model,
			time.Duration(cfg.DashScope.TimeoutSeconds)*time.Second,
		))

		// 4. CLIProxyAPI（最后 fallback）
		if cfg.CLIProxyAPI.Enabled || cfg.CLIProxyAPI.BaseURL != "" {
			providers = append(providers, ai.NewOpenAICompatProvider(
				ai.ProviderCLIProxyAPI,
				cfg.CLIProxyAPI.Enabled,
				cfg.CLIProxyAPI.BaseURL,
				cfg.CLIProxyAPI.APIKey,
				cfg.CLIProxyAPI.Model,
				time.Duration(cfg.CLIProxyAPI.TimeoutSeconds)*time.Second,
			))
		}

		aiRouter = ai.NewRouter(aiQuotaStore, providers...)
	})
}

func defaultChatFunc(ctx context.Context, messages []ai.ChatMessage) (*ai.ChatResponse, error) {
	initAIRouter()

	req := ai.ChatRequest{
		Messages:    messages,
		MaxTokens:   config.GlobalConfig.AI.Safety.MaxOutputTokens,
		Temperature: 0.7,
	}

	resp, err := aiRouter.Chat(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func firstEnabledProvider() ai.ProviderName {
	initAIRouter()
	for _, p := range aiRouter.Providers() {
		if p.Enabled() {
			return p.Name()
		}
	}
	return ai.ProviderDashScope // fallback
}

func defaultSpeechFunc(data []byte) (*ai.SpeechResult, error) {
	return ai.RecognizeSpeech(data)
}

func buildChatMessages(req AIChatRequest, baby *model.Baby, provider ai.ProviderName) []ai.ChatMessage {
	messages := []ai.ChatMessage{{Role: "system", Content: ai.BuildSystemPrompt(provider)}}
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

func familyIDFromBaby(baby *model.Baby) string {
	if baby == nil {
		return ""
	}
	return baby.FamilyID
}

func logAIUsageAsync(log model.AIUsageLog) {
	if mysql.DB == nil {
		return
	}
	go func(item model.AIUsageLog) {
		_ = mysql.DB.Create(&item).Error
	}(log)
}

func getUserTokenStats(userID string) *TokenQuotaStats {
	if mysql.DB == nil {
		return &TokenQuotaStats{}
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	stats := &TokenQuotaStats{}
	type tokenSum struct {
		InputTokens  int
		OutputTokens int
		CachedTokens int
		TotalTokens  int
	}
	var today tokenSum
	if err := mysql.DB.Model(&model.AIUsageLog{}).
		Select("COALESCE(SUM(input_tokens), 0) AS input_tokens, COALESCE(SUM(output_tokens), 0) AS output_tokens, COALESCE(SUM(cached_tokens), 0) AS cached_tokens, COALESCE(SUM(total_tokens), 0) AS total_tokens").
		Where("user_id = ? AND status = ? AND created_at >= ?", userID, "success", todayStart).
		Scan(&today).Error; err == nil {
		stats.TodayInputTokens = today.InputTokens
		stats.TodayOutputTokens = today.OutputTokens
		stats.TodayCachedTokens = today.CachedTokens
		stats.TodayTotalTokens = today.TotalTokens
	}

	var month tokenSum
	if err := mysql.DB.Model(&model.AIUsageLog{}).
		Select("COALESCE(SUM(total_tokens), 0) AS total_tokens").
		Where("user_id = ? AND status = ? AND created_at >= ?", userID, "success", monthStart).
		Scan(&month).Error; err == nil {
		stats.MonthTotalTokens = month.TotalTokens
	}
	return stats
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

// ============================================================================
// Redis QuotaStore — 实现 ai.QuotaStore 接口
// ============================================================================

// redisQuotaStore 基于 Redis 的额度与熔断存储
type redisQuotaStore struct{}

func newRedisQuotaStore() *redisQuotaStore {
	return &redisQuotaStore{}
}

func (s *redisQuotaStore) CanUseProvider(ctx context.Context, provider ai.ProviderName, req ai.ChatRequest) bool {
	// 检查熔断
	circuitKey := "ai:provider:" + string(provider) + ":circuit_until"
	val, err := gredis.Get(circuitKey)
	if err == nil && val != "" {
		until, parseErr := time.Parse(time.RFC3339, val)
		if parseErr == nil && time.Now().Before(until) {
			return false
		}
	}

	// Cloudflare 特有：检查 Neurons 预算
	if provider == ai.ProviderCloudflareWorkersAI {
		cfg := config.GlobalConfig.AI.Cloudflare
		neuronsUsed, _ := getCloudflareNeurons()
		if neuronsUsed >= cfg.HardNeuronBudget {
			return false
		}
		if neuronsUsed >= cfg.DailyNeuronBudget {
			return false
		}
	}

	return true
}

func (s *redisQuotaStore) RecordSuccess(ctx context.Context, provider ai.ProviderName, resp *ai.ChatResponse) {
	// 重置熔断计数
	failKey := "ai:provider:" + string(provider) + ":failures"
	_ = gredis.Del(failKey)
	_ = gredis.Del("ai:provider:" + string(provider) + ":circuit_until")

	// Cloudflare：记录 Neurons 消耗
	if provider == ai.ProviderCloudflareWorkersAI && resp.Usage.NeuronsEst > 0 {
		addCloudflareNeurons(resp.Usage.NeuronsEst)
	}
}

func (s *redisQuotaStore) RecordFailure(ctx context.Context, provider ai.ProviderName, err error) {
	failKey := "ai:provider:" + string(provider) + ":failures"
	val, incrErr := gredis.Incr(failKey)
	if incrErr != nil {
		return
	}
	_ = gredis.Expire(failKey, 300) // 5 分钟窗口

	// 连续失败 >= 3 次 → 熔断 60 秒
	if val >= 3 {
		circuitKey := "ai:provider:" + string(provider) + ":circuit_until"
		until := time.Now().Add(60 * time.Second).UTC().Format(time.RFC3339)
		_ = gredis.Set(circuitKey, until, 60)
	}
}

// ============================================================================
// Cloudflare Neurons 追踪
// ============================================================================

func cloudflareNeuronsKey() string {
	return "babygarden:ai:cf:neurons:" + time.Now().UTC().Format("2006-01-02")
}

func getCloudflareNeurons() (int, error) {
	val, err := gredis.Get(cloudflareNeuronsKey())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

func addCloudflareNeurons(neurons int) {
	key := cloudflareNeuronsKey()
	_, _ = gredis.IncrBy(key, int64(neurons))
	_ = gredis.Expire(key, 48*3600) // 48 小时过期
}

func nextUTCMidnight() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
}
