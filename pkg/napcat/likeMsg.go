package napcat

import "encoding/json"

type LikeMsg struct {
	Action string `json:"action"`
	Params struct {
		UserId int `json:"user_id"`
		Times  int `json:"times"`
	} `json:"params"`
}

func MarshalLikeMsg(userId int, times int) ([]byte, error) {
	likeMsg := LikeMsg{
		Action: "send_like",
		Params: struct {
			UserId int `json:"user_id"`
			Times  int `json:"times"`
		}{
			UserId: userId,
			Times:  times,
		},
	}
	return json.Marshal(likeMsg)
}
