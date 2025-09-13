package napcat

import "encoding/json"

type SendMessage struct {
	Action string        `json:"action"`
	Params MessageParams `json:"params"`
}

type MessageParams struct {
	GroupID string `json:"group_id"`
	Message string `json:"message"`
}

func Marshal(action string, groupID string, msgType string, msg string) ([]byte, error) {
	var data SendMessage
	data.Action = action
	data.Params.GroupID = groupID
	data.Params.Message = msg
	return json.Marshal(data)
}
