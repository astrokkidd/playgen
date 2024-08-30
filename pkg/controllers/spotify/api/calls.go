package api

import (
	"encoding/json"
	"net/http"

	"github.com/astrokkidd/playgen/pkg/models"
)

func GetAvailableSeedGenres(clientToken string) (*models.AvailableSeedGenres, error) {
	requestURL := "https://api.spotify.com/v1/recommendations/available-genre-seeds"

	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)

	req.Header.Add("Authorization", "Bearer "+clientToken)

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
