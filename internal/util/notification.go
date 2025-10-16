package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vearne/autotest/internal/config"
	slog "github.com/vearne/simplelog"
)

// NotificationService 通知服务
type NotificationService struct {
	config config.AutoTestConfig
	client *http.Client
}

// TestResult 测试结果
type TestResult struct {
	TotalTests  int           `json:"total_tests"`
	PassedTests int           `json:"passed_tests"`
	FailedTests int           `json:"failed_tests"`
	Duration    time.Duration `json:"duration"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	FailedCases []string      `json:"failed_cases,omitempty"`
}

// SlackMessage Slack消息格式
type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment Slack附件
type Attachment struct {
	Color  string  `json:"color"`
	Title  string  `json:"title"`
	Fields []Field `json:"fields"`
}

// Field Slack字段
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// NewNotificationService 创建通知服务
func NewNotificationService(cfg config.AutoTestConfig) *NotificationService {
	return &NotificationService{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendTestResult 发送测试结果通知
func (ns *NotificationService) SendTestResult(result TestResult) error {
	if !ns.config.Global.Notifications.Enabled {
		return nil
	}

	isSuccess := result.FailedTests == 0

	// 检查是否需要发送通知
	if isSuccess && !ns.config.Global.Notifications.OnSuccess {
		return nil
	}
	if !isSuccess && !ns.config.Global.Notifications.OnFailure {
		return nil
	}

	return ns.sendSlackNotification(result)
}

// sendSlackNotification 发送Slack通知
func (ns *NotificationService) sendSlackNotification(result TestResult) error {
	if ns.config.Global.Notifications.WebhookURL == "" {
		slog.Warn("Webhook URL not configured, skipping notification")
		return nil
	}

	isSuccess := result.FailedTests == 0
	color := "good" // 绿色
	status := "✅ 成功"

	if !isSuccess {
		color = "danger" // 红色
		status = "❌ 失败"
	}

	message := SlackMessage{
		Text: fmt.Sprintf("AutoTest 执行完成 - %s", status),
		Attachments: []Attachment{
			{
				Color: color,
				Title: "测试结果详情",
				Fields: []Field{
					{
						Title: "总测试数",
						Value: fmt.Sprintf("%d", result.TotalTests),
						Short: true,
					},
					{
						Title: "通过数",
						Value: fmt.Sprintf("%d", result.PassedTests),
						Short: true,
					},
					{
						Title: "失败数",
						Value: fmt.Sprintf("%d", result.FailedTests),
						Short: true,
					},
					{
						Title: "执行时间",
						Value: result.Duration.String(),
						Short: true,
					},
					{
						Title: "开始时间",
						Value: result.StartTime.Format("2006-01-02 15:04:05"),
						Short: true,
					},
					{
						Title: "结束时间",
						Value: result.EndTime.Format("2006-01-02 15:04:05"),
						Short: true,
					},
				},
			},
		},
	}

	// 如果有失败的测试用例，添加失败详情
	if len(result.FailedCases) > 0 {
		failedCasesText := ""
		for i, failedCase := range result.FailedCases {
			if i > 0 {
				failedCasesText += "\n"
			}
			failedCasesText += fmt.Sprintf("• %s", failedCase)

			// 限制显示的失败用例数量
			if i >= 9 {
				failedCasesText += fmt.Sprintf("\n... 还有 %d 个失败用例", len(result.FailedCases)-10)
				break
			}
		}

		message.Attachments[0].Fields = append(message.Attachments[0].Fields, Field{
			Title: "失败用例",
			Value: failedCasesText,
			Short: false,
		})
	}

	return ns.sendWebhook(message)
}

// sendWebhook 发送Webhook请求
func (ns *NotificationService) sendWebhook(message SlackMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal notification message: %w", err)
	}

	req, err := http.NewRequest("POST", ns.config.Global.Notifications.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ns.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook returned status code: %d", resp.StatusCode)
	}

	slog.Info("Notification sent successfully")
	return nil
}
