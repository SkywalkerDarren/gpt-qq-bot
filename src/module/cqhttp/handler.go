package cqhttp

import "OpenAIBot/src/module/cqhttp/model"

type Handler interface {
	OnEvent(event model.Event, action Action)
	OnMessage(msg model.MessageEvent, action Action)
	OnMessageGroup(msg model.MessageGroupEvent, action Action)
	OnMessagePrivate(msg model.MessagePrivateEvent, action Action)
	OnNotice(notice model.NoticeEvent, action Action)
	OnRequest(request model.RequestEvent, action Action)
	OnMetaEvent(metaEvent model.MetaEvent, action Action)
	OnMetaHeartbeat(metaHeartbeat model.MetaHeartbeatEvent, action Action)
	OnMetaLifecycle(metaLifecycle model.MetaLifecycleEvent, action Action)
}

type DefaultHandler struct {
}

func (d *DefaultHandler) OnEvent(event model.Event, action Action) {

}

func (d *DefaultHandler) OnMessage(msg model.MessageEvent, action Action) {

}

func (d *DefaultHandler) OnMessageGroup(msg model.MessageGroupEvent, action Action) {

}

func (d *DefaultHandler) OnMessagePrivate(msg model.MessagePrivateEvent, action Action) {

}

func (d *DefaultHandler) OnNotice(notice model.NoticeEvent, action Action) {

}

func (d *DefaultHandler) OnRequest(request model.RequestEvent, action Action) {

}

func (d *DefaultHandler) OnMetaEvent(metaEvent model.MetaEvent, action Action) {

}

func (d *DefaultHandler) OnMetaHeartbeat(metaHeartbeat model.MetaHeartbeatEvent, action Action) {

}

func (d *DefaultHandler) OnMetaLifecycle(metaLifecycle model.MetaLifecycleEvent, action Action) {

}

var _ Handler = (*DefaultHandler)(nil)
