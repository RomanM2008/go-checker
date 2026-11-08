package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Elapsed    time.Duration
	Up         bool
	Error      error
}

func checkStatus(client *http.Client, url string, expectedStatus int) Result {
	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)

	result := Result{
		URL:     url,
		Elapsed: elapsed,
	}

	if err != nil {
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Up = resp.StatusCode == expectedStatus
	return result
}

func printResult(result Result) {
	if result.Error != nil {
		fmt.Println(result.URL, "DOWN", result.Error)
		return
	}
	if result.Up {
		fmt.Println(result.URL, "UP", result.StatusCode, result.Elapsed)
	} else {
		fmt.Println(result.URL, "DOWN", result.StatusCode, result.Elapsed)
	}
}

func monitorURL(client *http.Client, n *Notifier, url string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	previousResult := checkStatus(client, url, 200)
	printResult(previousResult)
	notify(n, previousResult)

	for range ticker.C {
		result := checkStatus(client, url, 200)
		if result.Up != previousResult.Up {
			printResult(result)
			notify(n, result)
			previousResult = result
		}
	}
}

func notify(n *Notifier, result Result) {
	if err := n.Send(formatMessage(result)); err != nil {
		log.Println("failed to send notification:", err)
	}
}

func formatMessage(result Result) string {
	if result.Up {
		return fmt.Sprintf("🟢 %s is UP again", result.URL)
	}
	if result.Error != nil {
		return fmt.Sprintf("🔴 %s is DOWN: %v", result.URL, result.Error)
	}
	return fmt.Sprintf("🔴 %s is DOWN: unexpected status %d", result.URL, result.StatusCode)
}

func main() {
	urls := []string{
		"http://localhost:8080",
	}
	var client = &http.Client{Timeout: 10 * time.Second}

	tg_bot_token := os.Getenv("TELEGRAM_BOT_TOKEN")
	tg_chat_id := os.Getenv("TELEGRAM_CHAT_ID")

	notifier := NewNotifier(tg_bot_token, tg_chat_id)

	if tg_bot_token == "" {
		fmt.Println("TELEGRAM_BOT_TOKEN environment variable is not set")
		os.Exit(1)
	}
	if tg_chat_id == "" {
		fmt.Println("TELEGRAM_CHAT_ID environment variable is not set")
		os.Exit(1)
	}

	for _, url := range urls {
		go func() {
			monitorURL(client, notifier, url, 10*time.Second)
		}()
	}

	select {}
}
