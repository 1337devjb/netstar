package themoviedb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// to generate your structs from a json you can just use https://mholt.github.io/json-to-go/ !

type Client struct {
	http         *http.Client
	key          string
	lang         string
	includeAdult bool
}

type TVShow struct {
	PosterPath       string   `json:"poster_path"`
	Popularity       float64  `json:"popularity"`
	ID               int      `json:"id"`
	BackdropPath     string   `json:"backdrop_path"`
	VoteAverage      float64  `json:"vote_average"`
	Overview         string   `json:"overview"`
	FirstAirDate     string   `json:"first_air_date"`
	OriginCountry    []string `json:"origin_country"`
	GenreIds         []int    `json:"genre_ids"`
	OriginalLanguage string   `json:"original_language"`
	VoteCount        int      `json:"vote_count"`
	Name             string   `json:"name"`
	OriginalName     string   `json:"original_name"`
}

type TVShowDetails struct {
	BackdropPath string `json:"backdrop_path"`
	CreatedBy    []struct {
		ID          int    `json:"id"`
		CreditID    string `json:"credit_id"`
		Name        string `json:"name"`
		Gender      int    `json:"gender"`
		ProfilePath string `json:"profile_path"`
	} `json:"created_by"`
	EpisodeRunTime []int  `json:"episode_run_time"`
	FirstAirDate   string `json:"first_air_date"`
	Genres         []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage         string   `json:"homepage"`
	ID               int      `json:"id"`
	InProduction     bool     `json:"in_production"`
	Languages        []string `json:"languages"`
	LastAirDate      string   `json:"last_air_date"`
	LastEpisodeToAir struct {
		AirDate        string  `json:"air_date"`
		EpisodeNumber  int     `json:"episode_number"`
		ID             int     `json:"id"`
		Name           string  `json:"name"`
		Overview       string  `json:"overview"`
		ProductionCode string  `json:"production_code"`
		SeasonNumber   int     `json:"season_number"`
		StillPath      string  `json:"still_path"`
		VoteAverage    float64 `json:"vote_average"`
		VoteCount      int     `json:"vote_count"`
	} `json:"last_episode_to_air"`
	Name             string      `json:"name"`
	NextEpisodeToAir interface{} `json:"next_episode_to_air"`
	Networks         []struct {
		Name          string `json:"name"`
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		OriginCountry string `json:"origin_country"`
	} `json:"networks"`
	NumberOfEpisodes    int      `json:"number_of_episodes"`
	NumberOfSeasons     int      `json:"number_of_seasons"`
	OriginCountry       []string `json:"origin_country"`
	OriginalLanguage    string   `json:"original_language"`
	OriginalName        string   `json:"original_name"`
	Overview            string   `json:"overview"`
	Popularity          float64  `json:"popularity"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	Seasons []struct {
		AirDate      string `json:"air_date"`
		EpisodeCount int    `json:"episode_count"`
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Overview     string `json:"overview"`
		PosterPath   string `json:"poster_path"`
		SeasonNumber int    `json:"season_number"`
	} `json:"seasons"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Type        string  `json:"type"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

type TVSeasonDetails struct {
	AirDate  string `json:"air_date"`
	Episodes []struct {
		AirDate       string `json:"air_date"`
		EpisodeNumber int    `json:"episode_number"`
		Crew          []struct {
			Department         string  `json:"department"`
			Job                string  `json:"job"`
			CreditID           string  `json:"credit_id"`
			Adult              bool    `json:"adult"`
			Gender             int     `json:"gender"`
			ID                 int     `json:"id"`
			KnownForDepartment string  `json:"known_for_department"`
			Name               string  `json:"name"`
			OriginalName       string  `json:"original_name"`
			Popularity         float64 `json:"popularity"`
			ProfilePath        string  `json:"profile_path"`
		} `json:"crew"`
		GuestStars []struct {
			CreditID           string  `json:"credit_id"`
			Order              int     `json:"order"`
			Character          string  `json:"character"`
			Adult              bool    `json:"adult"`
			Gender             int     `json:"gender"`
			ID                 int     `json:"id"`
			KnownForDepartment string  `json:"known_for_department"`
			Name               string  `json:"name"`
			OriginalName       string  `json:"original_name"`
			Popularity         float64 `json:"popularity"`
			ProfilePath        string  `json:"profile_path"`
		} `json:"guest_stars"`
		ID             int     `json:"id"`
		Name           string  `json:"name"`
		Overview       string  `json:"overview"`
		ProductionCode string  `json:"production_code"`
		SeasonNumber   int     `json:"season_number"`
		StillPath      string  `json:"still_path"`
		VoteAverage    float64 `json:"vote_average"`
		VoteCount      int     `json:"vote_count"`
	} `json:"episodes"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	ID           int    `json:"id"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
	TVID         int
}

type TVEpisodeDetails struct {
	AirDate string `json:"air_date"`
	Crew    []struct {
		ID          int    `json:"id"`
		CreditID    string `json:"credit_id"`
		Name        string `json:"name"`
		Department  string `json:"department"`
		Job         string `json:"job"`
		ProfilePath string `json:"profile_path"`
	} `json:"crew"`
	EpisodeNumber int `json:"episode_number"`
	GuestStars    []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		CreditID    string `json:"credit_id"`
		Character   string `json:"character"`
		Order       int    `json:"order"`
		ProfilePath string `json:"profile_path"`
	} `json:"guest_stars"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ID             int     `json:"id"`
	ProductionCode string  `json:"production_code"`
	SeasonNumber   int     `json:"season_number"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
}

type Results struct {
	Page         int      `json:"page"`
	Results      []TVShow `json:"results"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}

