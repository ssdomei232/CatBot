package configs

import (
	"encoding/json"
	"os"
)

type Configs struct {
	OpenaiLikeUrl       string `json:"openai_like_url"`
	AiApiKey            string `json:"ai_api_key"`
	NapcatHost          string `json:"napcat_host"`
	NapcatToken         string `json:"napcat_token"`
	LLMModel            string `json:"llm_model"`
	NapcatWebsocketPort int    `json:"napcat_websocket_port"`
	NapcatHttpPort      int    `json:"napcat_http_port"`
	WebhookSecret       string `json:"webhook_secret"`
	Prompt              string `json:"prompt"`
}

func GetConfig() (config Configs, err error) {
	content, err := os.ReadFile("config.json")
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
