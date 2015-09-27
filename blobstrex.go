package blobstrex

import (
        "fmt"
        "html/template"
        "io"
        "net/http"
        "github.com/rwcarlsen/goexif/exif"
        "strconv"
        "appengine"
        "appengine/blobstore"
)

func init() {
        http.HandleFunc("/", handleRoot)
        http.HandleFunc("/serve/", handleServe)
        http.HandleFunc("/upload", handleUpload)
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
        w.WriteHeader(http.StatusInternalServerError)
        w.Header().Set("Content-Type", "text/plain")
        io.WriteString(w, "Internal Server Error")
        c.Errorf("%v", err)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        uploadURL, err := blobstore.UploadURL(c, "/upload", nil)
        if err != nil {
                serveError(c, w, err)
                return
        }
        w.Header().Set("Content-Type", "text/html")
        err = rootTemplate.Execute(w, uploadURL)
        if err != nil {
                c.Errorf("%v", err)
        }
}

var rootTemplate = template.Must(template.New("root").Parse(rootTemplateHTML))

const rootTemplateHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<link rel="stylesheet" href="css/upper.css">
<title>Upload your Photo</title>
</head>
<body>
<form action="{{.}}" method="POST" enctype="multipart/form-data">
Upload File: <input type="file" name="file"><br />
<input type="submit" name="submit" value="Submit">
</form></body></html>
`

func handleServe(w http.ResponseWriter, r *http.Request) {
        // Instantiate blobstore reader
        reader := blobstore.NewReader(appengine.NewContext(r), 
                                      appengine.BlobKey(r.FormValue("blobKey")))
        
        lat, lng, _ := getLatLng(reader)
        
        blobstore.Delete(appengine.NewContext(r), 
                         appengine.BlobKey(r.FormValue("blobKey")))
        
        if lat == "" {
                io.WriteString(w, "Sorry but your photo has no GeoTag information...")
                return
        }        

        s := "http://maps.googleapis.com/maps/api/staticmap?zoom=5&size=600x300&maptype=roadmap&amp;center="
        s = s + lat + "," + lng + "&markers=color:blue%7Clabel:I%7C" + lat + "," + lng

        img := "<img src='" + s + "' alt='map' />"
        fmt.Fprint(w, img)
        
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        blobs, _, err := blobstore.ParseUpload(r)
        if err != nil {
                serveError(c, w, err)
                return
        }
        file := blobs["file"]
        if len(file) == 0 {
                c.Errorf("no file uploaded")
                http.Redirect(w, r, "/", http.StatusFound)
                return
        }
        http.Redirect(w, r, "/serve/?blobKey="+string(file[0].BlobKey), 
                      http.StatusFound)
}

func getLatLng(f io.Reader) (string, string, error) {
	// Decode parses EXIF-encoded data from f 
	// and returns a queryable Exif object.
	x, err := exif.Decode(f)
	if err != nil {
		return "", "", err
	}
	
        lat, lng, _ := x.LatLong()
        latstr := strconv.FormatFloat(lat, 'f', 15, 64)
        lngstr := strconv.FormatFloat(lng, 'f', 15, 64)
        
	return latstr, lngstr, nil
}

