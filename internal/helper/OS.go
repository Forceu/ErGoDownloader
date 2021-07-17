package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// CreateDir creates the data folder if it does not exist
func CreateDir(name string) {
	if !FolderExists(name) {
		err := os.Mkdir(name, 0770)
		Check(err)
	}
}

// FolderExists returns true if a folder exists
func FolderExists(folder string) bool {
	_, err := os.Stat(folder)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// Check panics if error is not nil
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// FileExists returns true if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsDirectLink returns true and the extension if PNG or JPG
func IsDirectLink(url string) (bool, string) {
	if strings.HasSuffix(url, ".png") {
		return true, ".png"
	}
	if strings.HasSuffix(url, ".jpg") {
		return true, ".jpg"
	}
	if strings.HasSuffix(url, ".jpeg") {
		return true, ".jpeg"
	}
	return false, ""
}

// WriteLog appends a string to the logfile
func WriteLog(text string) {
	file, err := os.OpenFile("download/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(time.Now().UTC().Format(time.RFC1123) + "   " + text + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// DownloadFile downloads a file to download/+name
func DownloadFile(filename string, url string) error {
	name := getFileName(filename, 0)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create("download/" + name)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getFileName(name string, counter int) string {
	cleanName := sanitizeFileName(name)
	ext := filepath.Ext(cleanName)
	if len(cleanName) > 130 {
		cleanName = cleanName[:130] + ext
	}
	noExt := cleanName[0 : len(cleanName)-len(ext)]
	newName := cleanName
	if counter != 0 {
		newName = noExt + "_" + strconv.Itoa(counter) + ext
	}
	if FileExists("download/" + newName) {
		counter++
		return getFileName(name, counter)
	} else {
		return newName
	}
}

func sanitizeFileName(name string) string {
	var invalidChars = []string{
		"/",
		"<",
		">",
		"\"",
		"\\",
		"&",
		"$",
		"#",
		"{", "}", "[", "]", "=",
		";", "?", "%20", "%22",
		"%3c",   // <
		"%253c", // <
		"%3e",   // >
		"%0e",   //
		"%26",   // &
		"%24",   // $
		"%3f",   // ?
		"%3b",   // ;
		"%3d",   // =
	}

	for _, invalidChar := range invalidChars {
		name = strings.Replace(name, invalidChar, " ", -1)
	}
	return name
}
