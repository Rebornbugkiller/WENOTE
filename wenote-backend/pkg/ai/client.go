package ai

import "context"

// SummaryResult 摘要生成结果
type SummaryResult struct {
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
}

// Client AI 客户端接口
type Client interface {
	GenerateSummaryAndTags(ctx context.Context, content string, summaryLen int) (*SummaryResult, error)
}
