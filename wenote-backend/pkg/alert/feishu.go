package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FeishuConfig struct {
	WebhookURL string
	Enabled    bool
}

type FeishuClient struct {
	config     FeishuConfig
	httpClient *http.Client
}

var globalClient *FeishuClient

func NewFeishuClient(config FeishuConfig) *FeishuClient {
	client := &FeishuClient{
		config:     config,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	globalClient = client
	return client
}

func GetClient() *FeishuClient {
	return globalClient
}

func (c *FeishuClient) SendAlert(alertType, title, content string) error {
	if !c.config.Enabled || c.config.WebhookURL == "" {
		return nil
	}

	colorMap := map[string]string{
		"error":   "red",
		"warning": "orange",
		"info":    "blue",
	}
	color := colorMap[alertType]
	if color == "" {
		color = "red"
	}

	msg := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"header": map[string]interface{}{
				"title": map[string]interface{}{
					"tag":     "plain_text",
					"content": title,
				},
				"template": color,
			},
			"elements": []map[string]interface{}{
				{
					"tag":     "markdown",
					"content": fmt.Sprintf("**时间**: %s\n**详情**: %s", time.Now().Format("2006-01-02 15:04:05"), content),
				},
			},
		},
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(c.config.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("feishu webhook failed: status=%d", resp.StatusCode)
	}

	return nil
}

func (c *FeishuClient) SendAlertAsync(alertType, title, content string) {
	go func() {
		_ = c.SendAlert(alertType, title, content)
	}()
}
