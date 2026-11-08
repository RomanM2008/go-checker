package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	client := &http.Client{Timeout: 10 * time.Second}

	tgBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	tgChatId := os.Getenv("TELEGRAM_CHAT_ID")

	if tgBotToken == "" || tgChatId == "" {
		fmt.Println("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID environment variables must be set")
		os.Exit(1)
	}

	notifier := NewNotifier(tgBotToken, tgChatId)

	monitors := []Monitor{
		{URL: "https://google.com", ExpectedStatus: 200, Interval: 10 * time.Second},
		{URL: "http://localhost:8080", ExpectedStatus: 200, Interval: 1 * time.Second},
	}

	for _, m := range monitors {
		go runMonitor(client, notifier, m)
	}

	select {}
}
