package rendering

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log/slog"
)

type TemplateRenderer struct {
	Templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := t.Templates.Clone()
	if err != nil {
		return err
	}

	_, err = tmpl.ParseFiles(fmt.Sprintf("web/templates/%s", name))
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
