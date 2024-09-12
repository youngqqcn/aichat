package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
)

// handleChatSSE 流式聊天,
// 使用 SSE, 与deepseek也是SSE, 与客户端也是SSE
// 关于SSE(Server-Sents Event)
// https://en.wikipedia.org/wiki/Server-sent_events
func HandleChatStream(w http.ResponseWriter, r *http.Request) {
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

	// 发送请求到 deepseek
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

	// 创建流式聊天
	stream, err := client.CreateChatCompletionStream(
		r.Context(), // 与请求共用同一个 context, 当用户断开连接时，与deepseek的连接会同时断开
		openai.ChatCompletionRequest{
			Model:       chatModel,
			Stream:      true, // 流
			Temperature: 1.3,  // 关于temprature, 参考 https://platform.deepseek.com/api-docs/zh-cn/quick_start/parameter_settings
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
	defer stream.Close()

	// 手动设置响应头
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 返回数据

	// 统一用一个msgId
	msgId := uuid.NewString()
	msgBuf := make([]byte, 0)
	for {
		// select {
		// case <-r.Context().Done():
		// 	fmt.Println("客户端已断开，同时已断开与DeepSeek的连接")
		// 	return
		// default:
		// }

		response, err := stream.Recv()

		tmpResp := ChatResp{
			Id:     msgId,
			Finish: false,
			MsgBuf: string(msgBuf),
			Role:   openai.ChatMessageRoleAssistant,
		}

		// 处理结束
		if errors.Is(err, io.EOF) {
			fmt.Println()

			tmpResp.Finish = true
			tmpResp.Delta = ""
			tmpResp.MsgBuf = string(msgBuf)

			fmt.Fprintf(w, "%v", tmpResp.Marshal())
			w.(http.Flusher).Flush()
			return
		}

		// 处理客户端断开
		if errors.Is(err, context.Canceled) {
			fmt.Println("客户端已断开，同时已断开与DeepSeek的连接")
			return
		}

		if err != nil {

			fmt.Printf("Stream error: %v\n", err)
			return
		}

		// 打印新增的消息
		fmt.Printf("%v", response.Choices[0].Delta.Content)

		// 追加新增内容
		msgBuf = append(msgBuf, []byte(response.Choices[0].Delta.Content)...)

		tmpResp.Delta = response.Choices[0].Delta.Content
		tmpResp.MsgBuf = string(msgBuf)
		tmpResp.Finish = false

		// 只返回增量文本
		// fmt.Fprintf(w, "%v", tmpResp.Delta)

		// 以json格式返回
		fmt.Fprintf(w, "%v", tmpResp.Marshal())

		// 为了在postman中显示SSE的消息, https://github.com/postmanlabs/postman-app-support/issues/12448
		// fmt.Fprintf(w, "data: %v \n\n", tmpResp.Marshal())

		// 刷新
		w.(http.Flusher).Flush()
	}

}
