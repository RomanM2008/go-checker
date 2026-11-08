package main

import (
	"fmt"
	"log"
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

type Monitor struct {
	URL            string
	ExpectedStatus int
	Interval       time.Duration
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

func runMonitor(client *http.Client, n *Notifier, m Monitor) {
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	previousResult := checkStatus(client, m.URL, m.ExpectedStatus)
	printResult(previousResult)
	if !previousResult.Up {
		notify(n, previousResult)
	}

	for range ticker.C {
		result := checkStatus(client, m.URL, m.ExpectedStatus)
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
