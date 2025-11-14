package napcat

import (
	"encoding/json"
	"errors"
)

// 编码撤回消息
func MarshalDeleteMessage(messageID int) ([]byte, error) {
	message := WSMsg{
		Action: "delete_msg",
		Params: DeleteMsgParams{MessageId: messageID},
	}
	return json.Marshal(message)
}

// 编码群禁言消息
func MarshalGroupBan(groupID int, userID int, duration int) ([]byte, error) {
	message := WSMsg{
		Action: "set_group_ban",
		Params: GroupBanParams{
			GroupID:  groupID,
			UserID:   userID,
			Duration: duration,
		},
	}

	return json.Marshal(message)
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

// 编码群音乐消息
func MarshalGroupMusicMsg(groupID int, musicType string, musicID string) ([]byte, error) {
	if groupID <= 0 || musicID == "" {
		return nil, errors.New("invalid groupID, music ID")
	}

	msg := WSMsg{
		Action: ActionSendGroupMsg,
		Params: GroupMsgParams{
			GroupID: groupID,
			Message: []GroupMsgSegment{
				{
					Type: TypeMusic,
					Data: GroupMusicCardData{
						ID:   musicID,
						Type: musicType,
					},
				},
			},
		},
	}
	return json.Marshal(msg)
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
