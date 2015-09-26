package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// A struct is a collection of fields
// This is our type which matches the JSON object.
type IpRecord struct {
	// These two fields use the json: tag to specify which field they map to
	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`
	// These fields are mapped directly by name (note the different case)
	City string
	Ip   string
	// As these fields can be nullable, we use a pointer to a string rather than a string
	Lat *string
	Lng *string
}

func main() {
        http.HandleFunc("/", handler)
        fmt.Println("listening...")
        err := http.ListenAndServe(GetPort(), nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
                return
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
        ip := "198.252.210.32"

	// QueryEscape escapes the ip string so
	// it can be safely placed inside a URL query
	safeIp := url.QueryEscape(ip)

        url := fmt.Sprintf("http://api.hostip.info/get_json.php?position=true&ip=%s", safeIp)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record IpRecord
	
        // Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	fmt.Fprintf(w, "Latitude = %s and Longitude = %s", *record.Lat, *record.Lng)
}
