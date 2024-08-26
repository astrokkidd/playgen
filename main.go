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

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type Secrets struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

var (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "http://localhost:3000/callback"
	baseURI      = "/"
)

// Controller
func logIn(c echo.Context) error {
	readSecrets()

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

	var tokenHandler TokenResponse
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
	data := map[string]interface{}{
		"Title": "playgen",
	}
	return c.Render(http.StatusOK, "content.gotmpl", data)
}

func main() {
	port := flag.Int("p", 3000, "port")
	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.gotmpl")),
	}

	e.Renderer = renderer

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "static")

	e.GET("/", home)
	e.GET("/login", logIn)
	e.GET("/callback", callBack)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
