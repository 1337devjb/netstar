package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"bereths.com/netstar/themoviedb"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// declare template
var index = template.Must(template.ParseFiles("pages/index.html", "pages/base.html"))
var details = template.Must(template.ParseFiles("pages/details.html", "pages/base.html"))
var seasonDetails = template.Must(template.ParseFiles("pages/season_details.html", "pages/base.html"))
var episodeDetails = template.Must(template.ParseFiles("pages/episode_details.html", "pages/base.html"))

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *themoviedb.Results
}

// define the route urls here
func newRouter(themoviedbAPI *themoviedb.Client) *mux.Router {
	r := mux.NewRouter()
	// index
	r.HandleFunc("/", indexHandler).Methods("GET")
	// search like /search?q=Star Wars
	r.HandleFunc("/search", searchHandler(themoviedbAPI)).Methods("GET")
	// details like /search?id=1337
	r.HandleFunc("/details", tvShowDetailsHandler(themoviedbAPI)).Methods("GET")
	// details for seasion like /search?id=1337&seasonNumber=1
	r.HandleFunc("/details/season", seasonDetailsHandler(themoviedbAPI)).Methods("GET")
	// details for episode like /details/episode?id=1337&seasonNumber=1&episodeNumber=4
	r.HandleFunc("/details/episode", episodeDetailsHandler(themoviedbAPI)).Methods("GET")

	// declare static files
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// add static files prefix to router
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	return r
}

// index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	err := index.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

// handles the search a user executes
func searchHandler(themoviedbAPI *themoviedb.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := themoviedbAPI.SearchTVShows(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Search Query is: ", searchQuery)
		log.Println("Page is: ", page)
		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: results.TotalResults,
			Results:    results,
		}

		buf := &bytes.Buffer{}
		err = index.ExecuteTemplate(w, "base", search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}

// handles the tv show details if a user clicks on a tv show
func tvShowDetailsHandler(themoviedbAPI *themoviedb.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		id := params.Get("id")

		results, err := themoviedbAPI.GetTVShowDetails(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf := &bytes.Buffer{}
		err = details.ExecuteTemplate(w, "base", results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}

// handles the season details if a user klicks on a season
func seasonDetailsHandler(themoviedbAPI *themoviedb.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		id := params.Get("id")
		seasonNumber := params.Get("seasonNumber")

		result, err := themoviedbAPI.GetSeasonDetails(id, seasonNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf := &bytes.Buffer{}
		err = seasonDetails.ExecuteTemplate(w, "base", result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}

// handles the episode a user clicks
func episodeDetailsHandler(themoviedbAPI *themoviedb.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		id := params.Get("id")
		seasonNumber := params.Get("seasonNumber")
		episodeNumber := params.Get("episodeNumber")

		result, err := themoviedbAPI.GetEpisodeDetails(id, seasonNumber, episodeNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf := &bytes.Buffer{}
		err = episodeDetails.ExecuteTemplate(w, "base", result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}

func main() {

	// load .env if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// set port of the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // set default port to 3000 if no port is set
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set") // api key has to be set for this to work!
	}

	lang := os.Getenv("LANGUAGE")
	if lang == "" {
		lang = "en-US" // default language is en-US
	}

	includeAdultProperty := os.Getenv("INCLUDE_ADULT")
	if includeAdultProperty == "" {
		includeAdultProperty = "False"
	}

	includeAdult, err := strconv.ParseBool(includeAdultProperty)

	if err != nil {
		log.Fatal(err)
		includeAdult = false // in error case we set adult only to false for child protection ;)
	}

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := themoviedb.NewClient(themoviedbClient, apiKey, lang, includeAdult)

	// declare router
	r := newRouter(themoviedbAPI)

	// serve
	http.ListenAndServe(":"+port, r)
}
