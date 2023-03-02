package model

type WebsocketActionRequest struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo"`
}

// action enum
const (
	SendPrivateMsg = "send_private_msg"
	SendGroupMsg   = "send_group_msg"
)

type SendPrivateMsgParams struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

type SendGroupMsgParams struct {
	GroupID int64  `json:"group_id"`
	Message string `json:"message"`
}
