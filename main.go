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
var DB *sql.DB = nil

// TODO: create different structs for future/past travel data
// once tables begin to have different data

// Struct needs to change to match database tables
type NewTravelItem struct {
	ID          int    `json:"id"`
	Place       string `json:"place"`
	Description string `json:"description"`
}

func ConnectToDB() *sql.DB {
	// The `sql.Open` function opens a new `*sql.DB` instance. We specify the driver name
	// and the URI for our database. Here, we're using a Postgres URI
	db, err := sql.Open("pgx", "user= password= host=localhost port=5432 dbname=travel-db sslmode=disable")
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
	DB = ConnectToDB()

	// define handlers
	r.HandleFunc("/", mainPageHandler).Methods("GET")

	// Future Travels
	r.HandleFunc("/NextTravelIdeas", getNextTravelIdeasPage).Methods("GET")
	r.HandleFunc("/FutureTravelCards", getNextTravelCards).Methods("GET")
	r.HandleFunc("/NextTravelForm", getNextTravelFormPage).Methods("GET")
	r.HandleFunc("/SubmitNewTravelIdea", postNewTravel).Methods("POST")

	// Past Travels
	r.HandleFunc("/PastTravels", getPastTravelsPage).Methods("GET")
	r.HandleFunc("/PastTravelsCards", getPastTravelCards).Methods("GET")
	r.HandleFunc("/PastTravelForm", getPastTravelFormPage).Methods("GET")
	r.HandleFunc("/SubmitPastTravel", postPastTravel).Methods("POST")

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
	data, err := DB.Query("SELECT * FROM future_travels")
	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}

	// create array to hold data
	futureTravels := []NewTravelItem{}

	// iterate over the returned rows
	// we can go over to the next row by calling the `Next` method, which will
	// return `false` if there are no more rows
	for data.Next() {
		travel := NewTravelItem{}
		// create an instance of `Bird` and write the result of the current row into it
		if err := data.Scan(&travel.ID, &travel.Place, &travel.Description); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		// append the current instance to the slice of birds
		futureTravels = append(futureTravels, travel)
	}

	// Convert to json
	travelListBytes, err := json.Marshal(futureTravels)

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

	// TODO: need to have better way of giving ids
	row := DB.QueryRow("SELECT COUNT(*) FROM future_travels")
	if err := row.Scan(&newTravel.ID); err != nil {
		log.Fatalf("could not scan row: %v", err)
	}

	// Get the information about the travel from the form info
	newTravel.Place = r.Form.Get("place")
	newTravel.Description = r.Form.Get("description")
	newTravel.ID++ // increment id for new item

	fmt.Print(newTravel)

	// the `Exec` method returns a `Result` type instead of a `Row`
	// we follow the same argument pattern to add query params
	result, err := DB.Exec("INSERT INTO future_travels (id, place, description) VALUES ($1, $2, $3)", newTravel.ID, newTravel.Place, newTravel.Description)
	if err != nil {
		log.Fatalf("could not insert row: %v", err)
	}

	// the `Result` type has special methods like `RowsAffected` which returns the
	// total number of affected rows reported by the database
	// In this case, it will tell us the number of rows that were inserted using
	// the above query
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
	}
	// we can log how many rows were inserted
	fmt.Println("inserted", rowsAffected, "rows")

	//Finally, we redirect the user to the original HTMl page
	// (located at `/assets/`), using the http libraries `Redirect` method
	http.Redirect(w, r, "http://localhost:400/NextTravelIdeas", http.StatusFound)
}

func getPastTravelsPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./assets/PastTravels.html")
	t.Execute(w, nil)
}

func getPastTravelCards(w http.ResponseWriter, r *http.Request) {
	data, err := DB.Query("SELECT * FROM past_travels")
	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}

	// create array to hold data
	pastTravels := []NewTravelItem{}

	// iterate over the returned rows
	// we can go over to the next row by calling the `Next` method, which will
	// return `false` if there are no more rows
	for data.Next() {
		travel := NewTravelItem{}
		// create an instance of `Bird` and write the result of the current row into it
		if err := data.Scan(&travel.ID, &travel.Place, &travel.Description); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		// append the current instance to the slice of birds
		pastTravels = append(pastTravels, travel)
	}

	// Convert the "nextTravels" variable to json
	travelListBytes, err := json.Marshal(pastTravels)

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

func getPastTravelFormPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./assets/PastTravelForm.html")
	t.Execute(w, nil)
}

func postPastTravel(w http.ResponseWriter, r *http.Request) {
	// Crate a new travel object
	pastTravel := NewTravelItem{}

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

	// TODO: need to have better way of giving ids
	row := DB.QueryRow("SELECT COUNT(*) FROM past_travels")
	if err := row.Scan(&pastTravel.ID); err != nil {
		log.Fatalf("could not scan row: %v", err)
	}

	// Get the information about the travel from the form info
	pastTravel.Place = r.Form.Get("place")
	pastTravel.Description = r.Form.Get("description")
	pastTravel.ID++ // increment id for new item

	fmt.Print(pastTravel)

	// the `Exec` method returns a `Result` type instead of a `Row`
	// we follow the same argument pattern to add query params
	result, err := DB.Exec("INSERT INTO past_travels (id, place, description) VALUES ($1, $2, $3)", pastTravel.ID, pastTravel.Place, pastTravel.Description)
	if err != nil {
		log.Fatalf("could not insert row: %v", err)
	}

	// the `Result` type has special methods like `RowsAffected` which returns the
	// total number of affected rows reported by the database
	// In this case, it will tell us the number of rows that were inserted using
	// the above query
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
	}
	// we can log how many rows were inserted
	fmt.Println("inserted", rowsAffected, "rows")

	//Finally, we redirect the user to the original HTMl page
	// (located at `/assets/`), using the http libraries `Redirect` method
	http.Redirect(w, r, "http://localhost:400/PastTravels", http.StatusFound)
}
