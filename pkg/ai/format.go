package ai

import (
	"git.mmeiblog.cn/mei/aiComplain/configs"
)

type ComplainMessage struct {
	Taunt string `json:"Taunt"`
}

type AuditMessage struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
}

// SendComplain 发送返回一个ai的吐槽内容
func SendComplain(message string) (response string, err error) {
	response, err = NewClient(message, configs.COMPLAIN_PROMPT_TEMPLATE, configs.LLM_MODEL_LARGE)
	if err != nil {
		return "", err
	}
	return response, nil
}
