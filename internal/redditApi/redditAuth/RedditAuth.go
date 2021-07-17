package redditAuth

import (
	"encoding/json"
	"ergodownloader/internal/models"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	Application models.Application
)

// Init loads the credentials for the authentication
func Init() {
	Application = application
}

// GetOAuthUrl returns an url for authentication by the user
func GetOAuthUrl() (string, string) {
	state := strconv.Itoa(rand.Int())
	params := url.Values{}
	params.Add("client_id", Application.Id)
	params.Add("response_type", "code")
	params.Add("state", state)
	params.Add("redirect_uri", Application.RedirectUrl)
	params.Add("scope", "history,identity,save")
	params.Add("duration", "permanent")
	return "https://www.reddit.com/api/v1/authorize?" + params.Encode(), state
}

// GetInitialTokens requests a RefreshToken and an AccessToken with an auth code
func GetInitialTokens(authCode string) (string, string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("redirect_uri", Application.RedirectUrl)

	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("User-agent", Application.Useragent)

	req.SetBasicAuth(Application.Id, "")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var responseMap map[string]interface{}

	err = json.Unmarshal(response, &responseMap)
	if err != nil {
		return "", "", err
	}

	refreshToken := responseMap["refresh_token"]
	if refreshToken == nil || refreshToken.(string) == "" {
		return "", "", errors.New("invalid token")
	}

	accessToken := responseMap["access_token"]
	if accessToken == nil || accessToken.(string) == "" {
		return "", "", errors.New("invalid token")
	}

	return refreshToken.(string), accessToken.(string), nil
}

// GetUserName returns the username of the authorised user
func GetUserName(token string) (string, error) {

	req, err := http.NewRequest("GET", "https://oauth.reddit.com/api/v1/me", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("User-agent", Application.Useragent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var arbitraryJson map[string]interface{}

	err = json.Unmarshal(response, &arbitraryJson)
	if err != nil {
		return "", err
	}

	if arbitraryJson["name"] == nil {
		return "", errors.New("invalid response")
	}

	return arbitraryJson["name"].(string), nil
}
