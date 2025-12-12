package napcathttp

type GroupMessage struct {
	GroupID string    `json:"group_id"`
	Message []Message `json:"message"`
}

type Message struct {
	Type string `json:"type"`
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}
