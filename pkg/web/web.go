package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

var db *gorm.DB

// Start runs the web interface.
func Start(port int) {

	// Init Echo
	e := echo.New()

	// Get the database up and running.
	d, err := database.Config{Temp: true}.Connect() // Set to temp for development usage.
	if err != nil {
		log.Fatal(err)
	}

	defer d.Close()

	db = d

	// Get the middleware up and running.
	e.Static("/", "public")
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	// Set up the renderer
	e.Renderer = &Template{templates: template.Must(template.ParseGlob("templates/*"))}

	// Frontpage
	e.GET("/", displayPage)

	/////////////
	//// API ////
	/////////////

	// API
	a := e.Group("/api/v1")

	//// Songs
	s := a.Group("/song")

	s.Use(middleware.JWT([]byte(viper.GetString("icii.jwt"))))

	s.PUT("/", songCreate)
	s.PUT("/:id", songUpdate)
	s.GET("/:id", songGet)
	s.DELETE("/:id", songDelete)

	//// Users
	u := a.Group("/user")

	u.PUT("/", userCreate)
	u.PUT("/login", userLogin)
	// u.PUT("/:id", notImplemented)
	// u.GET("/:id", notImplemented)
	// u.DELETE("/:id", notImplemented)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))

}

func displayPage(c echo.Context) error {

	w := "World"
	if c.Param("id") != "" {
		w = c.Param("id")
	}

	return c.Render(http.StatusOK, "index", w)
}

func notImplemented(c echo.Context) error {
	return c.JSON(http.StatusNotFound, "this has not been implemented yet")
}
