package napcat

import "encoding/json"

type DeleteMessage struct {
	Action string `json:"action"`
	Params struct {
		MessageId int `json:"message_id"`
	} `json:"params"`
}

func MarshalDeleteMessage(messageID int) ([]byte, error) {
	message := DeleteMessage{
		Action: "delete_msg",
		Params: struct {
			MessageId int `json:"message_id"`
		}{MessageId: messageID},
	}
	return json.Marshal(message)
}
