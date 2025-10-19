package napcat

import (
	"encoding/json"
)

const (
	ActionSendGroupMsg = "send_group_msg"
)

// Websocket 消息基本结构
type WSMsg struct {
	Action string `json:"action"`
	Params any    `json:"params"`
}

// websocket 连接比较特别的发送文本消息的方式
type GroupTextMsgParams struct {
	GroupID int    `json:"group_id"`
	Message string `json:"message"`
}

// 群消息基本结构
type GroupMsgParams struct {
	GroupID int `json:"group_id"`
	Message any `json:"message"`
}

// 群消息内容基本结构
type GroupMsgSegment struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// 群@消息
type GroupAtMessageData struct {
	QQ int `json:"qq"`
}

// 群文本消息
type GroupTextMsgData struct {
	Text string `json:"text"`
}

// 群语音/视频消息(这两种消息的Data是一样的)
type GroupAudioVideoMsgData struct {
	File string `json:"file"`
}

// 群图片消息
type GroupImgMsgData struct {
	File    string `json:"file"`
	Summary string `json:"summary"`
}

// 编码群文本消息
func MarshalGroupTextMsg(groupID int, text string) ([]byte, error) {
	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupTextMsgParams{
			GroupID: groupID,
			Message: text,
		},
	}
	return json.Marshal(msg)
}

// 编码群@消息
func MarshalAtMsg(groupID int, qq int, text string) ([]byte, error) {
	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: "at",
					Data: GroupAtMessageData{QQ: qq},
				},
				{
					Type: "text",
					Data: GroupTextMsgData{Text: text},
				},
			},
		},
	}

	return json.Marshal(msg)
}

// 编码群语音消息
func MarshalGroupAudioMsg(groupID int, path string) ([]byte, error) {
	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: "record",
					Data: GroupAudioVideoMsgData{File: path},
				},
			},
		},
	}

	return json.Marshal(msg)
}

// 编码群图片消息
func MarshalGroupImgMsg(groupID int, imgUrl string) ([]byte, error) {
	groupImgMsg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: "image",
					Data: GroupImgMsgData{
						File:    imgUrl,
						Summary: "[图片]",
					},
				},
			},
		},
	}
	return json.Marshal(groupImgMsg)
}
