package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/sony/gobreaker"
)

// ZhipuConfig 智谱配置
type ZhipuConfig struct {
	APIKey     string
	Model      string
	BaseURL    string
	Timeout    int
	MaxRetries int
	RetryDelay int
}

// ZhipuClient 智谱客户端
type ZhipuClient struct {
	config     ZhipuConfig
	httpClient *http.Client
	breaker    *gobreaker.CircuitBreaker
}

// NewZhipuClient 创建智谱客户端
func NewZhipuClient(config ZhipuConfig) *ZhipuClient {
	// 设置默认值
	if config.Timeout <= 0 {
		config.Timeout = 30
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 2
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = 2
	}

	timeout := time.Duration(config.Timeout) * time.Second

	// 配置熔断器
	breakerSettings := gobreaker.Settings{
		Name:        "AI-Service",
		MaxRequests: 3,                // 半开状态最大请求数
		Interval:    60 * time.Second, // 统计周期
		Timeout:     timeout,          // 熔断后恢复时间
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// 失败率超过 50% 且请求数 >= 5 时触发熔断
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.5
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			slog.Warn("Circuit breaker state changed",
				"name", name,
				"from", from.String(),
				"to", to.String(),
			)
		},
	}

	return &ZhipuClient{
		config: config,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		breaker: gobreaker.NewCircuitBreaker(breakerSettings),
	}
}

// 智谱 API 请求结构
type zhipuRequest struct {
	Model    string         `json:"model"`
	Messages []zhipuMessage `json:"messages"`
}

type zhipuMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 智谱 API 响应结构
type zhipuResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// GenerateSummaryAndTags 生成摘要和标签
func (c *ZhipuClient) GenerateSummaryAndTags(ctx context.Context, content string, summaryLen int) (*SummaryResult, error) {
	prompt := buildPrompt(content, summaryLen)

	// 构建请求
	reqBody := zhipuRequest{
		Model: c.config.Model,
		Messages: []zhipuMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	// 发送请求(带重试)
	maxAttempts := c.config.MaxRetries + 1
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, err := c.doRequest(ctx, jsonData)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// 如果是熔断器打开或上下文取消，不重试
		if err == gobreaker.ErrOpenState || err == gobreaker.ErrTooManyRequests {
			slog.Warn("Circuit breaker active, stop retrying", "error", err)
			return nil, lastErr
		}
		if ctx.Err() != nil {
			slog.Warn("Context cancelled, stop retrying", "error", ctx.Err())
			return nil, lastErr
		}

		slog.Warn("AI request failed",
			"attempt", attempt,
			"error", err,
		)

		// 最后一次尝试不需要等待
		if attempt < maxAttempts {
			// 指数退避
			delay := time.Duration(c.config.RetryDelay*attempt) * time.Second
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("all retries failed: %w", lastErr)
}

// doRequest 执行单次请求（带熔断器保护）
func (c *ZhipuClient) doRequest(ctx context.Context, jsonData []byte) (*SummaryResult, error) {
	// 通过熔断器执行请求
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.makeHTTPRequest(ctx, jsonData)
	})

	if err != nil {
		// 判断是否为熔断器打开错误
		if err == gobreaker.ErrOpenState {
			slog.Warn("Circuit breaker is open, request rejected")
			return nil, fmt.Errorf("AI 服务暂时不可用，请稍后重试")
		}
		return nil, err
	}

	return result.(*SummaryResult), nil
}

// makeHTTPRequest 执行实际的 HTTP 请求
func (c *ZhipuClient) makeHTTPRequest(ctx context.Context, jsonData []byte) (*SummaryResult, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 解析响应
	var zhipuResp zhipuResponse
	if err := json.Unmarshal(body, &zhipuResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	if zhipuResp.Error.Message != "" {
		return nil, fmt.Errorf("API error: %s", zhipuResp.Error.Message)
	}

	if len(zhipuResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	// 解析 JSON 结果
	var result SummaryResult
	contentStr := zhipuResp.Choices[0].Message.Content

	// 清理 markdown 代码块标记（智谱 AI 可能返回 ```json ... ``` 格式）
	contentStr = cleanMarkdownCodeBlock(contentStr)

	if err := json.Unmarshal([]byte(contentStr), &result); err != nil {
		return nil, fmt.Errorf("parse AI result failed: %w, content=%s", err, contentStr)
	}

	return &result, nil
}

// cleanMarkdownCodeBlock 清理 markdown 代码块标记
func cleanMarkdownCodeBlock(content string) string {
	// 移除开头的 ```json 或 ```
	if len(content) > 3 && content[:3] == "```" {
		// 找到第一个换行符
		start := 0
		for i := 3; i < len(content); i++ {
			if content[i] == '\n' {
				start = i + 1
				break
			}
		}
		content = content[start:]
	}

	// 移除结尾的 ```
	if len(content) > 3 && content[len(content)-3:] == "```" {
		content = content[:len(content)-3]
	}

	// 去除首尾空白
	content = strings.TrimSpace(content)

	return content
}

// buildPrompt 构建 Prompt
func buildPrompt(content string, summaryLen int) string {
	return fmt.Sprintf(`请对以下笔记内容进行分析：

1. 生成一段不超过 %d 字的摘要，概括核心内容
2. 提取 3-5 个关键词作为标签建议

笔记内容：
%s

请以 JSON 格式返回(只返回 JSON,不要其他文字)：
{"summary": "摘要内容", "tags": ["标签1", "标签2", "标签3"]}`, summaryLen, content)
}
