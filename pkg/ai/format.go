package ai

import (
	"log"

	"git.mmeiblog.cn/mei/CatBot/configs"
)

type ComplainMessage struct {
	Taunt string `json:"Taunt"`
}

type AuditMessage struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
}

// SendComplain 返还ai的话
func SendComplain(message string) (response string, err error) {
	Config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	response, err = NewClient(message, Config.Prompt, Config.LLMModel)
	if err != nil {
		return "", err
	}
	return response, nil
}
