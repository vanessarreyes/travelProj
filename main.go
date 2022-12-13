// This is the name of our package
// Everything with this package name can see everything
// else inside the same package, regardless of the file they are in
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	// "fmt" has methods for formatted I/O operations (like printing to the console)

	// THe "net/http" library has methods to implement HTTP clients and servers

	"fmt"
	"net/http"
	"text/template"
)

var port = 400

type Page struct {
	Title string
}

func main() {
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))))

	// define handlers
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/NextTravelIdeas", nextTravelIdeasHandler)
	http.HandleFunc("/PastTravels", pastTravelsHandler)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./assets/index.html")
	t.Execute(w, nil)
}

func nextTravelIdeasHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./assets/NextTravelIdeas.html")
	t.Execute(w, nil)
}

func pastTravelsHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./assets/PastTravels.html")
	t.Execute(w, nil)
}
