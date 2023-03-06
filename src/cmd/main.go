package main

import (
	"OpenAIBot/src/config"
	"OpenAIBot/src/controller"
	"OpenAIBot/src/module/cqhttp"
	"OpenAIBot/src/module/cqhttp/model"
	"log"
	"os"
	"os/signal"
)

func main() {

	cfg := config.GetConfigFromJsonFile()

	chatBot := controller.NewChatBot(
		cfg.TriggerPrefix,
		cfg.MaxTokens,
		cfg.OpenAIKey,
		cfg.DefaultPrompt,
	)
	client := cqhttp.NewCQHttp(
		&model.CQHttpConfig{
			Host:  cfg.WebsocketHost,
			Path:  "/",
			Token: cfg.WebsocketToken,
		},
		chatBot,
	)

	err := client.Run()
	if err != nil {
		log.Println(err)
		return
	}

	defer client.Close()

	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			return
		}
	}
}
