package napcat

import "encoding/json"

type DeleteMsgParams struct {
	MessageId int `json:"message_id"`
}

func MarshalDeleteMessage(messageID int) ([]byte, error) {
	message := WSMsg{
		Action: "delete_msg",
		Params: DeleteMsgParams{MessageId: messageID},
	}
	return json.Marshal(message)
}
