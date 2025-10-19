package napcat

import "encoding/json"

type DeleteMsgParams struct {
	MessageId int `json:"message_id"`
}

// 编码撤回消息
func MarshalDeleteMessage(messageID int) ([]byte, error) {
	message := WSMsg{
		Action: "delete_msg",
		Params: DeleteMsgParams{MessageId: messageID},
	}
	return json.Marshal(message)
}
