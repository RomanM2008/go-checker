package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func checkStatus(client *http.Client, url string) (int, time.Duration, error) {
	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		return 0, elapsed, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, elapsed, nil
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
			code, elapsed, err := checkStatus(client, u)
			if err != nil {
				fmt.Println(u, "ERROR:", err)
				return
			}
			fmt.Println(u, code, elapsed)
		}(url)
	}

	wg.Wait()
}
