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

	e.HideBanner = true

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

	var err error

	// Get the database up and running.
	db, err = c.Connect()
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

	//// Single User
	u := a.Group("/user")

	u.POST("/", userCreate)
	u.POST("/login/", userLogin)

	// Edit User information
	i := u.Group("/edit")

	useJWT(i)

	i.POST("/", userUpdate)
	i.GET("/", userRetrieve)
	i.DELETE("/", userDelete)

	//// Station

	s := a.Group("/station")

	useJWT(s)

	s.POST("/", stationCreate)

	si := s.Group("/:station")

	si.POST("/", notImplemented)
	si.GET("/", notImplemented)
	si.DELETE("/", notImplemented)

	//// Station Users
	r := si.Group("/user")

	r.POST("/", notImplemented)
	r.POST("/:user/", notImplemented)
	r.GET("/:user/", notImplemented)
	r.DELETE("/:user/", notImplemented)

	//// Tracks

	t := a.Group("/track")

	useJWT(t)

	t.POST("/", trackCreate)
	t.POST("/:track/", trackUpdate)
	t.GET("/:track/", trackGet)
	t.DELETE("/:track/", trackDelete)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))

}

func useJWT(c *echo.Group) {
	c.Use(middleware.JWT([]byte(viper.GetString("icii.jwt"))))
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
