package main

import (
	"OpenAIBot/src/model"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func main() {
	tk := os.Getenv("OPENAI_KEY")
	cfg := gogpt.DefaultConfig(tk)
	client := gogpt.NewClientWithConfig(cfg)
	ctx := context.Background()

	baseMsgs := []gogpt.ChatCompletionMessage{
		{"system", "你是一个说中文的AI助手，名字叫做\"聊天狗屁通\""},
		// {"user", "Translate Text to English: \"你好\""},
		// {"assistant", "Hello\nHi"},
		// {"user", "Translate Text to English: \"我是中国人\""},
		// {"assistant", "I am Chinese"},
		// {"system", "你是个小说家。"},
	}

	var msgs []gogpt.ChatCompletionMessage

	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:5701", Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": {"Bearer bot"}})
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("recv: %s", message)
			var event model.Event
			_ = json.Unmarshal(message, &event)
			if event.PostType == model.PostTypeMessage {
				// do something with message
				var msg model.MessageEvent
				_ = json.Unmarshal(message, &msg)
				if msg.MessageType == model.MessageTypeGroup {
					// do something with group message
					var groupMsg model.MessageGroupEvent
					_ = json.Unmarshal(message, &groupMsg)
					prefix := "[CQ:at,qq=2869015936]"
					if strings.HasPrefix(groupMsg.Message, prefix) {
						msgs = append(msgs, gogpt.ChatCompletionMessage{
							Role:    "user",
							Content: strings.TrimSpace(strings.Split(groupMsg.Message, prefix)[1]),
						})
						merge := mergeSlice(baseMsgs, msgs)
						response, err := client.CreateChatCompletion(ctx, gogpt.ChatCompletionRequest{
							Model:    gogpt.GPT3Dot5Turbo,
							Messages: merge,
						})
						if err != nil {
							return
						}
						output := response.Choices[0].Message
						msgs = append(msgs, output)

						if len(msgs) >= 2 && response.Usage.TotalTokens >= 500 {
							msgs = msgs[2:]
						}

						// send message
						sendMsg := model.WebsocketActionRequest{
							Action: model.SendGroupMsg,
							Params: model.SendGroupMsgParams{
								GroupID: groupMsg.GroupID,
								Message: output.Content,
							},
							Echo: "echo",
						}
						encoded, _ := json.Marshal(sendMsg)
						err = c.WriteMessage(websocket.TextMessage, encoded)
						if err != nil {
							return
						}
						log.Println("send:", string(encoded))
					}

				} else if msg.MessageType == model.MessageTypePrivate {
					// do something with private message
					var privateMsg model.MessagePrivateEvent
					_ = json.Unmarshal(message, &privateMsg)
				}
			} else if event.PostType == model.PostTypeNotice {
				// do something with notice
			} else if event.PostType == model.PostTypeRequest {
				// do something with request
			} else if event.PostType == model.PostTypeMetaEvent {
				// do something with meta event
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func mergeSlice(a, b []gogpt.ChatCompletionMessage) []gogpt.ChatCompletionMessage {
	c := make([]gogpt.ChatCompletionMessage, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}
