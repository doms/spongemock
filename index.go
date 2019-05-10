package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
)

// SlackRequest - structure of a request sent from Slack
type SlackRequest struct {
	Token       string `json:"token,omitempty"`
	Command     string `json:"command,omitempty"`
	Text        string `json:"text,omitempty"`
	ResponseURL string `json:"response_url,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	UserName    string `json:"user_name,omitempty"`
	TeamID      string `json:"team_id,omitempty"`
	ChannelID   string `json:"channel_id,omitempty"`
}

// SlackResponse - structure of a response sent to Slack
type SlackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// Handler - handles the request and sends back mocking response
func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	params, err := url.ParseQuery(string(body))
	if err != nil {
		log.Fatalln(err)
	}

	content := SlackRequest{
		Text:        params.Get("text"),
		ResponseURL: params.Get("response_url"),
	}

	sendSlackNotification(
		content.ResponseURL,
		&SlackResponse{
			ResponseType: "in_channel",
			Text:         SpongeMock(content.Text),
		},
	)
}

// SpongeMock - converts strings into sponges
func SpongeMock(sentence string) string {
	words := strings.Split(sentence, " ")
	for i, word := range words {
		if word == "" {
			continue
		}

		// ignore mentions (users, channels)
		if word[0] == '#' || word[0] == '<' || word[0] == '@' {
			continue
		}

		buf := []rune(strings.ToLower(word))
		for i, char := range buf {
			if buf[i] > unicode.MaxASCII {
				continue
			}

			if i%2 == 0 {
				if char != 'i' {
					buf[i] -= 32
				}
			} else {
				if char == 'l' {
					buf[i] -= 32
				}
			}
		}

		words[i] = string(buf)
	}

	return strings.Join(words, " ")
}

// sendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func sendSlackNotification(webhookURL string, payload *SlackResponse) error {
	slackBody, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
