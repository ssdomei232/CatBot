package napcat

import (
	"encoding/json"
	"errors"
)

const (
	ActionSendGroupMsg = "send_group_msg"
	SummaryImage       = "[图片]"
)

type GroupImgMsg struct {
	Action string `json:"action"`
	Params struct {
		GroupID int                 `json:"group_id"`
		Message []GroupImgMsgParams `json:"message"`
	} `json:"params"`
}

type GroupImgMsgParams struct {
	Type string `json:"type"`
	Data struct {
		File    string `json:"file"`
		Summary string `json:"summary"`
	} `json:"data"`
}

func MarshalGroupImgMsg(groupID int, imgUrl string) ([]byte, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid group id")
	}
	if imgUrl == "" {
		return nil, errors.New("image url is empty")
	}

	groupImgMsg := GroupImgMsg{
		Action: ActionSendGroupMsg,
	}
	groupImgMsg.Params.GroupID = groupID
	groupImgMsg.Params.Message = []GroupImgMsgParams{
		{
			Type: "image",
			Data: struct {
				File    string `json:"file"`
				Summary string `json:"summary"`
			}{
				File:    imgUrl,
				Summary: SummaryImage,
			},
		},
	}
	return json.Marshal(groupImgMsg)
}
