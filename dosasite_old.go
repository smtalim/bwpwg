package main

import (
        "fmt"
        "net/http"
        "os"
)

func main() {
        fs := http.FileServer(http.Dir("public"))
        http.Handle("/", fs)

        fmt.Println("Listening...")
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
		port = "3000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

