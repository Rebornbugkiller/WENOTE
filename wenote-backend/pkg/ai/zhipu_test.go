package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestZhipuClientSuccess 测试成功请求
func TestZhipuClientSuccess(t *testing.T) {
	// Mock 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := zhipuResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{
					Message: struct {
						Content string `json:"content"`
					}{
						Content: `{"summary":"测试摘要","tags":["Go","测试"]}`,
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewZhipuClient(ZhipuConfig{
		APIKey:  "test-key",
		Model:   "glm-4-flash",
		BaseURL: server.URL,
	})

	ctx := context.Background()
	result, err := client.GenerateSummaryAndTags(ctx, "测试内容", 100)

	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if result.Summary != "测试摘要" {
		t.Errorf("Expected summary='测试摘要', got='%s'", result.Summary)
	}

	if len(result.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(result.Tags))
	}

	if result.Tags[0] != "Go" || result.Tags[1] != "测试" {
		t.Errorf("Unexpected tags: %v", result.Tags)
	}
}

// TestZhipuClientRetry 测试重试机制
func TestZhipuClientRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 第三次成功
		response := zhipuResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{
					Message: struct {
						Content string `json:"content"`
					}{
						Content: `{"summary":"重试成功","tags":["重试"]}`,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewZhipuClient(ZhipuConfig{
		APIKey:  "test-key",
		Model:   "glm-4-flash",
		BaseURL: server.URL,
	})

	ctx := context.Background()
	result, err := client.GenerateSummaryAndTags(ctx, "测试内容", 100)

	if err != nil {
		t.Fatalf("Expected success after retry, got error: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	if result.Summary != "重试成功" {
		t.Errorf("Unexpected summary: %s", result.Summary)
	}
}

// TestZhipuClientTimeout 测试超时
func TestZhipuClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewZhipuClient(ZhipuConfig{
		APIKey:  "test-key",
		Model:   "glm-4-flash",
		BaseURL: server.URL,
		Timeout: 1, // 1秒超时
	})

	ctx := context.Background()
	_, err := client.GenerateSummaryAndTags(ctx, "测试内容", 100)

	if err == nil {
		t.Error("Expected timeout error")
	}
}

// TestZhipuClientAPIError 测试 API 错误
func TestZhipuClientAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := zhipuResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "INVALID_KEY",
				Message: "Invalid API key",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewZhipuClient(ZhipuConfig{
		APIKey:  "invalid-key",
		Model:   "glm-4-flash",
		BaseURL: server.URL,
	})

	ctx := context.Background()
	_, err := client.GenerateSummaryAndTags(ctx, "测试内容", 100)

	if err == nil {
		t.Error("Expected API error")
	}

	if err.Error() != "all retries failed: API error: Invalid API key" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestBuildPrompt 测试 Prompt 构建
func TestBuildPrompt(t *testing.T) {
	content := "这是测试内容"
	summaryLen := 100

	prompt := buildPrompt(content, summaryLen)

	if prompt == "" {
		t.Error("Prompt should not be empty")
	}

	// 检查是否包含关键信息
	if !contains(prompt, "100") {
		t.Error("Prompt should contain summary length")
	}

	if !contains(prompt, content) {
		t.Error("Prompt should contain note content")
	}

	if !contains(prompt, "JSON") {
		t.Error("Prompt should mention JSON format")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
