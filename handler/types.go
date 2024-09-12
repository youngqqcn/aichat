package handler

import (
	"encoding/json"
)

// 定义与 JSON 数据匹配的结构体
type ChatReq struct {
	Msg string `json:"msg"`
}

func (r *ChatReq) Marshal() string {
	b, _ := json.Marshal(r)
	return string(b)
}

type ChatResp struct {
	Id     string `json:"id"`
	Role   string `json:"role"`
	MsgBuf string `json:"msgBuf"` // 包含所有已完成的msg
	Delta  string `json:"delta"`  // 新增
	Finish bool   `json:"finish"`
}

func (r *ChatResp) Marshal() string {
	b, _ := json.Marshal(r)
	return string(b)
}
