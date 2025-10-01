package napcat

import (
	"encoding/json"
)

type AtMessageData struct {
	QQ int `json:"qq"`
}

type MessageData struct {
	Text string `json:"text"`
}

type MessageSegment struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type GroupMessage struct {
	Action string             `json:"action"`
	Params GroupMessageParams `json:"params"`
}

type GroupMessageParams struct {
	GroupID int              `json:"group_id"`
	Message []MessageSegment `json:"message"`
}

func MarshalAtMsg(groupID int, qq int, text string) ([]byte, error) {
	msg := GroupMessage{
		Action: "send_group_msg",
		Params: GroupMessageParams{
			GroupID: groupID,
			Message: []MessageSegment{
				{
					Type: "at",
					Data: AtMessageData{QQ: qq},
				},
				{
					Type: "text",
					Data: MessageData{Text: text},
				},
			},
		},
	}

	return json.Marshal(msg)
}
