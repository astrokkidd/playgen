package main

import (
	"flag"
	"fmt"

	"github.com/astrokkidd/playgen/pkg/controllers"
	"github.com/astrokkidd/playgen/pkg/controllers/spotify/api"
	"github.com/astrokkidd/playgen/pkg/controllers/spotify/auth"
	"github.com/astrokkidd/playgen/pkg/rendering"
	"github.com/astrokkidd/playgen/pkg/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := flag.Int("p", 3000, "port")
	e := echo.New()

	AuthConfig := services.NewAuthConfig()
	authHandler := auth.Handler{
		Config: AuthConfig,
	}

	clientToken, err := authHandler.ClientLogIn()
	if err != nil {
		panic(err)
	}

	APIConfig := services.NewAPIConfig(*clientToken)
	apiHandler := api.Handler{
		Config: APIConfig,
	}

	availableSeedGenres, err := apiHandler.GetAvailableSeedGenres()
	if err != nil {
		panic(err)
	}

	homeController := controllers.HomeController{
		AvailableSeedGenres: *availableSeedGenres,
	}

	renderer := &rendering.TemplateRenderer{
		Templates: services.LoadTemplates(),
	}

	e.Renderer = renderer

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		//middleware.CSRF(),
		middleware.Secure(),
		middleware.Static("web/static"),
	)

	e.GET("/", homeController.RenderHome)
	e.GET("/login", authHandler.LogIn)
	e.GET("/callback", authHandler.CallBack)
	e.POST("/generate", apiHandler.GetRecommendations)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
