# aichat

使用 deepseek 作为 AI 聊天大模型， deepseek 性价比非常高

-   DeepSeek API 文档: https://platform.deepseek.com/api-docs/zh-cn/
-   DeepSeek 开发者平台： https://platform.deepseek.com/

## 数据流

![](./arch.png)

## 启动服务

-   新建`.env`文件, 设置 `DEEPSEEK_API_KEY`

-   启动服务：

```
go run server/main.go
```

-   启动客户端(非流式), 需要一直等待全部完成:

让 ai 写一首诗

```
curl -s -N -X POST localhost:8080/chat --data '
{
    "msg":"Please write a poem about your kingdom"
}' | jq  --unbuffered

```

响应示例：

```json
{
    "id": "7380f171-c05d-43bf-a6f4-56138356a8ba",
    "role": "assistant",
    "msgBuf": "In a land where dreams take flight,\nLies my kingdom, pure and bright.\nA butterfly's gentle wings,\nShaped our realm, where wonders sing.\n\nGolden fields of shimmering grace,\nWhere the rivers softly trace\nPatterns of the stars above,\nIn this haven, peace is love.\n\nThe barrier's power, strong and true,\nKeeps our secrets safe from view.\nYet within, the heart's content,\nWhere every soul is heaven-sent.\n\nLush gardens bloom with every hue,\nWhispering tales of ancient blue.\nMountains rise in emerald green,\nGuarding secrets, yet unseen.\n\nIn the heart of this butterfly,\nLives a people, kind and free.\nWith laughter bright and songs so sweet,\nIn our kingdom, joy we greet.\n\nSo come, dear friend, let's take a stroll,\nThrough the land where dreams unroll.\nIn the butterfly's gentle hold,\nFind a peace that's never old.",
    "delta": "",
    "finish": true
}
```

-   启动客户端（流式）, 立即响应

```
curl -s -N -X POST localhost:8080/chat-stream --data '
{
    "msg":"Please write a poem about your kingdom"
}' | jq  --unbuffered
```
