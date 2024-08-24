package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

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

	router := http.NewServeMux()

	router.HandleFunc("POST /login", logIn)

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.gotmpl")),
	}

	e.Renderer = renderer

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "static")

	e.GET("/", home)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))

}

func logIn(w http.ResponseWriter, r *http.Request) {

}
