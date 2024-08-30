package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"html/template"
	"os"
)

type Config struct {
	ClientID     string
	ClientSecret string
	ClientToken  string
	RedirectURI  string
	BaseURI      string
}

type Secrets struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func readSecrets() Secrets {
	file, _ := os.Open("secrets.json")

	var secrets Secrets
	//err :=
	json.NewDecoder(file).Decode(&secrets)

	return secrets
}

func CookieExists(c echo.Context, name string) int {
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

func ReadCookieValue(c echo.Context, name string) string {
	cookie, err := c.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func LoadTemplates() *template.Template {
	return template.Must(template.ParseGlob("web/templates/*.gotmpl"))
}

func NewConfig() *Config {
	secrets := readSecrets()
	return &Config{
		ClientID:     secrets.ClientID,
		ClientSecret: secrets.ClientSecret,
		RedirectURI:  "http://localhost:3000/callback",
		BaseURI:      "/",
	}
}
