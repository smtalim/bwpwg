package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Add a handler to handle serving static files from a specified directory
	// The reason for using StripPrefix is that you can change the served 
	// directory as you please, but keep the reference in HTML the same.
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	
	http.HandleFunc("/", root)
	http.HandleFunc("/upper", upper)
        fmt.Println("listening...")
        err := http.ListenAndServe(GetPort(), nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
                return
        }
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, rootForm)
}

const rootForm = `
  <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <link rel="stylesheet" href="css/upper.css">
        <title>String Upper</title>
      </head>
      <body>
        <h1>String Upper</h1>
        <p>The String Upper Service will accept a string from you and 
           return you the Uppercase version of the original string. Have fun!</p>
        <form action="/upper" method="post" accept-charset="utf-8">
	  <input type="text" name="str" value="Type a string..." id="str">
	  <input type="submit" value=".. and change to uppercase!">
        </form>
      </body>
    </html>
`
var upperTemplate = template.Must(template.New("upper").Parse(upperTemplateHTML))

func upper(w http.ResponseWriter, r *http.Request) {
        strEntered := r.FormValue("str")
        strUpper := strings.ToUpper(strEntered)
        err := upperTemplate.Execute(w, strUpper)
        if err != nil {
	        http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

const upperTemplateHTML = ` 
<!DOCTYPE html>
  <html>
    <head>
      <meta charset="utf-8">
      <link rel="stylesheet" href="css/upper.css">
      <title>String Upper Results</title>
    </head>
    <body>
      <h1>String Upper Results</h1>
      <p>The Uppercase of the string that you had entered is:</p>
      <pre>{{html .}}</pre>
    </body>
  </html>
`

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



