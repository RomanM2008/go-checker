package main

import (
	"fmt"
	"net/http"
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

func PrintResult(result Result) {
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

func monitorURL(client *http.Client, url string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		result := checkStatus(client, url, 200)
		PrintResult(result)
		<-ticker.C
	}
}

func main() {
	urls := []string{
		"https://google.com",
		"https://gobyexample.com",
		"https://mail.google.com",
	}
	var client = &http.Client{Timeout: 10 * time.Second}

	for _, url := range urls {
		go func(u string) {
			monitorURL(client, u, 10*time.Second)
		}(url)
	}

	select {}
}
