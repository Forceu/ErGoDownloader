package configuration

import (
	"encoding/json"
	"ergodownloader/internal/configuration/clientConfig"
	"ergodownloader/internal/configuration/defaultConfig"
	"ergodownloader/internal/helper"
	"ergodownloader/internal/models"
	"ergodownloader/internal/redditApi/redditAuth"
	"fmt"
	"os"
)

const configDir = "config"
const configPath = configDir + "/config.json"

// Load loads a saved config and creates it on first start
func Load() {
	helper.CreateDir(configDir)
	helper.CreateDir("download")
	if !helper.FileExists(configPath) {
		clientConfig.Credentials = defaultConfig.Generate()
		save()
	}
	file, err := os.Open(configPath)
	helper.Check(err)
	defer file.Close()
	decoder := json.NewDecoder(file)
	credentials := models.Configuration{}
	err = decoder.Decode(&credentials)
	clientConfig.Credentials = credentials
	helper.Check(err)
	redditAuth.Init()
}

// SetMaxPost sets post from which on no further posts will be considered
func SetMaxPost(id string) {
	clientConfig.Credentials.MaxPost = id
	save()
}

// Save the configuration as a json file
func save() {
	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error reading configuration:", err)
		os.Exit(1)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(&clientConfig.Credentials)
	if err != nil {
		fmt.Println("Error writing configuration:", err)
		os.Exit(1)
	}
}
