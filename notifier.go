package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Notifier struct {
	client *http.Client
	token  string
	chatID string
}

func NewNotifier(token, chatID string) *Notifier {
	return &Notifier{
		client: &http.Client{Timeout: 10 * time.Second},
		token:  token,
		chatID: chatID,
	}
}

func (n *Notifier) Send(text string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.token)

	data := url.Values{}
	data.Set("chat_id", n.chatID)
	data.Set("text", text)

	resp, err := n.client.PostForm(endpoint, data)
	if err != nil {
		return errors.New("telegram request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram request failed with status: %s", resp.Status)
	}

	return nil
}
