package handler

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"yuanzi-backend/model"
	"yuanzi-backend/mysql"
	"yuanzi-backend/pkg/ai"
	"yuanzi-backend/pkg/gredis"
)

func TestAIChatSuccess(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("aichat"), "AI用户")
	family := createTestFamily(t, admin.ID, "AI家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "AI宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)
	defer cleanupAIChatRecords(t, admin.ID)
	clearAIQuota(t, admin.ID)

	aiChatFunc = func(_ context.Context, messages []ai.ChatMessage) (*ai.ChatResponse, error) {
		return &ai.ChatResponse{
			Content:  "测试回答",
			Provider: ai.ProviderDashScope,
			Model:    "qwen-turbo",
			Usage: ai.Usage{
				InputTokens:  5,
				OutputTokens: 7,
				CachedTokens: 2,
				TotalTokens:  12,
			},
		}, nil
	}
	defer resetAIHandlers()

	body := mustMarshal(t, AIChatRequest{Question: "宝宝喝多少奶?", BabyID: &baby.ID})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/ai/chat", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("userId", admin.ID)

	AIChat(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("AI问答失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Code int            `json:"code"`
		Data AIChatResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if response.Data.Answer != "测试回答" || response.Data.TokensUsed != 12 || response.Data.InputTokens != 5 || response.Data.OutputTokens != 7 || response.Data.CachedTokens != 2 || response.Data.Provider != string(ai.ProviderDashScope) {
		t.Fatalf("AI问答响应错误: %+v", response.Data)
	}

	var record model.AIChatRecord
	if err := mysql.DB.Where("user_id = ? AND question = ?", admin.ID, "宝宝喝多少奶?").First(&record).Error; err != nil {
		t.Fatalf("问答记录未保存: %v", err)
	}
}

func TestAIChatStreamSuccess(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("aistream"), "AI流用户")
	family := createTestFamily(t, admin.ID, "AI流家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "AI流宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)
	defer cleanupAIChatRecords(t, admin.ID)
	clearAIQuota(t, admin.ID)

	aiChatFunc = func(_ context.Context, messages []ai.ChatMessage) (*ai.ChatResponse, error) {
		return &ai.ChatResponse{
			Content:  "睡眠平稳，奶量正常，排泄规律。",
			Provider: ai.ProviderDashScope,
			Model:    "qwen-turbo",
			Usage:    ai.Usage{InputTokens: 6, OutputTokens: 8, TotalTokens: 14},
		}, nil
	}
	defer resetAIHandlers()

	body := mustMarshal(t, AIChatRequest{Question: "分析近一周趋势", BabyID: &baby.ID})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/ai/chat/stream", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("userId", admin.ID)

	AIChatStream(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("AI流式问答失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); !strings.Contains(got, "text/event-stream") {
		t.Fatalf("Content-Type 错误: %s", got)
	}
	bodyText := recorder.Body.String()
	if !strings.Contains(bodyText, "event: delta") || !strings.Contains(bodyText, "event: done") {
		t.Fatalf("流式响应缺少事件: %s", bodyText)
	}

	var record model.AIChatRecord
	if err := mysql.DB.Where("user_id = ? AND question = ?", admin.ID, "分析近一周趋势").First(&record).Error; err != nil {
		t.Fatalf("流式问答记录未保存: %v", err)
	}
}

func TestListAIChats(t *testing.T) {
	setupFamilyTestDB(t)

	admin := createTestUser(t, uniquePhone("aihist"), "AI历史用户")
	family := createTestFamily(t, admin.ID, "AI历史家庭")
	member := createTestFamilyMember(t, family.ID, admin.ID, model.FamilyRoleAdmin)
	baby := createTestBaby(t, family.ID, "AI历史宝宝")
	defer cleanupFamilies(t, family.ID)
	defer cleanupUsers(t, admin.ID)
	defer cleanupMembers(t, member.ID)
	defer cleanupAIChatRecords(t, admin.ID)

	record := model.AIChatRecord{UserID: admin.ID, BabyID: &baby.ID, Question: "今天喝奶如何?", Answer: "正常记录。", TokensUsed: 8, Model: "test", CreatedAt: time.Now()}
	if err := mysql.DB.Create(&record).Error; err != nil {
		t.Fatalf("创建AI历史失败: %v", err)
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/ai/chats?baby_id="+baby.ID, nil)
	ctx.Set("userId", admin.ID)

	ListAIChats(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取AI历史失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var response struct {
		Code int `json:"code"`
		Data struct {
			List []AIChatHistoryResponse `json:"list"`
		} `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if len(response.Data.List) != 1 || response.Data.List[0].Question != record.Question {
		t.Fatalf("AI历史返回错误: %+v", response.Data.List)
	}
}

func TestSpeechRecognizeSuccess(t *testing.T) {
	setupFamilyTestDB(t)

	user := createTestUser(t, uniquePhone("speech"), "语音用户")
	defer cleanupUsers(t, user.ID)
	clearAIQuota(t, user.ID)

	speechRecognizeFunc = func(data []byte) (*ai.SpeechResult, error) {
		return &ai.SpeechResult{Text: "识别结果", Confidence: 0.9}, nil
	}
	defer resetAIHandlers()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("audio", "test.wav")
	_, _ = part.Write([]byte("audio"))
	writer.Close()

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/ai/speech/recognize", body)
	ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
	ctx.Set("userId", user.ID)

	SpeechRecognize(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("语音识别失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestGetAIQuota(t *testing.T) {
	setupFamilyTestDB(t)

	user := createTestUser(t, uniquePhone("quota"), "配额用户")
	defer cleanupUsers(t, user.ID)
	clearAIQuota(t, user.ID)

	_ = gredis.SetEx(quotaKey(user.ID, quotaTypeChat), "3", 3600)
	_ = gredis.SetEx(quotaKey(user.ID, quotaTypeSpeech), "2", 3600)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/ai/quota", nil)
	ctx.Set("userId", user.ID)

	GetAIQuota(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取配额失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestGetAIHistory(t *testing.T) {
	setupFamilyTestDB(t)

	user := createTestUser(t, uniquePhone("aihist"), "AI历史用户")
	defer cleanupUsers(t, user.ID)
	defer cleanupAIChatRecords(t, user.ID)

	record := model.AIChatRecord{
		UserID:     user.ID,
		Question:   "怎么安排午睡？",
		Answer:     "先观察困倦信号。",
		TokensUsed: 8,
		Model:      "test-model",
		CreatedAt:  time.Now(),
	}
	if err := mysql.DB.Create(&record).Error; err != nil {
		t.Fatalf("创建AI历史失败: %v", err)
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/ai/history?page=1&page_size=10", nil)
	ctx.Set("userId", user.ID)

	ListAIChats(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("获取AI历史失败: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var response struct {
		Data model.ListResponse `json:"data"`
	}
	decodeResponse(t, recorder.Body.Bytes(), &response)
	if response.Data.Pagination.Total < 1 {
		t.Fatalf("AI历史为空: %+v", response.Data)
	}
}

func cleanupAIChatRecords(t *testing.T, userIDs ...string) {
	t.Helper()
	if len(userIDs) == 0 {
		return
	}
	_ = mysql.DB.Where("user_id IN ?", userIDs).Delete(&model.AIChatRecord{}).Error
}

func clearAIQuota(t *testing.T, userID string) {
	t.Helper()
	_ = gredis.Del(quotaKey(userID, quotaTypeChat))
	_ = gredis.Del(quotaKey(userID, quotaTypeSpeech))
}

func resetAIHandlers() {
	aiChatFunc = defaultChatFunc
	speechRecognizeFunc = defaultSpeechFunc
}
