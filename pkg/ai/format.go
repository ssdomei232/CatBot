package ai

import (
	"encoding/json"

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
func SendComplain(message string) (ComplainMessage, error) {
	response, err := NewClient(message, configs.COMPLAIN_PROMPT_TEMPLATE, configs.LLM_MODEL_LARGE)
	if err != nil {
		return ComplainMessage{}, err
	}
	var format ComplainMessage
	err = json.Unmarshal([]byte(response), &format)
	if err != nil {
		return ComplainMessage{}, err
	}
	return format, nil
}

// SendAudit 检测内容是否合规
func SendAudit(message string) (AuditMessage, error) {
	response, err := NewClient(message, configs.AUDIT_PROMPT_TEMPLATE, configs.LLM_MODEL_TINY)
	if err != nil {
		return AuditMessage{}, err
	}
	var format AuditMessage
	err = json.Unmarshal([]byte(response), &format)
	if err != nil {
		return AuditMessage{}, err
	}
	return format, nil
}
