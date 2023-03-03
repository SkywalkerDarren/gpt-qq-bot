package cqhttp

import (
	"OpenAIBot/src/module/cqhttp/model"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type CQHttp struct {
	cfg     *model.CQHttpConfig
	ws      *websocket.Conn
	handler Handler
	action  Action
}

func NewCQHttp(cfg *model.CQHttpConfig, handler Handler) *CQHttp {
	return &CQHttp{cfg: cfg, handler: handler, action: Action{}}
}

func (c *CQHttp) Run() error {
	u := url.URL{Scheme: "ws", Host: c.cfg.Host, Path: c.cfg.Path}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": {"Bearer " + c.cfg.Token}})
	if err != nil {
		return err
	}
	c.action.ws = ws
	c.ws = ws

	go c.wsHandler()

	return nil
}

func (c *CQHttp) wsHandler() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
		}

		var event model.Event
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println("json unmarshal:", err)
			continue
		}
		c.handler.OnEvent(event, c.action)
		if event.PostType == model.PostTypeMessage {
			var msg model.MessageEvent
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("json unmarshal:", err)
				continue
			}
			c.handler.OnMessage(msg, c.action)
			if msg.MessageType == model.MessageTypeGroup {
				var groupMsg model.MessageGroupEvent
				err = json.Unmarshal(message, &groupMsg)
				if err != nil {
					log.Println("json unmarshal:", err)
					continue
				}
				c.handler.OnMessageGroup(groupMsg, c.action)
			} else if msg.MessageType == model.MessageTypePrivate {
				var privateMsg model.MessagePrivateEvent
				err = json.Unmarshal(message, &privateMsg)
				if err != nil {
					log.Println("json unmarshal:", err)
					continue
				}
				c.handler.OnMessagePrivate(privateMsg, c.action)
			}
		} else if event.PostType == model.PostTypeNotice {
			var notice model.NoticeEvent
			err = json.Unmarshal(message, &notice)
			if err != nil {
				log.Println("json unmarshal:", err)
				continue
			}
			c.handler.OnNotice(notice, c.action)
		} else if event.PostType == model.PostTypeRequest {
			var request model.RequestEvent
			err = json.Unmarshal(message, &request)
			if err != nil {
				log.Println("json unmarshal:", err)
				continue
			}
			c.handler.OnRequest(request, c.action)
		} else if event.PostType == model.PostTypeMetaEvent {
			var metaEvent model.MetaEvent
			err = json.Unmarshal(message, &metaEvent)
			if err != nil {
				log.Println("json unmarshal:", err)
				continue
			}
			c.handler.OnMetaEvent(metaEvent, c.action)
			if metaEvent.MetaEventType == model.MetaEventTypeHeartbeat {
				var metaHeartbeat model.MetaHeartbeatEvent
				err = json.Unmarshal(message, &metaHeartbeat)
				if err != nil {
					log.Println("json unmarshal:", err)
					continue
				}
				c.handler.OnMetaHeartbeat(metaHeartbeat, c.action)
			} else if metaEvent.MetaEventType == model.MetaEventTypeLifecycle {
				var metaLifecycle model.MetaLifecycleEvent
				err = json.Unmarshal(message, &metaLifecycle)
				if err != nil {
					log.Println("json unmarshal:", err)
					continue
				}
				c.handler.OnMetaLifecycle(metaLifecycle, c.action)
			}
		}
	}
}

func (c *CQHttp) Close() {
	c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.ws.Close()
}
