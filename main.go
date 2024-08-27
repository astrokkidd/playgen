package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"crypto/rand"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

type UserTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type ClientTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Secrets struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

type AvailableSeedGenres struct {
	Genres []string `json:"genres"`
}

var (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "http://localhost:3000/callback"
	baseURI      = "/"
)

var availableSeedGenres AvailableSeedGenres

// Controller
func logIn(c echo.Context) error {
	state := generateRandomString(16)
	scope := "user-read-private user-read-email"

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("scope", scope)
	params.Add("redirect_uri", redirectURI)
	params.Add("state", state)

	authURL := "https://accounts.spotify.com/authorize?" + params.Encode()
	return c.Redirect(http.StatusFound, authURL)
}

func clientLogIn(c echo.Context) error {

	readSecrets()

	params := url.Values{}
	params.Add("grant_type", "client_credentials")

	tokenURL := "https://accounts.spotify.com/api/token"

	req, _ := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(params.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var tokenHandler ClientTokenResponse
	err = json.NewDecoder(res.Body).Decode(&tokenHandler)
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "Playgen-Client-Token"
	cookie.Value = tokenHandler.AccessToken
	cookie.Expires = time.Now().Add(time.Hour)

	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/getgenres")
}

// Controller
func callBack(c echo.Context) error {

	params := url.Values{}
	params.Add("code", c.Request().URL.Query()["code"][0])
	params.Add("redirect_uri", redirectURI)
	params.Add("grant_type", "authorization_code")

	tokenURL := "https://accounts.spotify.com/api/token"

	req, _ := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(params.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var tokenHandler UserTokenResponse
	err = json.NewDecoder(res.Body).Decode(&tokenHandler)
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "Playgen-Access-Token"
	cookie.Value = tokenHandler.AccessToken
	cookie.Expires = time.Now().Add(time.Hour)

	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, baseURI)
}

// Controller
func getAvailableSeedGenres(c echo.Context) error {
	requestURL := "https://api.spotify.com/v1/recommendations/available-genre-seeds"

	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)

	req.Header.Add("Authorization", "Bearer "+readCookieValue(c, "Playgen-Client-Token"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(res.Body).Decode(&availableSeedGenres)
	if err != nil {
		return err
	}

	fmt.Printf("\n\n%s\n\n", availableSeedGenres)

	return c.Redirect(http.StatusFound, baseURI)
}

// Util
func generateRandomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Util
func readSecrets() {
	file, _ := os.Open("secrets.json")

	var secrets Secrets

	//err :=
	json.NewDecoder(file).Decode(&secrets)

	clientID = secrets.ClientID
	clientSecret = secrets.ClientSecret
}

// Util
func cookieExists(c echo.Context, name string) int {
	cookie, err := c.Cookie(name)
	if err != nil {
		return -1
	}
	if cookie.Value == "" {
		return 0
	} else {
		return 1
	}
}

// Util
func readCookieValue(c echo.Context, name string) string {
	cookie, err := c.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// Service
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := t.templates.Clone()
	if err != nil {
		return err
	}

	_, err = tmpl.ParseFiles(fmt.Sprintf("templates/%s", name))
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		slog.Error("template failed to render", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func home(c echo.Context) error {
	flag := cookieExists(c, "Playgen-Client-Token")
	if flag < 1 {
		clientLogIn(c)
	}

	data := map[string]interface{}{
		"Title":               "playgen",
		"AvailableSeedGenres": availableSeedGenres,
	}
	return c.Render(http.StatusOK, "content.gotmpl", data)
}

func main() {
	port := flag.Int("p", 3000, "port")
	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.gotmpl")),
	}

	fmt.Printf("\n\n%s\n\n", availableSeedGenres)

	e.Renderer = renderer

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.CSRF(),
		middleware.Secure(),
		middleware.Static("static"),
	)

	//e.Static("/", "static")

	e.GET("/", home)
	e.GET("/login", logIn)
	e.GET("/callback", callBack)
	e.GET("/getgenres", getAvailableSeedGenres)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
