package main

import (
	"fmt"
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
        http.HandleFunc("/", handler)
        fmt.Println("listening...")
        err := http.ListenAndServe(GetPort(), nil)
        if err != nil {
                panic(err)
        }
}

func handler(w http.ResponseWriter, r *http.Request) {
        // In the open command window set the following for Heroku:
        // heroku config:set MONGOHQ_URL
           =mongodb://IndianGuru:password@troup.mongohq.com:10080/godata
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
        err = collection.Find(bson.M{"name": "Stefan Klaste"}).One(&result)
        if err != nil {
                panic(err)
        }

        fmt.Fprintf(w, "Output from a Go program on Heroku that accesses MongoDB
                        database on MongoHQ\n\nEmail id: %s", result.Email)
}

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
        var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, 
		            defaulting to " + port)
	}
	return ":" + port
}