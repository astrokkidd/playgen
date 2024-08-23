package main

import (
    "html/template"
    "net/http"
    "flag"
    "fmt"
    "io"
    "log/slog"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		slog.Error("template failed to render", slog.String("error", err.Error()))
		return err
	}
	return nil;
}

func home(c echo.Context) error {
	data := map[string]interface{}{
		"Title": "playgen",
	}
	return c.Render(http.StatusOK, "layout.gotmpl", data)
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

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
		
}
