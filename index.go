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

	content := SlackRequest{}
	content.Text = params.Get("text")
	content.ResponseURL = params.Get("response_url")

	splitContent := strings.Split(content.Text, " ")

	var t []string
	for _, word := range splitContent {
		t = append(t, SpongeMock(word))
	}

	mockifiedText := strings.Join(t, " ")

	payload := SlackResponse{
		ResponseType: "in_channel",
		Text:         mockifiedText,
	}

	sendSlackNotification(content.ResponseURL, &payload)
}

// SpongeMock - converts strings into sponges
func SpongeMock(content string) string {
	if content == "" {
		return ""
	}

	// ignore mentions (users, channels)
	if strings.Index("<@#", string(content[0])) != -1 {
		return content
	}

	buf := []rune(strings.ToLower(content))

	for i, char := range buf {
		if buf[i] > unicode.MaxASCII {
			continue
		}

		if i%2 == 0 {
			if string(char) != "i" {
				buf[i] -= 32
			}
		} else {
			if string(char) == "l" {
				buf[i] -= 32
			}
		}
	}

	return string(buf)
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
