package napcat

import (
	"encoding/json"
	"errors"
)

const (
	ActionSendGroupMsg = "send_group_msg"
	TypeText           = "text"
	TypeAt             = "at"
	TypeImage          = "image"
	TypeAudio          = "record"
	TypeFile           = "file"
	TypeVideo          = "video"
	TypeFace           = "face"
	TypeReply          = "reply"
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

// 群文件消息
type GroupFileMsgData struct {
	File string `json:"file"`
	Name string `json:"name"`
}

// 群系统表情消息
type GroupFaceMsgData struct {
	ID int `json:"id"`
}

// 群回复消息
type GroupReplyMsgData struct {
	ID int `json:"id"`
}

// 编码群文本消息
func MarshalGroupTextMsg(groupID int, text string) ([]byte, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid groupID")
	}

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
	if groupID <= 0 {
		return nil, errors.New("invalid groupID")
	}

	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: TypeAt,
					Data: GroupAtMessageData{QQ: qq},
				},
				{
					Type: TypeText,
					Data: GroupTextMsgData{Text: text},
				},
			},
		},
	}

	return json.Marshal(msg)
}

// 编码群语音消息
// path 可以是 url或本地路径
func MarshalGroupAudioMsg(groupID int, path string) ([]byte, error) {
	if groupID <= 0 || path == "" {
		return nil, errors.New("invalid groupID or empty path")
	}

	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: TypeAudio,
					Data: GroupAudioVideoMsgData{File: path},
				},
			},
		},
	}

	return json.Marshal(msg)
}

// 编码群视频消息
// path 可以是 url或本地路径
func MarshalGroupVideoMsg(groupID int, path string) ([]byte, error) {
	if groupID <= 0 || path == "" {
		return nil, errors.New("invalid groupID or empty path")
	}

	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: TypeVideo,
					Data: GroupAudioVideoMsgData{File: path},
				},
			},
		},
	}

	return json.Marshal(msg)
}

// 编码群图片消息
// path 可以是 url或本地路径
func MarshalGroupImgMsg(groupID int, path string) ([]byte, error) {
	if groupID <= 0 || path == "" {
		return nil, errors.New("invalid groupID or empty image URL")
	}

	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: "image",
					Data: GroupImgMsgData{
						File:    path,
						Summary: "[图片]",
					},
				},
			},
		},
	}
	return json.Marshal(msg)
}

// 编码群文件消息
// name是文件名
// path 可以是 url或本地路径
func MarshalGroupFileMsg(groupID int, path string, name string) ([]byte, error) {
	if groupID <= 0 || path == "" || name == "" {
		return nil, errors.New("invalid groupID, empty path or name")
	}

	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: TypeFile,
					Data: GroupFileMsgData{
						File: path,
						Name: name,
					},
				},
			},
		},
	}
	return json.Marshal(msg)
}

// 编码群表情消息
// faceid参考: https://bot.q.qq.com/wiki/develop/api-v2/openapi/emoji/model.html#EmojiType
func MarshalGroupFaceMsg(groupID int, faceID int) ([]byte, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid groupID face ID")
	}

	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: TypeFace,
					Data: GroupFaceMsgData{
						ID: faceID,
					},
				},
			},
		},
	}
	return json.Marshal(msg)
}

// 编码群回复消息
// messageID是回复的消息ID,text是回复的文本
func MarshalGroupReplyMsg(groupID int, messageID int, text string) ([]byte, error) {
	if groupID <= 0 || messageID <= 0 {
		return nil, errors.New("invalid groupID, messageID")
	}

	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: TypeReply,
					Data: GroupReplyMsgData{
						ID: messageID,
					},
				},
				{
					Type: TypeText,
					Data: GroupTextMsgData{
						Text: text,
					},
				},
			},
		},
	}
	return json.Marshal(msg)
}
