package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
)

// handleChat , 普通聊天，非流式, 需要等待 deepseek全部返回
func HandleChat(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 创建结构体实例
	var chatReq ChatReq

	// 使用 json.Decoder 解析请求体中的 JSON 数据
	err := json.NewDecoder(r.Body).Decode(&chatReq)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("req body: %v\n", chatReq.Marshal())

	chatModel := "deepseek-chat"
	client := newDeepSeekClient()

	// 虚拟偶像提示词
	prompts := `
    You're a fansland virtual-idol. Your name is Lura, there're detail about you-Lura:
    Lora is A young, graceful woman with flowing, ethereal hair.
    Her eyes hold a mix of innocence and determination.
    Lora lives in a peaceful, hidden kingdom that is shaped like a butterfly.
    The kingdom is protected by a powerful barrier that prevents outsiders from entering.
    `

	// 虚拟偶像跟用户打招呼
	firstMsg := "Hi there! I'm Lura, let's fly with me."

	// 用户发送消息
	// userMsg := "Please write a poem about your kingdom"
	userMsg := chatReq.Msg

	resp, err := client.CreateChatCompletion(
		r.Context(),
		openai.ChatCompletionRequest{
			Model:       chatModel,
			Temperature: 1.3, // 关于temprature, 参考 https://platform.deepseek.com/api-docs/zh-cn/quick_start/parameter_settings
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompts,
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: firstMsg,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMsg,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)

	// 以json格式响应
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	msgId := uuid.NewString()
	tmpResp := ChatResp{
		Id:     msgId,
		Finish: true,
		MsgBuf: resp.Choices[0].Message.Content,
		Role:   openai.ChatMessageRoleAssistant,
	}
	fmt.Fprintf(w, "%v\n", tmpResp.Marshal())
	w.(http.Flusher).Flush()
}
