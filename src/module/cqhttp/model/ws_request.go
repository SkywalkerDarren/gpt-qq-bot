package model

import "github.com/google/uuid"

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

func NewSendPrivateMsg(userID int64, message string) *WebsocketActionRequest {
	return &WebsocketActionRequest{
		Action: SendPrivateMsg,
		Params: SendPrivateMsgParams{
			UserID:  userID,
			Message: message,
		},
		Echo: uuid.NewString(),
	}
}

type SendGroupMsgParams struct {
	GroupID int64  `json:"group_id"`
	Message string `json:"message"`
}

func NewSendGroupMsg(groupID int64, message string) *WebsocketActionRequest {
	return &WebsocketActionRequest{
		Action: SendGroupMsg,
		Params: SendGroupMsgParams{
			GroupID: groupID,
			Message: message,
		},
		Echo: uuid.NewString(),
	}
}
