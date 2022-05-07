package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"bereths.com/netstar/themoviedb"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
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

type Config struct {
	APIKey       string `mapstructure:API_KEY`
	Language     string `mapstructure:LANGUAGE`
	IncludeAdult bool   `mapstructure:INCLUDE_ADULT`
	Port         string `mapstructure:PORT`
}

// define the route urls here
func NewRouter(themoviedbAPI *themoviedb.Client) *mux.Router {
	r := mux.NewRouter()
	// index
	r.HandleFunc("/", IndexHandler).Methods("GET")
	// search like /search?q=Star Wars
	r.HandleFunc("/search", SearchHandler(themoviedbAPI)).Methods("GET")
	// details like /search?id=1337
	r.HandleFunc("/details", TVShowDetailsHandler(themoviedbAPI)).Methods("GET")
	// details for seasion like /search?id=1337&seasonNumber=1
	r.HandleFunc("/details/season", SeasonDetailsHandler(themoviedbAPI)).Methods("GET")
	// details for episode like /details/episode?id=1337&seasonNumber=1&episodeNumber=4
	r.HandleFunc("/details/episode", EpisodeDetailsHandler(themoviedbAPI)).Methods("GET")

	// declare static files
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// add static files prefix to router
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	return r
}

// index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	err := index.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

// handles the search a user executes
func SearchHandler(themoviedbAPI *themoviedb.Client) http.HandlerFunc {
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
func TVShowDetailsHandler(themoviedbAPI *themoviedb.Client) http.HandlerFunc {
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
func SeasonDetailsHandler(themoviedbAPI *themoviedb.Client) http.HandlerFunc {
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
func EpisodeDetailsHandler(themoviedbAPI *themoviedb.Client) http.HandlerFunc {
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

// loads a config file with the name ".env" from a given path and returns a config or
// an error if config could not be read
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}

func main() {

	config, err := LoadConfig(".")

	if err != nil {
		log.Fatalf("Could not read config!")
	}

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := themoviedb.NewClient(themoviedbClient, config.APIKey, config.Language, config.IncludeAdult)

	// declare router
	r := NewRouter(themoviedbAPI)

	// serve
	http.ListenAndServe(":"+config.Port, r)
}
