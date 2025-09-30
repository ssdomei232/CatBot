package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"git.mmeiblog.cn/mei/CatBot/configs"
)

// 定义请求格式（兼容 OpenAI 格式）
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Client 封装接口调用
type Client struct {
	APIKey  string
	BaseURL string
	Model   string
}

// Send 发送对话请求
func (c *Client) Send(messages []ChatMessage) (string, error) {
	reqBody := ChatRequest{
		Model:    c.Model,
		Messages: messages,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http error: %s, body: %s", resp.Status, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}
	return chatResp.Choices[0].Message.Content, nil
}

func NewClient(message string, prompt string, model string) (string, error) {
	Config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	var c Client
	c.BaseURL = Config.OpenaiLikeUrl
	c.Model = model
	c.APIKey = Config.AiApiKey
	msgs := []ChatMessage{
		{Role: "system", Content: prompt},
		{Role: "user", Content: message},
	}
	response, err := c.Send(msgs)
	if err != nil {
		return "", err
	}
	return response, nil
}