type Result interface {
	Results | TVShowDetails | TVSeasonDetails | TVEpisodeDetails
}

var apiURL = "https://api.themoviedb.org/3"

// creates a new client to use. if lang is empty default language will be en-US
func NewClient(httpClient *http.Client, key string, lang string, includeAdult bool) *Client {
	if lang == "" {
		lang = "en-US"
	}

	return &Client{httpClient, key, lang, includeAdult}
}

func (c *Client) SearchTVShows(query, page string) (*Results, error) {
	endpoint := fmt.Sprintf(apiURL+"/search/tv?query=%s&api_key=%s&language=%s&page=%s&include_adult=%s", url.QueryEscape(query), c.key, c.lang, page, strconv.FormatBool(c.includeAdult))
	return SendRequest[Results](endpoint, c)
}

func (c *Client) GetTVShowDetails(id string) (*TVShowDetails, error) {
	endpoint := fmt.Sprintf(apiURL+"/tv/%s?api_key=%s&language=%s", id, c.key, c.lang)
	return SendRequest[TVShowDetails](endpoint, c)
}

func (c *Client) GetSeasonDetails(id string, seasonNumber string) (*TVSeasonDetails, error) {
	endpoint := fmt.Sprintf(apiURL+"/tv/%s/season/%s?api_key=%s&language=%s", id, seasonNumber, c.key, c.lang)
	details, error := SendRequest[TVSeasonDetails](endpoint, c)
	if details != nil {
		intID, err := strconv.Atoi(id)
		if err == nil {
			details.TVID = intID
		}
	}
	return details, error
}

func (c *Client) GetEpisodeDetails(id string, seasonNumber string, episodeNumber string) (*TVEpisodeDetails, error) {
	endpoint := fmt.Sprintf(apiURL+"/tv/%s/season/%s/episode/%s?api_key=%s&language=%s", id, seasonNumber, episodeNumber, c.key, c.lang)
	return SendRequest[TVEpisodeDetails](endpoint, c)
}

// Generic function to send a simple get request and get a result of T.
// in case we got any error we'll return the error
func SendRequest[T Result](endpoint string, c *Client) (*T, error) {
	body, shouldReturnError, error := GetResponse(c, endpoint)
	if shouldReturnError {
		return new(T), error
	}

	res := new(T)
	return res, json.Unmarshal(body, res)
}

// Gets the response from the get request if no error occurs and the
// http status is 200 it will return the body of the get request
func GetResponse(c *Client, endpoint string) ([]byte, bool, error) {
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, true, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, true, fmt.Errorf(string(body))
	}
	return body, false, nil
}
