package main

import (
	"OpenAIBot/src/controller"
	"OpenAIBot/src/module/cqhttp"
	"OpenAIBot/src/module/cqhttp/model"
	"log"
	"os"
	"os/signal"
)

func main() {

	chatBot := controller.NewChatBot()
	client := cqhttp.NewCQHttp(
		&model.CQHttpConfig{
			Host:  "localhost:5701",
			Path:  "/",
			Token: "bot",
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
