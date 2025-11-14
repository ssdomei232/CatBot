package napcat

type DeleteMsgParams struct {
	MessageId int `json:"message_id"`
}

type GroupBanParams struct {
	GroupID  int `json:"group_id"`
	UserID   int `json:"user_id"`
	Duration int `json:"duration"`
}

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
	TypeMusic          = "music"
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

// 群音乐卡片消息
type GroupMusicCardData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type LikeMsgParams struct {
	UserID int `json:"user_id"`
	Times  int `json:"times"`
}
