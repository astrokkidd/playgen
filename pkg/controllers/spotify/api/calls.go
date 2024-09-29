package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"
	//"io"
	
	//"strings"
	"github.com/astrokkidd/playgen/pkg/models"
	"github.com/astrokkidd/playgen/pkg/services"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Config *services.Config
}

func (h *Handler) GetAvailableSeedGenres() (*models.AvailableSeedGenres, error) {

	requestURL := "https://api.spotify.com/v1/recommendations/available-genre-seeds"

	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)

	req.Header.Add("Authorization", "Bearer "+ h.Config.ClientToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var availableSeedGenres models.AvailableSeedGenres

	err = json.NewDecoder(res.Body).Decode(&availableSeedGenres)
	if err != nil {
		return nil, err
	}

	return &availableSeedGenres, nil
}

func (h *Handler) GetRecommendations(c echo.Context) (error) {
	var request models.Recommendation
	if err := c.Bind(&request); err != nil {
		return err
	}

	requestURL := "https://api.spotify.com/v1/recommendations"

	params := url.Values{}
	params.Add("limit", request.Limit)
	params.Add("seed_genres", request.Genres)

	fmt.Printf("\n\n%s\n\n", request.Genres)

	params.Add("target_acousticness", fmt.Sprintf("%d", request.Acousticness/100))
	params.Add("target_danceability", fmt.Sprintf("%d", request.Danceability/100))
	params.Add("target_energy", fmt.Sprintf("%d", request.Energy/100))
	params.Add("target_instrumentalness", fmt.Sprintf("%d", request.Instrumentalness/100))
	params.Add("target_popularity", fmt.Sprintf("%d", request.Popularity/100))
	params.Add("target_valence", fmt.Sprintf("%d", request.Valence/100))

	//params.Add("limit", "10")
	//params.Add("seed_genres", "edm")
	//params.Add("target_acousticness", "0.5")
	//params.Add("target_danceability", "0.5")
	//params.Add("target_energy", "0.5")
	//params.Add("target_instrumentalness", "0.5")
	//params.Add("target_popularity", "0.5")
	//params.Add("target_valence", "0.5")
	
	req, _ := http.NewRequest("GET", requestURL+"?"+params.Encode(), nil)
    req.Header.Add("Authorization", "Bearer "+ h.Config.ClientToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	//defer res.Body.Close() // Ensure the response body is closed after reading

    /* Read the response body
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return err
    }*/

    // Print the response status code and body
    //fmt.Printf("Status: %s\n", res.Status)
    //fmt.Printf("Response Body: %s\n", string(body))

	var tracks models.RecommendationsResponse

	err = json.NewDecoder(res.Body).Decode(&tracks)
	if err != nil {
		return err
	}

	//return nil
	return c.JSON(http.StatusOK, tracks.Tracks)
}

/*func searchArtists(clientToken string) (*models.[]Artist, error) {
	requestURL := "https://api.spotify.com/v1/search"

	params := url.Values{}
	params.Add("q", query)
	params.Add("type", "artist")
	params.Add("limit", "10")

	req, _ := http.NewRequest("GET", spotifyURL+"?"+params.Encode(), nil)
    req.Header.Add("Authorization", "Bearer "+clientToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

}*/