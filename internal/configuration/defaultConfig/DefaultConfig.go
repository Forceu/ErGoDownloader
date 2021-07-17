package defaultConfig

import (
	"ergodownloader/internal/models"
	"ergodownloader/internal/redditApi"
	"ergodownloader/internal/redditApi/redditAuth"
	"ergodownloader/internal/redditApi/redditAuth/webserver"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/browser"
	"log"
	"os"
)

// Generate creates the initial config and user authorisation
func Generate() models.Configuration {
	redditAuth.Init()
	url, state := redditAuth.GetOAuthUrl()
	codeInput := make(chan string)
	go webserver.Start(state, codeInput)
	fmt.Println("ErGo Downloader needs to authenticate with Reddit. A browser window will be opened for authentication.")
	err := browser.OpenURL(url)
	if err != nil {
		fmt.Println("Please open " + url + " to authenticate")
	}

	result := <-codeInput
	webserver.Stop()
	if result == "error" {
		fmt.Println("Error receiving authorisation. Please try again.")
		os.Exit(1)
	}
	refreshToken, accessToken, err := redditAuth.GetInitialTokens(result)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	username, err := redditAuth.GetUserName(accessToken)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Logged in as " + username)

	maxId := ShowSelectLastPost(accessToken, username)

	return models.Configuration{
		Token:    refreshToken,
		Username: username,
		MaxPost:  maxId,
	}
}

// ShowSelectLastPost shows a dialog to select the post until which other posts are considered
func ShowSelectLastPost(token, username string) string {
	fmt.Println("Gathering all saved posts, this might take a while...")
	options := []string{"None"}
	optionsId := []string{""}
	var afterId string
	for {
		posts, err := redditApi.GetSavedPosts(token, username, afterId)
		if err != nil {
			log.Fatal(err)
		}
		for _, post := range posts {
			child := post.(map[string]interface{})
			childData := child["data"].(map[string]interface{})
			if childData["title"] != nil && childData["name"] != nil {
				options = append(options, childData["title"].(string)+" ("+childData["name"].(string)+")")
				optionsId = append(optionsId, childData["name"].(string))
			}
			afterId = childData["name"].(string)
		}
		if len(posts) < 100 {
			break
		}
	}
	var result int
	prompt := &survey.Select{
		Message: "Please choose the first post, after which results are not considered. Press Enter to consider all posts:",
		Options: options,
		Default: 0,
	}
	survey.AskOne(prompt, &result, nil)
	if result < 1 {
		return ""
	}

	return optionsId[result]
}
