package handler

import (
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// 创建 DeepSeek 的 Client
// 参考deepseek的官方文档: https://platform.deepseek.com/api-docs/zh-cn/
func newDeepSeekClient() *openai.Client {
	authToken := os.Getenv("DEEPSEEK_API_KEY")
	baseUrl := "https://api.deepseek.com"
	cfg := openai.DefaultConfig(authToken)
	cfg.BaseURL = baseUrl
	return openai.NewClientWithConfig(cfg)
}
