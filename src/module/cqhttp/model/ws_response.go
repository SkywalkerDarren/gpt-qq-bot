package model

// {"post_type":"message","message_type":"group","time":1677757343,"self_id":2869015936,"sub_type":"normal","anonymous":null,"group_id":734130058,"message_seq":11733,"user_id":477720730,"message_id":-1155659933,"font":0,"message":"啊","raw_message":"啊","sender":{"age":0,"area":"","card":"","level":"","nickname":"SkywalkerDarren","role":"owner","sex":"unknown","title":"","user_id":477720730}}
// {"post_type":"message","message_type":"private","time":1677757663,"self_id":2869015936,"sub_type":"friend","user_id":477720730,"target_id":2869015936,"message":"我","raw_message":"我","font":0,"sender":{"age":0,"nickname":"SkywalkerDarren","sex":"unknown","user_id":477720730},"message_id":1950249556}
// {"post_type":"message","message_type":"private","time":1677758385,"self_id":2869015936,"sub_type":"friend","raw_message":"[CQ:image,file=0a73b62c9c8c9afab69295d99a76929a.image,url=https://c2cpicdw.qpic.cn/offpic_new/477720730//477720730-1939516307-0A73B62C9C8C9AFAB69295D99A76929A/0?term=2\u0026amp;is_origin=0]","font":0,"sender":{"age":0,"nickname":"SkywalkerDarren","sex":"unknown","user_id":477720730},"message_id":-726405240,"user_id":477720730,"target_id":2869015936,"message":"[CQ:image,file=0a73b62c9c8c9afab69295d99a76929a.image,url=https://c2cpicdw.qpic.cn/offpic_new/477720730//477720730-1939516307-0A73B62C9C8C9AFAB69295D99A76929A/0?term=2\u0026amp;is_origin=0]"}
type Event struct {
	PostType PostType `json:"post_type"`
}

type PostType string

const (
	PostTypeMessage   PostType = "message"
	PostTypeRequest   PostType = "request"
	PostTypeNotice    PostType = "notice"
	PostTypeMetaEvent PostType = "meta_event"
)

type MessageEvent struct {
	Event
	MessageType MessageType `json:"message_type"`
	Time        int64       `json:"time"`
	SelfID      int64       `json:"self_id"`
	SubType     SubType     `json:"sub_type"`
	UserID      int64       `json:"user_id"`
	Message     string      `json:"message"`
}

type MessageType string

const (
	MessageTypeGroup   MessageType = "group"
	MessageTypePrivate MessageType = "private"
)

type SubType string

const (
	SubTypeFriend    SubType = "friend"
	SubTypeNormal    SubType = "normal"
	SubTypeAnonymous SubType = "anonymous"
	SubTypeGroupSelf SubType = "group_self"
	SubTypeGroup     SubType = "group"
	SubTypeNotice    SubType = "notice"
)

type MessageGroupEvent struct {
	MessageEvent
	GroupID    int64 `json:"group_id"`
	MessageSeq int64 `json:"message_seq"`
	MessageID  int64 `json:"message_id"`
}

type MessagePrivateEvent struct {
	MessageEvent
	TargetID int64 `json:"target_id"`
}

type Sender struct {
	Age      int64  `json:"age"`
	Area     string `json:"area"`
	Card     string `json:"card"`
	Level    string `json:"level"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Sex      string `json:"sex"`
	Title    string `json:"title"`
	UserID   int64  `json:"user_id"`
}

type NoticeEvent struct {
	Event
}

type RequestEvent struct {
	Event
}

// recv: {"_post_method":2,"meta_event_type":"lifecycle","post_type":"meta_event","self_id":2869015936,"sub_type":"connect","time":1677750317}
// recv: {"post_type":"meta_event","meta_event_type":"heartbeat","time":1677758872,"self_id":2869015936,"status":
// {"app_enabled":true,"app_good":true,"app_initialized":true,"good":true,"online":true,"plugins_good":null,"stat":
// {"packet_received":494,"packet_sent":464,"packet_lost":0,"message_received":20,"message_sent":9,"disconnect_times":0,"lost_times":0,"last_message_time":1677758385}},"interval":5000}

type MetaEvent struct {
	Event
	MetaEventType MetaEventType `json:"meta_event_type"`
	SelfID        int64         `json:"self_id"`
	Time          int64         `json:"time"`
}

type MetaEventType string

const (
	MetaEventTypeHeartbeat MetaEventType = "heartbeat"
	MetaEventTypeLifecycle MetaEventType = "lifecycle"
)

type MetaLifecycleEvent struct {
	MetaEvent
	SubType LifecycleSubType `json:"sub_type"`
}

type LifecycleSubType string

const (
	LifecycleSubTypeConnect    LifecycleSubType = "connect"
	LifecycleSubTypeDisconnect LifecycleSubType = "disconnect"
)

type MetaHeartbeatEvent struct {
	MetaEvent
	Status Status `json:"status"`
}

type Status struct {
	AppEnabled     bool        `json:"app_enabled"`
	AppGood        bool        `json:"app_good"`
	AppInitialized bool        `json:"app_initialized"`
	Good           bool        `json:"good"`
	Online         bool        `json:"online"`
	PluginsGood    interface{} `json:"plugins_good"`
	Stat           Stat        `json:"stat"`
	Interval       int64       `json:"interval"`
}

type Stat struct {
	PacketReceived  int64 `json:"packet_received"`
	PacketSent      int64 `json:"packet_sent"`
	PacketLost      int64 `json:"packet_lost"`
	MessageReceived int64 `json:"message_received"`
	MessageSent     int64 `json:"message_sent"`
	DisconnectTimes int64 `json:"disconnect_times"`
	LostTimes       int64 `json:"lost_times"`
	LastMessageTime int64 `json:"last_message_time"`
}

type ActionResponse struct {
}
