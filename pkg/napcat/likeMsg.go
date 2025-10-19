package napcat

import "encoding/json"

type LikeMsgParams struct {
	UserID int `json:"user_id"`
	Times  int `json:"times"`
}

// 编码点赞消息
func MarshalLikeMsg(userId int, times int) ([]byte, error) {
	likeMsg := WSMsg{
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
