package main

import (
	"ergodownloader/internal/configuration"
	"ergodownloader/internal/configuration/clientConfig"
	"ergodownloader/internal/redditApi"
	"flag"
	"fmt"
	"log"
	"os"
)

var Version = "v1.0"
var BuildTime = "Not set"

func main() {
	unsave, showVersion, noLogFile := parseFlags()
	if showVersion {
		fmt.Printf("ErGo Downloader %s\nBuild Time: %s\n", Version, BuildTime)
		os.Exit(0)
	}
	fmt.Println("################")
	fmt.Println("ErGo Downloader")
	fmt.Println("################")
	fmt.Println()
	configuration.Load()
	token, err := redditApi.GetToken()
	if err != nil {
		log.Println(err)
		fmt.Println("Unable to get the login token. Please delete configuration file and try again.")
		return
	}
	var firstId string
	var afterId string
	for {
		posts, err := redditApi.GetSavedPosts(token, clientConfig.Credentials.Username, afterId)
		if err != nil {
			fmt.Println(err)
			return
		}
		var firstPost string
		firstPost, afterId, err = redditApi.DownloadSavedPosts(posts, token, unsave, !noLogFile)
		if firstId == "" {
			firstId = firstPost
		}
		if afterId == "" || err != nil {
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}
	if !unsave {
		configuration.SetMaxPost(firstId)
	}
}

func parseFlags() (bool, bool, bool) {
	isUnsave := flag.Bool("unsave", false, "Pass to unsave downloaded posts")
	noLog := flag.Bool("no-log", false, "Disable log file")
	showVersion := flag.Bool("v", false, "Show version")
	flag.Parse()
	return *isUnsave, *showVersion, *noLog
}
