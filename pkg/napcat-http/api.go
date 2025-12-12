package napcathttp

func (c *Client) SendGroupMsg(groupID string, message string) error {
	req := GroupMessage{
		GroupID: groupID,
		Message: []Message{
			{
				Type: "text",
				Data: struct {
					Text string `json:"text"`
				}{Text: message},
			},
		},
	}
	err := c.DoRequest("POST", "/send_group_msg", req, nil)

	return err
}
