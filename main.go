package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/youngqqcn/aichat_go/handler"
)

func main() {
    // 如果.env不存在则创建
	_, err := os.Stat(".env")
	if err != nil {
		if os.IsNotExist(err) {
			f, e := os.Create(".env")
			if e == nil {
				f.WriteString("DEEPSEEK_API_KEY=xxx")
				f.Close()
			}
		}
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("缺少 .env 文件")
	}

	// 参考deepseek的官方文档: https://platform.deepseek.com/api-docs/zh-cn/
	if !strings.HasPrefix(os.Getenv("DEEPSEEK_API_KEY"), "sk-") {
		log.Fatal("请在 .env中填写正确deepseek的api key, DEEPSEEK_API_KEY")
	}

	http.HandleFunc("/chat-stream", handler.HandleChatStream)
	http.HandleFunc("/chat", handler.HandleChat)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
