package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

var srv = &http.Server{
	Addr:         "127.0.0.1:19995",
	ReadTimeout:  2 * time.Second,
	WriteTimeout: 2 * time.Second,
}

var stateCode string
var eventChannel chan string

// Start starts the webserver required for authorisation
func Start(code string, channel chan string) {
	stateCode = code
	eventChannel = channel
	http.HandleFunc("/auth", handlerToken)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handlerToken(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	err := values.Get("error")
	code := values.Get("code")
	state := values.Get("state")

	if err != "" {
		fmt.Println("Received error: " + err)
		_, _ = w.Write([]byte("Error authenticating. Please try again."))
		eventChannel <- "error"
		return
	}
	if state == "" || state != stateCode {
		fmt.Println("Error: Incorrect state code supplied!")
		_, _ = w.Write([]byte("WARNING: Incorrect auth code was supplied. Please try again."))
		eventChannel <- "error"
		return
	}
	if code != "" {
		_, _ = w.Write([]byte("Authorisation successful. You can now close this window."))
		fmt.Println("Authorisation granted")
		eventChannel <- code
	}
}

// Stop stops the webserver gracefully
func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
