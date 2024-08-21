package main

import (
	"html/template"
	"net/http"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	t := &Template {
		templates: template.Must(template.ParseGlob("public/views/index.html")),
	}
	
	e := echo.New()
	e.Renderer = t

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		data := map[string]interface{}{
			"Message": "Hello world!",
		}
		return c.Render(http.StatusOK, "index.html", data)
	})

	e.Logger.Fatal(e.Start(":8080"))

}
