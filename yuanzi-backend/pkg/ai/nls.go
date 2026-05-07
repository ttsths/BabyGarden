package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"yuanzi-backend/config"

	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
)

type SpeechResult struct {
	Text       string
	Confidence float64
}

type nlsSpeechPayload struct {
	Result     string  `json:"result"`
	Confidence float64 `json:"confidence"`
}

// RecognizeSpeech 语音识别（阿里云 NLS SDK）
func RecognizeSpeech(audio []byte) (*SpeechResult, error) {
	if len(audio) == 0 {
		return nil, errors.New("audio empty")
	}
	cfg := config.GlobalConfig.AI
	if cfg.NLSAppKey == "" || cfg.NLSAccessKeyID == "" || cfg.NLSAccessKeySecret == "" {
		return nil, errors.New("nls config missing")
	}

	url := nls.DEFAULT_URL
	if cfg.NLSEndpoint != "" {
		url = cfg.NLSEndpoint
	}
	config, err := nls.NewConnectionConfigWithAKInfoDefault(url, cfg.NLSAppKey, cfg.NLSAccessKeyID, cfg.NLSAccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("nls config error: %w", err)
	}
	
	logger := nls.DefaultNlsLog() // 默认日志
	done := make(chan struct{})
	result := &SpeechResult{}
	var resultErr error

	onCompleted := func(text string, _ interface{}) {
		parseSpeechResult(text, result)
		close(done)
	}
	onFailed := func(text string, _ interface{}) {
		resultErr = fmt.Errorf("nls failed: %s", text)
		close(done)
	}

	sr, err := nls.NewSpeechRecognition(config, logger, onFailed, nil, nil, onCompleted, nil, nil)
	if err != nil {
		return nil, err
	}
	defer sr.Shutdown()

	param := nls.DefaultSpeechRecognitionParam()
	param.Format = "wav"
	param.SampleRate = 16000
	if _, err := sr.Start(param, nil); err != nil {
		return nil, err
	}
	if err := sr.SendAudioData(audio); err != nil {
		return nil, err
	}
	_, _ = sr.Stop()

	select {
	case <-done:
		if resultErr != nil {
			return nil, resultErr
		}
		if result.Text == "" {
			result.Text = ""
		}
		return result, nil
	case <-time.After(30 * time.Second):
		return nil, errors.New("nls timeout")
	}
}

func parseSpeechResult(raw string, result *SpeechResult) {
	if result == nil || raw == "" {
		return
	}
	var payload nlsSpeechPayload
	if err := json.Unmarshal([]byte(raw), &payload); err == nil {
		if payload.Result != "" {
			result.Text = payload.Result
			result.Confidence = payload.Confidence
			return
		}
	}
	// fallback
	result.Text = raw
}
