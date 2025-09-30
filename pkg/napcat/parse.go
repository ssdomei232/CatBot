package napcat

import "encoding/json"

type Message struct {
	SelfID      int    `json:"self_id"`      // 机器人qq号
	UserID      int    `json:"user_id"`      // 发送者qq号
	Time        int    `json:"time"`         // 消息时间戳 *
	MessageID   int    `json:"message_id"`   // 消息id
	MessageSeq  int    `json:"message_seq"`  // 值和消息id一样，不知道是什么
	ReadID      string `json:"read_id"`      // 值和消息id一样，不知道是什么
	MessageType string `json:"message_type"` // 消息发送位置类型(群/私)
	Sender      Sender `json:"sender"`       // 发送者信息
	RawMessage  string `json:"raw_message"`  // 原始消息 *
	Font        int    `json:"font"`
	SubType     string `json:"sub_type"`
	Message     []struct {
		Type string `json:"type"`
		Data struct {
			Text string `json:"text"`
		} `json:"data"`
	} `json:"message"`
	MessageFormat string `json:"message_format"`
	PostType      string `json:"post_type"`
	GroupID       int    `json:"group_id"`   // 群聊id *
	GroupName     string `json:"group_name"` // 群名 *
}

type Sender struct {
	UserID   int    `json:"user_id"`  // 发送者qq号 *
	Nickname string `json:"nickname"` // 昵称 *
	Card     string `json:"card"`
	Role     string `json:"role"` // 地位
}

func Parse(message []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
