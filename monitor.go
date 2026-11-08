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

type StatusChange struct {
	URL            string
	ExpectedStatus int
	Up             bool
	StatusCode     int
	Error          string
	DownFor        time.Duration
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

func newStatusChange(m Monitor, r Result) StatusChange {
	st := StatusChange{
		URL:            m.URL,
		ExpectedStatus: m.ExpectedStatus,
		Up:             r.Up,
		StatusCode:     r.StatusCode,
	}

	if r.Error != nil {
		st.Error = r.Error.Error()
	}
	return st
}

func runMonitor(client *http.Client, n *Notifier, m Monitor) {
	var downSince time.Time
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	wasUp := true

	for {
		result := checkStatus(client, m.URL, m.ExpectedStatus)
		printResult(result)
		if result.Up != wasUp {
			st := newStatusChange(m, result)

			if result.Up {
				st.DownFor = time.Since(downSince)
			} else {
				downSince = time.Now()
			}

			notify(n, st)
			wasUp = result.Up
		}
		<-ticker.C
	}
}

func notify(n *Notifier, st StatusChange) {
	if err := n.Send(formatMessage(st)); err != nil {
		log.Println("failed to send notification:", err)
	}
}

func formatMessage(st StatusChange) string {
	if st.Up {
		if st.DownFor > 0 {
			return fmt.Sprintf("🟢 %s is UP again. Down for %v", st.URL, st.DownFor.Round(time.Second))
		}
		return fmt.Sprintf("🟢 %s is UP", st.URL)
	}
	if st.Error != "" {
		return fmt.Sprintf("🔴 %s is DOWN: %s", st.URL, st.Error)
	}
	return fmt.Sprintf("🔴 %s is DOWN: unexpected status %d, expected %d", st.URL, st.StatusCode, st.ExpectedStatus)
}
