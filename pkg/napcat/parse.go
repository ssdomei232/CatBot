package napcat

import "encoding/json"

type Message struct {
	SelfID      int    `json:"self_id"`
	UserID      int    `json:"user_id"`
	Time        int    `json:"time"`
	MessageID   int    `json:"message_id"`
	MessageSeq  int    `json:"message_seq"`
	ReadID      string `json:"read_id"`
	MessageType string `json:"message_type"`
	Sender      Sender `json:"sender"`
	RawMessage  string `json:"raw_message"`
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
	GroupID       int    `json:"group_id"`
	GroupName     string `json:"group_name"`
}

type Sender struct {
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card"`
	Role     string `json:"role"`
}

func Parse(message []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
