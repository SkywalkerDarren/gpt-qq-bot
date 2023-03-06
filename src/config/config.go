package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	// 机器人配置
	TriggerPrefix  string `json:"trigger_prefix"`
	DefaultPrompt  string `json:"default_prompt"`
	WebsocketHost  string `json:"websocket_host"`
	WebsocketToken string `json:"websocket_token"`
	OpenAIKey      string `json:"openai_key"`
	MaxTokens      int    `json:"max_tokens"`
}

func GetConfigFromJsonFile() *Config {
	var config Config

	// check if config.json exists
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		log.Println("config.json 不存在，正在创建...")
		openFile, err := os.OpenFile("config.json", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer openFile.Close()
		defaultConfig := Config{
			TriggerPrefix:  "/chat",
			DefaultPrompt:  "你是个说中文的机器人",
			WebsocketHost:  "localhost:6700",
			WebsocketToken: "cqhttp配置的token",
			OpenAIKey:      "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			MaxTokens:      1024,
		}
		data, err := json.MarshalIndent(defaultConfig, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		_, err = openFile.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("请在 config.json 中填写配置信息后重启程序")
		os.Exit(0)
	}

	openFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer openFile.Close()
	data, err := io.ReadAll(openFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
