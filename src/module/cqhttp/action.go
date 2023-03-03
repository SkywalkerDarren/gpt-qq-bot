package cqhttp

import (
	"OpenAIBot/src/module/cqhttp/model"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Action struct {
	ws *websocket.Conn
}

func (a *Action) SendGroupMsg(groupID int64, msg string) {
	err := a.ws.WriteJSON(model.NewSendGroupMsg(groupID, msg))
	if err != nil {
		log.Println("send group msg:", err)
		return
	}
}

func (a *Action) SendPrivateMsg(userID int64, msg string) {
	data, err := json.Marshal(model.NewSendPrivateMsg(userID, msg))
	log.Println("SendPrivateMsg", string(data))
	if err != nil {
		log.Println("send private msg:", err)
		return
	}
	err = a.ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("send private msg:", err)
		return
	}
}
