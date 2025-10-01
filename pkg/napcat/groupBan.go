package napcat

import "encoding/json"

type GroupBanMessage struct {
	Action string         `json:"action"`
	Params GroupBanParams `json:"params"`
}

type GroupBanParams struct {
	GroupID  int `json:"group_id"`
	UserID   int `json:"user_id"`
	Duration int `json:"duration"`
}

func MarshalGroupBan(groupID int, userID int, duration int) ([]byte, error) {
	message := GroupBanMessage{
		Action: "set_group_ban",
		Params: GroupBanParams{
			GroupID:  groupID,
			UserID:   userID,
			Duration: duration,
		},
	}

	return json.Marshal(message)
}
