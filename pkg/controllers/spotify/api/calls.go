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
	params.Add("target_acousticness", fmt.Sprintf("%f", request.Acousticness))
	params.Add("target_danceability", fmt.Sprintf("%f", request.Danceability))
	params.Add("target_energy", fmt.Sprintf("%f", request.Energy))
	params.Add("target_instrumentalness", fmt.Sprintf("%f", request.Instrumentalness))
	params.Add("target_popularity", fmt.Sprintf("%d", request.Popularity))
	params.Add("target_valence", fmt.Sprintf("%f", request.Valence))

	fmt.Printf("\n\n%s\n\n", params.Encode())
	
	req, _ := http.NewRequest("GET", requestURL+"?"+params.Encode(), nil)
    req.Header.Add("Authorization", "Bearer "+ h.Config.ClientToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var tracks models.RecommendationsResponse

	err = json.NewDecoder(res.Body).Decode(&tracks)
	if err != nil {
		return err
	}

	//return nil
	return c.JSON(http.StatusOK, tracks)
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