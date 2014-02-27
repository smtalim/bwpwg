//http://stackoverflow.com/questions/14018864/how-do-i-call-mongodb-crud-method-from-handler
//http://blog.gopheracademy.com/day-11-martini
//https://github.com/codegangsta/go-advent-martini/blob/master/templates/list.tmpl
//http://docs.mongohq.com/languages/golang.html
//http://virantha.com/2013/09/29/quick-web-app-with-go-ember-js-and-mongodb/

package main

import (
	"fmt"
	"html/template"
        "labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
	"net/http"
	"os"
)

type Person struct {
        Name string
        Email string
}

func main() {
	http.HandleFunc("/", root)
        http.HandleFunc("/display", display)
        fmt.Println("listening...")
        err := http.ListenAndServe(GetPort(), nil)
        if err != nil {
                panic(err)
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
        <title>Your details</title>
        <style>
          body {
            background-color: #C2A7F2;
            font-family: sans-serif;
          }
          h1 {
            color: #2A1959;
            border-bottom: 2px solid #2A1959;
          }
          h3 {
            color: #474B94;
            font-size: 1.2em;
          }
          #footer {
            clear: both;
            border-top: 1px solid #2A1959;
            text-align: left; 
            height: 50px; 
            font-size: 70%;
            width: 100%;
          }          
        </style>        
      </head>
      <body>
        <h1>A Fun Go App on Heroku to access MongoDB on MongoHQ</h1>
        <h3>Enter the Name</h3>
        <form action="/display" method="post" accept-charset="utf-8">
	  <input type="text" name="name" value="Type the Name..." id="name">
	  <input type="submit" value=".. and query database!">
	</form>
        <div id="footer">
          <p><b>@ Copyright: RubyLearning 2014</b></p>
        </div>	
      </body>
    </html>
`

var displayTemplate = template.Must(template.New("display").Parse(displayTemplateHTML))

func display(w http.ResponseWriter, r *http.Request) {
        // In the open command window set the following for Heroku:
        // heroku config:set MONGOHQ_URL=mongodb://IndianGuru:password@troup.mongohq.com:10080/godata
        uri := os.Getenv("MONGOHQ_URL")
        if uri == "" {
                fmt.Println("no connection string provided")
                os.Exit(1)
        }
 
        sess, err := mgo.Dial(uri)
        if err != nil {
                fmt.Printf("Can't connect to mongo, go error %v\n", err)
                os.Exit(1)
        }
        defer sess.Close()
        
        sess.SetSafe(&mgo.Safe{})
        
        collection := sess.DB("godata").C("user")

        result := Person{}
        //err = 
        collection.Find(bson.M{"name": r.FormValue("name")}).One(&result)

        if result.Email != "" {
                errn := displayTemplate.Execute(w, "The email id you wanted is: " + result.Email)
                if errn != nil {
                        http.Error(w, errn.Error(), http.StatusInternalServerError)
                } 
        } else {
                displayTemplate.Execute(w, "Sorry... The email id you wanted does not exist.")
        }
}

const displayTemplateHTML = ` 
<!DOCTYPE html>
  <html>
    <head>
      <meta charset="utf-8">
      <title>Results</title>
        <style>
          body {
            background-color: #C2A7F2;
            font-family: sans-serif;
          }
          h1 {
            color: #2A1959;
            border-bottom: 2px solid #2A1959;
          }
          h3 {
            color: #474B94;
            font-size: 1.2em;
          }
          #footer {
            clear: both;
            border-top: 1px solid #2A1959;
            text-align: left; 
            height: 50px; 
            font-size: 70%;
            width: 100%;
          }          
        </style>        
    </head>
    <body>
      <h1>A Fun Go App on Heroku to access MongoDB on MongoHQ</h1>
      <p><b>{{html .}}</b></p>
      <p><a href="/">Start again!</a></p>
      <div id="footer">
        <p><b>@ Copyright: RubyLearning 2014</b></p>
      </div>	
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