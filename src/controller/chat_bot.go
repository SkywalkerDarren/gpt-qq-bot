package controller

import (
	"OpenAIBot/src/module/cqhttp"
	"OpenAIBot/src/module/cqhttp/model"
	"context"
	"os"
	"strings"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type ChatBot struct {
	cqhttp.DefaultHandler
	openAI             *gogpt.Client
	msgsGroupMap       map[int64][]gogpt.ChatCompletionMessage
	baseMsgsGroupMap   map[int64][]gogpt.ChatCompletionMessage
	baseSetting        gogpt.ChatCompletionMessage
	msgsPrivateMap     map[int64][]gogpt.ChatCompletionMessage
	baseMsgsPrivateMap map[int64][]gogpt.ChatCompletionMessage
	maxContextTokens   int
}

func NewChatBot() *ChatBot {
	tk := os.Getenv("OPENAI_KEY")
	cfg := gogpt.DefaultConfig(tk)
	client := gogpt.NewClientWithConfig(cfg)
	baseSetting := gogpt.ChatCompletionMessage{Role: "system", Content: "你是一个说中文的AI聊天机器人，回答尽可能简短，名字叫做\"聊天狗屁通\""}
	maxContextTokens := 500
	return &ChatBot{
		msgsGroupMap:       make(map[int64][]gogpt.ChatCompletionMessage),
		baseMsgsGroupMap:   make(map[int64][]gogpt.ChatCompletionMessage),
		msgsPrivateMap:     make(map[int64][]gogpt.ChatCompletionMessage),
		baseMsgsPrivateMap: make(map[int64][]gogpt.ChatCompletionMessage),
		openAI:             client,
		baseSetting:        baseSetting,
		maxContextTokens:   maxContextTokens,
	}
}

func (c *ChatBot) OnMessageGroup(msg model.MessageGroupEvent, action cqhttp.Action) {
	prefix := "[CQ:at,qq=2869015936]"
	if strings.HasPrefix(msg.Message, prefix) {

		msgs, ok := c.msgsGroupMap[msg.GroupID]
		if !ok {
			msgs = []gogpt.ChatCompletionMessage{}
			c.msgsGroupMap[msg.GroupID] = msgs
		}

		baseMsgs, ok := c.baseMsgsGroupMap[msg.GroupID]
		if !ok {
			baseMsgs = []gogpt.ChatCompletionMessage{}
			c.baseMsgsGroupMap[msg.GroupID] = baseMsgs
		}

		content := strings.TrimSpace(strings.Split(msg.Message, prefix)[1])
		if content == "/help" {
			action.SendGroupMsg(msg.GroupID, "/clear 清空历史记录\n/set 设置人设\n/help 帮助\n其他就是聊天")
		} else if content == "/clear" {
			c.msgsGroupMap[msg.GroupID] = []gogpt.ChatCompletionMessage{}
			action.SendGroupMsg(msg.GroupID, "/历史记录清空成功")
		} else if strings.HasPrefix(content, "/set") {
			baseMsgs = []gogpt.ChatCompletionMessage{
				c.baseSetting,
				{Role: "system", Content: strings.TrimSpace(strings.Split(content, "/set")[1])},
			}
			c.baseMsgsGroupMap[msg.GroupID] = baseMsgs
			c.msgsGroupMap[msg.GroupID] = []gogpt.ChatCompletionMessage{}
			action.SendGroupMsg(msg.GroupID, "/人设已更新")
		} else {
			msgs = append(msgs, gogpt.ChatCompletionMessage{
				Role:    "user",
				Content: content,
			})
			merge := mergeSlice(baseMsgs, msgs)
			response, err := c.openAI.CreateChatCompletion(context.Background(), gogpt.ChatCompletionRequest{
				Model:    gogpt.GPT3Dot5Turbo,
				Messages: merge,
			})
			if err != nil {
				action.SendGroupMsg(msg.GroupID, "/再说一遍")
				return
			}
			output := response.Choices[0].Message
			msgs = append(msgs, output)

			if len(msgs) >= 2 && response.Usage.TotalTokens >= c.maxContextTokens {
				msgs = msgs[2:]
			}
			c.msgsGroupMap[msg.GroupID] = msgs

			// send message
			action.SendGroupMsg(msg.GroupID, output.Content)
		}
	}
}

func (c *ChatBot) OnMessagePrivate(msg model.MessagePrivateEvent, action cqhttp.Action) {
	msgs, ok := c.msgsPrivateMap[msg.UserID]
	if !ok {
		msgs = []gogpt.ChatCompletionMessage{}
		c.msgsPrivateMap[msg.UserID] = msgs
	}

	baseMsgs, ok := c.baseMsgsPrivateMap[msg.UserID]
	if !ok {
		baseMsgs = []gogpt.ChatCompletionMessage{}
		c.baseMsgsPrivateMap[msg.UserID] = baseMsgs
	}

	content := msg.Message
	if content == "/help" {
		action.SendPrivateMsg(msg.UserID, "/clear 清空历史记录\n/set 设置人设\n/help 帮助\n其他就是聊天")
	} else if content == "/clear" {
		c.msgsPrivateMap[msg.UserID] = []gogpt.ChatCompletionMessage{}
		action.SendPrivateMsg(msg.UserID, "/历史记录清空成功")
	} else if strings.HasPrefix(content, "/set") {
		baseMsgs = []gogpt.ChatCompletionMessage{
			c.baseSetting,
			{Role: "system", Content: strings.TrimSpace(strings.Split(content, "/set")[1])},
		}
		c.baseMsgsPrivateMap[msg.UserID] = baseMsgs
		c.msgsPrivateMap[msg.UserID] = []gogpt.ChatCompletionMessage{}
		action.SendPrivateMsg(msg.UserID, "/人设已更新")
	} else {
		msgs = append(msgs, gogpt.ChatCompletionMessage{
			Role:    "user",
			Content: content,
		})
		merge := mergeSlice(baseMsgs, msgs)
		response, err := c.openAI.CreateChatCompletion(context.Background(), gogpt.ChatCompletionRequest{
			Model:    gogpt.GPT3Dot5Turbo,
			Messages: merge,
		})
		if err != nil {
			action.SendPrivateMsg(msg.UserID, "/再说一遍")
			return
		}
		output := response.Choices[0].Message
		msgs = append(msgs, output)

		if len(msgs) >= 2 && response.Usage.TotalTokens >= c.maxContextTokens {
			msgs = msgs[2:]
		}
		c.msgsPrivateMap[msg.UserID] = msgs

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
