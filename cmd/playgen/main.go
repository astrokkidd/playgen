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

var (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "http://localhost:3000/callback"
	baseURI      = "/"
)

func main() {
	port := flag.Int("p", 3000, "port")
	e := echo.New()

	config := services.NewConfig()

	authHandler := auth.Handler{
		Config: config,
	}

	clientToken, err := authHandler.ClientLogIn()
	if err != nil {
		panic(err)
	}

	availableSeedGenres, err := api.GetAvailableSeedGenres(*clientToken)
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
		middleware.CSRF(),
		middleware.Secure(),
		middleware.Static("web/static"),
	)

	e.GET("/", homeController.RenderHome)
	e.GET("/login", authHandler.LogIn)
	e.GET("/callback", authHandler.CallBack)
	//e.GET("/getgenres", controllers.getAvailableSeedGenres)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
