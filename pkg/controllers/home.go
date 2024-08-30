package controllers

import (
	"github.com/astrokkidd/playgen/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HomeController struct {
	AvailableSeedGenres models.AvailableSeedGenres
}

func (home *HomeController) RenderHome(c echo.Context) error {
	data := map[string]interface{}{
		"Title":               "playgen",
		"AvailableSeedGenres": home.AvailableSeedGenres,
	}
	return c.Render(http.StatusOK, "content.gotmpl", data)
}
