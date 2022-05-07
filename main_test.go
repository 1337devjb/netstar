package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bereths.com/netstar/themoviedb"
	"github.com/stretchr/testify/assert"
)

func GetValidClient() *themoviedb.Client {
	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := themoviedb.NewClient(themoviedbClient, "9764715ae63b1cec4967f30ea42f8ee7", "de-DE", true)
	return themoviedbAPI
}

func TestIndexHandler(t *testing.T) {

	request, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(IndexHandler)

	hf.ServeHTTP(recorder, request)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRouter(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := themoviedb.NewClient(themoviedbClient, "1234", "de-DE", false)

	r := NewRouter(themoviedbAPI)

	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "")

	// Handle any unexpected error
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", resp.StatusCode)
	}

}

func TestStaticFileServer(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := themoviedb.NewClient(themoviedbClient, "1234", "de-DE", false)

	r := NewRouter(themoviedbAPI)
	mockServer := httptest.NewServer(r)

	// get request to assets
	resp, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", resp.StatusCode)
	}

	// result shoud be an html
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}

}

func TestSearchHandlerWithInvalidAPIKey(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := themoviedb.NewClient(themoviedbClient, "1234", "de-DE", false)

	request, err := http.NewRequest("GET", "/search?q=Star%20Wars", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(SearchHandler(themoviedbAPI))

	hf.ServeHTTP(recorder, request)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestSearchHandlerWithValidAPIKey(t *testing.T) {

	themoviedbAPI := GetValidClient()

	r := NewRouter(themoviedbAPI)

	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/search?q=Star%20Wars")

	if err != nil {
		assert.Fail(t, "There should be no error in valid search")
	}

	// result shoud be an html
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	assert.Equal(t, expectedContentType, contentType, "Wrong content type, expected %s, got %s", expectedContentType, contentType)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Status should be %s, got %d", http.StatusOK, resp.StatusCode)
}

func TestTVShowDetailsHandlerWithValidAPIKey(t *testing.T) {

	themoviedbAPI := GetValidClient()

	r := NewRouter(themoviedbAPI)

	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/details?id=1399")

	if err != nil {
		assert.Fail(t, "There should be no error in valid search")
	}

	// result shoud be an html
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	assert.Equal(t, expectedContentType, contentType, "Wrong content type, expected %s, got %s", expectedContentType, contentType)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Status should be %s, got %d", http.StatusOK, resp.StatusCode)
}

func TestSeasonDetailsHandlerWithValidAPIKey(t *testing.T) {

	themoviedbAPI := GetValidClient()

	r := NewRouter(themoviedbAPI)

	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/details/season?id=1399&seasonNumber=1")

	if err != nil {
		assert.Fail(t, "There should be no error in valid search")
	}

	// result shoud be an html
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	assert.Equal(t, expectedContentType, contentType, "Wrong content type, expected %s, got %s", expectedContentType, contentType)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Status should be %s, got %d", http.StatusOK, resp.StatusCode)
}

func TestEpisodeDetailsHandlerWithValidAPIKey(t *testing.T) {

	themoviedbAPI := GetValidClient()

	r := NewRouter(themoviedbAPI)

	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/details/episode?id=1399&seasonNumber=1&episodeNumber=1")

	if err != nil {
		assert.Fail(t, "There should be no error in valid search")
	}

	// result shoud be an html
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	assert.Equal(t, expectedContentType, contentType, "Wrong content type, expected %s, got %s", expectedContentType, contentType)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Status should be %s, got %d", http.StatusOK, resp.StatusCode)
}

func TestInvalidURLWithValidAPIKey(t *testing.T) {

	themoviedbAPI := GetValidClient()

	r := NewRouter(themoviedbAPI)

	mockServer := httptest.NewServer(r)

	ExecuteURL(mockServer.URL+"/search", t)
	ExecuteURL(mockServer.URL+"/details", t)
	ExecuteURL(mockServer.URL+"/details/season", t)
	ExecuteURL(mockServer.URL+"/details/episode", t)
}

func ExecuteURL(url string, t *testing.T) {
	resp, _ := http.Get(url)

	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/plain; charset=utf-8"

	assert.Equal(t, expectedContentType, contentType, "Wrong content type, expected %s, got %s", expectedContentType, contentType)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "Status should be %d, got %d", http.StatusInternalServerError, resp.StatusCode)
}

func TestRunMain(t *testing.T) {
	main()
}
