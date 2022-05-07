package themoviedb

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func GetValidClient() *Client {
	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := NewClient(themoviedbClient, "9764715ae63b1cec4967f30ea42f8ee7", "de-DE", true)
	return themoviedbAPI
}

func TestClient(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := NewClient(themoviedbClient, "1234", "de-DE", false)

	if themoviedbAPI == nil {
		t.Errorf("Client is null!")
	}

	assert.Equal(t, themoviedbAPI.lang, "de-DE", "lang should equal de-DE")
	assert.Equal(t, themoviedbAPI.key, "1234", "key should equal 1234")
	assert.Equal(t, themoviedbAPI.includeAdult, false, "incloude adult should be false")
}

func TestSearchTVShowsWithInvalidClientAndSearch(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := NewClient(themoviedbClient, "", "", false)

	if themoviedbAPI == nil {
		t.Errorf("Client is null!")
	}

	_, err := themoviedbAPI.SearchTVShows("", "a")

	if err == nil {
		assert.Fail(t, "error should be raised with an invalid query")
	}
	themoviedbAPI = GetValidClient()
	_, err = themoviedbAPI.SearchTVShows("", "a")

	if err == nil {
		assert.Fail(t, "Search with valid client but invalid search should fail!")
	}
}

func TestSearchTVShowsWithValidClientAndSearch(t *testing.T) {

	themoviedbAPI := GetValidClient()

	result, err := themoviedbAPI.SearchTVShows("Game of Thrones", "1")

	if err != nil {
		assert.Fail(t, "Valid search should not throw errors!")
	}

	if result == nil {
		assert.Fail(t, "Result should not be nil for valid search")
	}

	assert.True(t, len(result.Results) > 0, "Result list should be > 0")

}

func TestGetTVShowDetailsWithInvalidClientAndSearch(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := NewClient(themoviedbClient, "", "", false)

	if themoviedbAPI == nil {
		t.Errorf("Client is null!")
	}

	_, err := themoviedbAPI.GetTVShowDetails("")

	if err == nil {
		assert.Fail(t, "error should be raised with an invalid client")
	}
	themoviedbAPI = GetValidClient()
	_, err = themoviedbAPI.GetTVShowDetails("")

	if err == nil {
		assert.Fail(t, "Get details with valid client but invalid id should fail!")
	}
}

func TestGetTVShowDetailsWithValidClientAndSearch(t *testing.T) {

	themoviedbAPI := GetValidClient()

	result, err := themoviedbAPI.GetTVShowDetails("1399")

	if err != nil {
		assert.Fail(t, "Valid id should not throw errors!")
	}

	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, result.Name, "Game of Thrones")

}

func TestGetSeasonDetailsWithInvalidClientAndSearch(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := NewClient(themoviedbClient, "", "", false)

	if themoviedbAPI == nil {
		t.Errorf("Client is null!")
	}

	_, err := themoviedbAPI.GetSeasonDetails("", "")

	if err == nil {
		assert.Fail(t, "error should be raised with an invalid client")
	}
	themoviedbAPI = GetValidClient()
	_, err = themoviedbAPI.GetSeasonDetails("", "")

	if err == nil {
		assert.Fail(t, "Get season details with valid client but invalid id should fail!")
	}

	_, err = themoviedbAPI.GetSeasonDetails("1399", "-1")

	if err == nil {
		assert.Fail(t, "Get season details with valid client but invalid id should fail!")
	}
}

func TestGetSeasonDetailsWithValidClientAndSearch(t *testing.T) {

	themoviedbAPI := GetValidClient()

	result, err := themoviedbAPI.GetSeasonDetails("1399", "1")

	if err != nil {
		assert.Fail(t, "Valid id and seasonNumber should not throw errors!")
	}

	if result == nil {
		assert.Fail(t, "Result should not be nil for valid search")
	}

	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, result.Name, "Staffel 1")
	assert.Equal(t, result.TVID, 1399)

}

func TestGetEpisodeDetailsWithInvalidClientAndSearch(t *testing.T) {

	themoviedbClient := &http.Client{Timeout: 10 * time.Second}
	themoviedbAPI := NewClient(themoviedbClient, "", "", false)

	if themoviedbAPI == nil {
		t.Errorf("Client is null!")
	}

	_, err := themoviedbAPI.GetEpisodeDetails("", "", "")

	if err == nil {
		assert.Fail(t, "error should be raised with an invalid client")
	}
	themoviedbAPI = GetValidClient()
	_, err = themoviedbAPI.GetEpisodeDetails("", "", "")

	if err == nil {
		assert.Fail(t, "Get episode details with valid client but invalid id should fail!")
	}

	_, err = themoviedbAPI.GetEpisodeDetails("1399", "1", "-1")

	if err == nil {
		assert.Fail(t, "Get episode details with valid client but invalid id should fail!")
	}
}

func TestGetEpisodeDetailsWithValidClientAndSearch(t *testing.T) {

	themoviedbAPI := GetValidClient()

	result, err := themoviedbAPI.GetEpisodeDetails("1399", "1", "1")

	if err != nil {
		assert.Fail(t, "Valid id, seasonNumber and episode should not throw errors!")
	}

	if result == nil {
		assert.Fail(t, "Result should not be nil for valid search")
	}

	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, result.Name, "Der Winter naht")

}
