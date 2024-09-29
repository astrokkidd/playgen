package models

type UserTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type ClientTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type AvailableSeedGenres struct {
	Genres []string `json:"genres"`
}

type Recommendation struct {
	Acousticness     float32      `json:"target_acousticness"`
	Danceability     float32      `json:"target_danceability"`
	Energy           float32      `json:"target_energy"`
	Instrumentalness float32      `json:"target_instrumentalness"`
	Popularity       int      `json:"target_popularity"`
	Valence          float32      `json:"target_valence"`
	Limit	         string      `json:"limit"`
	Genres           string `json:"seed_genres"`
}

type Image struct {
	Url string `json:"url"`
	Height int `json:"height"`
	Width int `json:"width"`
}

type Album struct {
	Name string `json:"name"`
	Images []Image `json:"images"`
}

type Artist struct {
	Name string `json:"name"`
	ID string `json:"id"`
}

type Track struct {
	Name string `json:"name"`
	URI string `json:"uri"`
	URL struct { SpotifyURL string `json:"spotify"`} `json:"external_urls"`
	Artists []Artist `json:"artists"`
	Album Album `json:"album"`
}

type RecommendationsResponse struct {
	Tracks []Track `json:"tracks"`
}