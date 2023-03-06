package controller

import (
	"OpenAIBot/src/module/cqhttp"
	"OpenAIBot/src/module/cqhttp/model"
	"context"
	"log"
	"strings"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type ChatBot struct {
	cqhttp.DefaultHandler
	openAI             *gogpt.Client
	msgsGroupMap       map[int64]*[]gogpt.ChatCompletionMessage
	baseMsgsGroupMap   map[int64]*[]gogpt.ChatCompletionMessage
	groupTriggerPrefix string
	baseSetting        gogpt.ChatCompletionMessage
	msgsPrivateMap     map[int64]*[]gogpt.ChatCompletionMessage
	baseMsgsPrivateMap map[int64]*[]gogpt.ChatCompletionMessage
	maxContextTokens   int
}

func NewChatBot(trigger string, maxTokenUsage int, openAIToken string, prompt string) *ChatBot {
	cfg := gogpt.DefaultConfig(openAIToken)
	client := gogpt.NewClientWithConfig(cfg)
	baseSetting := gogpt.ChatCompletionMessage{Role: "system", Content: prompt}
	maxContextTokens := maxTokenUsage
	return &ChatBot{
		msgsGroupMap:       make(map[int64]*[]gogpt.ChatCompletionMessage),
		baseMsgsGroupMap:   make(map[int64]*[]gogpt.ChatCompletionMessage),
		groupTriggerPrefix: trigger,
		msgsPrivateMap:     make(map[int64]*[]gogpt.ChatCompletionMessage),
		baseMsgsPrivateMap: make(map[int64]*[]gogpt.ChatCompletionMessage),
		openAI:             client,
		baseSetting:        baseSetting,
		maxContextTokens:   maxContextTokens,
	}
}

func (c *ChatBot) OnMessage(msg model.MessageEvent, action cqhttp.Action) {
	log.Printf("OnMessage: %+v\n", msg)
}

func (c *ChatBot) OnNotice(msg model.NoticeEvent, action cqhttp.Action) {
	log.Printf("OnNotice: %+v\n", msg)
}

func (c *ChatBot) OnRequest(msg model.RequestEvent, action cqhttp.Action) {
	log.Printf("OnRequest: %+v\n", msg)
}

func (c *ChatBot) OnMetaLifecycle(msg model.MetaLifecycleEvent, action cqhttp.Action) {
	log.Printf("OnMetaLifecycle: %+v\n", msg)
}

func (c *ChatBot) OnMessageGroup(msg model.MessageGroupEvent, action cqhttp.Action) {
	prefix := c.groupTriggerPrefix
	if strings.HasPrefix(msg.Message, prefix) {

		msgs, ok := c.msgsGroupMap[msg.GroupID]
		if !ok {
			msgs = &[]gogpt.ChatCompletionMessage{}
			c.msgsGroupMap[msg.GroupID] = msgs
		}

		baseMsgs, ok := c.baseMsgsGroupMap[msg.GroupID]
		if !ok {
			baseMsgs = &[]gogpt.ChatCompletionMessage{}
			c.baseMsgsGroupMap[msg.GroupID] = baseMsgs
		}

		content := strings.TrimSpace(strings.Split(msg.Message, prefix)[1])
		if content == "/help" {
			action.SendGroupMsg(msg.GroupID, "/help: 帮助\n/reset: 重置\n/clear: 清空历史记录\n/set: 设置人设")
		} else if content == "/reset" {
			*baseMsgs = []gogpt.ChatCompletionMessage{
				c.baseSetting,
			}
			*msgs = []gogpt.ChatCompletionMessage{}
			action.SendGroupMsg(msg.GroupID, "/重置成功")
		} else if content == "/clear" {
			*msgs = []gogpt.ChatCompletionMessage{}
			action.SendGroupMsg(msg.GroupID, "/历史记录清空成功")
		} else if strings.HasPrefix(content, "/set") {
			*baseMsgs = []gogpt.ChatCompletionMessage{
				{Role: "system", Content: strings.TrimSpace(strings.Split(content, "/set")[1])},
			}
			*msgs = []gogpt.ChatCompletionMessage{}
			action.SendGroupMsg(msg.GroupID, "/人设已更新")
		} else if strings.HasPrefix(content, "/addset") {
			*baseMsgs = append(*baseMsgs, gogpt.ChatCompletionMessage{Role: "system", Content: strings.TrimSpace(strings.Split(content, "/addset")[1])})
			*msgs = []gogpt.ChatCompletionMessage{}
			action.SendGroupMsg(msg.GroupID, "/人设已追加")
		} else {
			*msgs = append(*msgs, gogpt.ChatCompletionMessage{
				Role:    "user",
				Content: content,
			})
			merge := mergeSlice(*baseMsgs, *msgs)
			response, err := c.openAI.CreateChatCompletion(context.Background(), gogpt.ChatCompletionRequest{
				Model:    gogpt.GPT3Dot5Turbo,
				Messages: merge,
			})
			if err != nil {
				action.SendGroupMsg(msg.GroupID, "/再说一遍")
				return
			}
			output := response.Choices[0].Message
			*msgs = append(*msgs, output)

			if len(*msgs) >= 2 && response.Usage.TotalTokens >= c.maxContextTokens {
				*msgs = (*msgs)[2:]
			}

			// send message
			action.SendGroupMsg(msg.GroupID, output.Content)
		}
	}
}

func (c *ChatBot) OnMessagePrivate(msg model.MessagePrivateEvent, action cqhttp.Action) {
	msgs, ok := c.msgsPrivateMap[msg.UserID]
	if !ok {
		msgs = &[]gogpt.ChatCompletionMessage{}
		c.msgsPrivateMap[msg.UserID] = msgs
	}

	baseMsgs, ok := c.baseMsgsPrivateMap[msg.UserID]
	if !ok {
		baseMsgs = &[]gogpt.ChatCompletionMessage{}
		c.baseMsgsPrivateMap[msg.UserID] = baseMsgs
	}

	content := msg.Message
	if content == "/help" {
		action.SendPrivateMsg(msg.UserID, "/help: 帮助\n/reset: 重置\n/clear: 清空历史记录\n/set: 设置人设")
	} else if content == "/reset" {
		*baseMsgs = []gogpt.ChatCompletionMessage{
			c.baseSetting,
		}
		*msgs = []gogpt.ChatCompletionMessage{}
		action.SendPrivateMsg(msg.UserID, "/重置成功")
	} else if content == "/clear" {
		*msgs = []gogpt.ChatCompletionMessage{}
		action.SendPrivateMsg(msg.UserID, "/历史记录清空成功")
	} else if strings.HasPrefix(content, "/set") {
		*baseMsgs = []gogpt.ChatCompletionMessage{
			{Role: "system", Content: strings.TrimSpace(strings.Split(content, "/set")[1])},
		}
		*msgs = []gogpt.ChatCompletionMessage{}
		action.SendPrivateMsg(msg.UserID, "/人设已更新")
	} else if strings.HasPrefix(content, "/addset") {
		*baseMsgs = append(*baseMsgs, gogpt.ChatCompletionMessage{Role: "system", Content: strings.TrimSpace(strings.Split(content, "/addset")[1])})
		*msgs = []gogpt.ChatCompletionMessage{}
		action.SendPrivateMsg(msg.UserID, "/人设已追加")
	} else {
		*msgs = append(*msgs, gogpt.ChatCompletionMessage{
			Role:    "user",
			Content: content,
		})
		merge := mergeSlice(*baseMsgs, *msgs)
		response, err := c.openAI.CreateChatCompletion(context.Background(), gogpt.ChatCompletionRequest{
			Model:    gogpt.GPT3Dot5Turbo,
			Messages: merge,
		})
		if err != nil {
			action.SendPrivateMsg(msg.UserID, "/再说一遍")
			return
		}
		output := response.Choices[0].Message
		*msgs = append(*msgs, output)

		if len(*msgs) >= 2 && response.Usage.TotalTokens >= c.maxContextTokens {
			*msgs = (*msgs)[2:]
		}

		// send message
		action.SendPrivateMsg(msg.UserID, output.Content)
	}
}

func mergeSlice(a, b []gogpt.ChatCompletionMessage) []gogpt.ChatCompletionMessage {
	c := make([]gogpt.ChatCompletionMessage, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}
