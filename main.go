package main

import (
	"fmt"
	"net/http"
	"sync"
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

func main() {
	urls := []string{
		"https://google.com",
		"https://gobyexample.com",
		"https://mail.google.com",
	}
	var wg sync.WaitGroup
	var client = &http.Client{Timeout: 10 * time.Second}

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := checkStatus(client, u, 200)
			if result.Error != nil {
				fmt.Println(result.URL, "DOWN", result.Error)
				return
			}
			if result.Up {
				fmt.Println(result.URL, "UP", result.StatusCode, result.Elapsed)
			} else {
				fmt.Println(result.URL, "DOWN", result.StatusCode, result.Elapsed)
			}
		}(url)
	}

	wg.Wait()
}
