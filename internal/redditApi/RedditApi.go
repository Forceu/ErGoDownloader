package redditApi

import (
	"encoding/json"
	"ergodownloader/internal/configuration/clientConfig"
	"ergodownloader/internal/helper"
	"ergodownloader/internal/redditApi/redditAuth"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetToken gets an accessToken with a refresh token
func GetToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", clientConfig.Credentials.Token)

	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("User-agent", redditAuth.Application.Useragent)

	req.SetBasicAuth(redditAuth.Application.Id, "")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseMap map[string]interface{}

	err = json.Unmarshal(response, &responseMap)
	if err != nil {
		return "", err
	}

	token := responseMap["access_token"]
	if token == nil || token.(string) == "" {
		return "", errors.New("invalid login")
	}

	return token.(string), nil
}

// GetSavedPosts gets the last 100 posts and downloads them
func GetSavedPosts(token, username, afterId string) ([]interface{}, error) {

	var arbitraryJson map[string]interface{}

	req, err := http.NewRequest("GET", "https://oauth.reddit.com/user/"+username+"/saved?limit=100&after="+afterId, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("User-agent", redditAuth.Application.Useragent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &arbitraryJson)
	if err != nil {
		return nil, err
	}

	data := arbitraryJson["data"].(map[string]interface{})
	children := data["children"].([]interface{})

	return children, nil
}

// DownloadSavedPosts downloads the given posts and unsaves them if requested
func DownloadSavedPosts(posts []interface{}, token string, unsave, logging bool) (string, string, error) {
	var firstPost, lastPost string
	for _, value := range posts {
		child := value.(map[string]interface{})
		childData := child["data"].(map[string]interface{})
		if firstPost == "" {
			firstPost = childData["name"].(string)
		}
		lastPost = childData["name"].(string)

		if lastPost == clientConfig.Credentials.MaxPost {
			return firstPost, "", nil
		}

		if childData["title"] == nil {
			continue
		}
		if childData["url"] == nil {
			continue
		}
		isValidLink, extension := helper.IsDirectLink(childData["url"].(string))
		if isValidLink {
			duplicate, err := helper.DownloadFile(childData["title"].(string)+extension, childData["url"].(string))
			if err != nil {
				return "", "", err
			}
			if duplicate {
				fmt.Println("Duplicate, not saving: " + childData["title"].(string))
			} else {
				fmt.Println("Download complete: " + childData["title"].(string))
				if logging {
					helper.WriteLog(childData["title"].(string) + ":  https://reddit.com" + childData["permalink"].(string))
				}
			}
			if unsave {
				err = UnsavePost(token, childData["name"].(string))
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println("Unable to download: " + childData["url"].(string))
		}
	}
	if len(posts) < 100 {
		return firstPost, "", nil
	} else {
		return firstPost, lastPost, nil
	}
}

// UnsavePost removes a submission from the saved posts
func UnsavePost(token, id string) error {
	data := url.Values{}
	data.Set("id", id)

	req, err := http.NewRequest("POST", "https://oauth.reddit.com/api/unsave", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("User-agent", redditAuth.Application.Useragent)
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
