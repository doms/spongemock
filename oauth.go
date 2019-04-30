package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// OauthResponse - the response body from https://slack.com/api/oauth.access
type OauthResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

// AuthResponse - the response body from https://slack.com/api/auth.test
type AuthResponse struct {
	OK    bool   `json:"ok"`
	URL   string `json:"url"`
	Error string `json:"error"`
}

// Auth - handles the request and send back mocking response
func Auth(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params, err := url.ParseQuery(string(body))
	if err != nil {
		log.Fatalln(err)
	}

	// Extract code from request url
	code := params.Get("code")

	// Compose authHeader by encoding the string ${client_id}:${client_secret}
	clientID := os.Getenv("SLACK_CLIENT_ID")
	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")

	data := fmt.Sprintf("%s:%s", clientID, clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	authorization := fmt.Sprintf("Basic %s", encoded)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// redirect to slack workspace url
	if link, err := buildSlackURL(client, code, authorization); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, link, http.StatusFound)
	}
}

func buildSlackURL(client *http.Client, code, authorization string) (string, error) {
	// Hit oauth.access for access_token
	request, err := http.NewRequest(http.MethodPost, "https://slack.com/api/oauth.access", strings.NewReader(code))
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Authorization", authorization)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	oauthResponse := OauthResponse{}
	err = json.Unmarshal(body, &oauthResponse)
	if err != nil {
		log.Fatal(err)
	}

	accessToken := oauthResponse.AccessToken

	params := url.Values{}
	params.Set("token", accessToken)

	// Hit auth.text for slack domain
	request, err = http.NewRequest(http.MethodPost, "https://slack.com/api/auth.test", strings.NewReader(params.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err = client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	authResponse := AuthResponse{}
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		log.Fatalln(err)
	}

	if authResponse.OK {
		return authResponse.URL, nil
	}

	return "", errors.New(authResponse.Error)
}
