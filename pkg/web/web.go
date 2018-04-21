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

	//// Start the database
	c := database.Config{
		Temp: viper.GetBool("database.temp"),
	}

	// If this is not a temporary database.
	if !c.Temp {
		c.Database = viper.GetString("database.database")
		c.Host = viper.GetString("database.host")
		c.Password = viper.GetString("database.password")
		c.Port = viper.GetInt("database.port")
		c.Postgres = viper.GetBool("database.postgres")
		c.User = viper.GetString("database.user")
	}

	// Get the database up and running.
	db, err := c.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Get the middleware up and running.
	e.Static("/", "public")
	e.Pre(middleware.AddTrailingSlash())
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

	//// tracks
	s := a.Group("/track")

	s.Use(middleware.JWT([]byte(viper.GetString("icii.jwt"))))

	s.POST("/", trackCreate)
	s.PUT("/:id/", trackUpdate)
	s.GET("/:id/", trackGet)
	s.DELETE("/:id/", trackDelete)

	//// Users
	u := a.Group("/user")

	u.PUT("/", userCreate)
	u.PUT("/login/", userLogin)
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
