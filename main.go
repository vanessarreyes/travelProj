// This is the name of our package
// Everything with this package name can see everything
// else inside the same package, regardless of the file they are in
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	// "fmt" has methods for formatted I/O operations (like printing to the console)

	// THe "net/http" library has methods to implement HTTP clients and servers

	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"database/sql"
	"log"

	// we have to import the driver, but don't use it in our code
	// so we use the `_` symbol
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/gorilla/mux"
)

var port = 400

type NewTravelItem struct {
	Place       string `json:"place"`
	Description string `json:"description"`
}

var nextTravels = []NewTravelItem{
	NewTravelItem{
		Place:       "Europe",
		Description: "I want to go here",
	},
	NewTravelItem{
		Place:       "Las Vegas",
		Description: "I want to go here",
	},
	NewTravelItem{
		Place:       "Europe",
		Description: "I want to go here",
	},
	NewTravelItem{
		Place:       "Las Vegas",
		Description: "I want to go here",
	},
	NewTravelItem{
		Place:       "Europe",
		Description: "I want to go here",
	},
}

func ConnectToDB() *sql.DB {
	// The `sql.Open` function opens a new `*sql.DB` instance. We specify the driver name
	// and the URI for our database. Here, we're using a Postgres URI
	db, err := sql.Open("pgx", "user=postgres password=2400 host=localhost port=5432 dbname=travel-db sslmode=verify-ca pool_max_conns=10")
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// To verify the connection to our database instance, we can call the `Ping`
	// method. If no error is returned, we can assume a successful connection
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	fmt.Println("database is reachable")

	return db
}

func main() {
	r := mux.NewRouter()

	// connect to DB
	// ConnectToDB()

	// define handlers
	r.HandleFunc("/", mainPageHandler).Methods("GET")

	// Future Travels
	r.HandleFunc("/NextTravelIdeas", getNextTravelIdeasPage).Methods("GET")
	r.HandleFunc("/TravelIdeasCards", getNextTravelCards).Methods("GET")
	r.HandleFunc("/NextTravelForm", getNextTravelFormPage).Methods("GET")
	r.HandleFunc("/SubmitNewTravelIdea", postNewTravel).Methods("POST")

	// Past Travels
	r.HandleFunc("/PastTravels", getPastTravelsPage).Methods("GET")

	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))).Methods("GET")

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./assets/index.html")
	t.Execute(w, nil)
}

func getNextTravelIdeasPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./assets/NextTravelIdeas.html")
	t.Execute(w, nil)
}

func getNextTravelCards(w http.ResponseWriter, r *http.Request) {
	// Convert the "nextTravels" variable to json
	travelListBytes, err := json.Marshal(nextTravels)

	// If there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If all goes well, write the JSON list of birds to the response
	w.Write(travelListBytes)
}

func getNextTravelFormPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./assets/NextTravelForm.html")
	t.Execute(w, nil)
}

func postNewTravel(w http.ResponseWriter, r *http.Request) {
	// Crate a new travel object
	newTravel := NewTravelItem{}

	// We send all our data as HTML form data
	// the `ParseForm` method of the request, parses the
	// form values
	err := r.ParseForm()

	// In case of any error, we respond with an error to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the information about the travel from the form info
	newTravel.Place = r.Form.Get("place")
	newTravel.Description = r.Form.Get("description")

	// Append our existing list of travel with a new entry
	nextTravels = append(nextTravels, newTravel)

	//Finally, we redirect the user to the original HTMl page
	// (located at `/assets/`), using the http libraries `Redirect` method
	http.Redirect(w, r, "http://localhost:400", http.StatusFound)
}

func getPastTravelsPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./assets/PastTravels.html")
	t.Execute(w, nil)
}
