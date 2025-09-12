package internal

import "git.mmeiblog.cn/mei/aiComplain/pkg/ai"

type SendAndAuditResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
}

func SendAndAudit(message string) (SendAndAuditResponse, error) {
	AuditResponse, err := ai.SendAudit(message)
	var ComplainResponse ai.ComplainMessage
	if err != nil {
		return SendAndAuditResponse{
			Data:    "",
			Message: "审核失败",
		}, err
	}
	if AuditResponse.Approved {
		ComplainResponse, err = ai.SendComplain(message)
		if err != nil {
			return SendAndAuditResponse{
				Data:    "",
				Message: "AI吐槽失败",
			}, err
		}
	}
	return SendAndAuditResponse{
		Data:    ComplainResponse.Taunt,
		Message: "AI吐槽成功",
	}, nil
}
