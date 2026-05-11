package ai

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// QuotaStore 额度与熔断存储接口
type QuotaStore interface {
	CanUseProvider(ctx context.Context, provider ProviderName, req ChatRequest) bool
	RecordSuccess(ctx context.Context, provider ProviderName, resp *ChatResponse)
	RecordFailure(ctx context.Context, provider ProviderName, err error)
}

// Router AI 路由器：按优先级链依次尝试 Provider，失败自动 fallback
type Router struct {
	providers []Provider
	quota     QuotaStore
}

// NewRouter 创建路由器，provider 顺序决定优先级
func NewRouter(quota QuotaStore, providers ...Provider) *Router {
	return &Router{providers: providers, quota: quota}
}

// Chat 按优先级链调用 AI，失败自动 fallback
func (r *Router) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	var errs []error

	for _, p := range r.providers {
		if !p.Enabled() {
			continue
		}
		if r.quota != nil && !r.quota.CanUseProvider(ctx, p.Name(), req) {
			continue
		}

		start := time.Now()
		resp, err := p.Chat(ctx, req)
		latency := time.Since(start)

		if err == nil && resp != nil && resp.Content != "" {
			if r.quota != nil {
				r.quota.RecordSuccess(ctx, p.Name(), resp)
			}
			_ = latency
			return resp, nil
		}

		if err == nil {
			err = fmt.Errorf("%s returned empty response", p.Name())
		}
		if r.quota != nil {
			r.quota.RecordFailure(ctx, p.Name(), err)
		}
		errs = append(errs, fmt.Errorf("%s: %w", p.Name(), err))
	}

	if len(errs) == 0 {
		return nil, errors.New("no provider available")
	}
	return nil, fmt.Errorf("all providers failed: %w", errors.Join(errs...))
}
