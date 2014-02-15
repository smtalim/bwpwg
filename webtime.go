package main

import (
	"fmt"
	"net/http" //package for http based web programs
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
        fmt.Println("listening...")
        err := http.ListenAndServe(GetPort(), nil)
        if err != nil {
                panic(err)
        }
}

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
        var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "The current time here is %s", time.Now())
}

