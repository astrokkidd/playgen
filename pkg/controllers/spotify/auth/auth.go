package auth

import (
	"github.com/astrokkidd/playgen/pkg/models"
	"github.com/astrokkidd/playgen/pkg/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"

	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type Handler struct {
	Config *services.Config
}

func (h *Handler) LogIn(c echo.Context) error {
	state := services.GenerateRandomString(16)
	scope := "user-read-private user-read-email"

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", h.Config.ClientID)
	params.Add("scope", scope)
	params.Add("redirect_uri", h.Config.RedirectURI)
	params.Add("state", state)

	authURL := "https://accounts.spotify.com/authorize?" + params.Encode()
	return c.Redirect(http.StatusFound, authURL)
}

func (h *Handler) CallBack(c echo.Context) error {

	params := url.Values{}
	params.Add("code", c.Request().URL.Query()["code"][0])
	params.Add("redirect_uri", h.Config.RedirectURI)
	params.Add("grant_type", "authorization_code")

	tokenURL := "https://accounts.spotify.com/api/token"

	req, _ := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(params.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(h.Config.ClientID+":"+h.Config.ClientSecret)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var tokenHandler models.UserTokenResponse
	err = json.NewDecoder(res.Body).Decode(&tokenHandler)
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "Playgen-Access-Token"
	cookie.Value = tokenHandler.AccessToken
	cookie.Expires = time.Now().Add(time.Hour)

	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, h.Config.BaseURI)
}

func (h *Handler) ClientLogIn() (*string, error) {

	params := url.Values{}
	params.Add("grant_type", "client_credentials")

	tokenURL := "https://accounts.spotify.com/api/token"

	req, _ := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(params.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(h.Config.ClientID+":"+h.Config.ClientSecret)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var tokenHandler models.ClientTokenResponse
	err = json.NewDecoder(res.Body).Decode(&tokenHandler)
	if err != nil {
		return nil, err
	}

	return &tokenHandler.AccessToken, nil
}
