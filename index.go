package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
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

// SlackResponse - structure of a resposne sent to Slack
type SlackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// Handler - handles the request and send back mocking response
func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

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
		if strings.Index("<@#", string(word[0])) != -1 {
			t = append(t, word)
		} else {
			t = append(t, spongeMock(word))
		}
	}

	mockifiedText := strings.Join(t, " ")

	payload := SlackResponse{
		ResponseType: "in_channel",
		Text:         mockifiedText,
	}

	sendSlackNotification(content.ResponseURL, &payload)
}

func spongeMock(content string) string {
	originalContent := strings.ToLower(content)
	var res []string

	shouldUpcase := true

	for _, char := range originalContent {
		if char >= 97 && char < 123 {
			if shouldUpcase {
				if string(char) == "i" {
					res = append(res, string(char))
				} else {
					res = append(res, string(char-rune(32)))
				}

				shouldUpcase = false
			} else {
				if string(char) == "l" {
					res = append(res, string(char-rune(32)))
				} else {
					res = append(res, string(char))
				}

				shouldUpcase = true
			}
		} else {
			res = append(res, string(char))
		}
	}

	return strings.Join(res, "")
}

// sendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func sendSlackNotification(webhookURL string, payload *SlackResponse) error {
	slackBody, _ := json.Marshal(payload)
	fmt.Println(string(slackBody))
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
