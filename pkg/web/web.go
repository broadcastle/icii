package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"broadcastle.co/code/icii/pkg/ice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// var db *gorm.DB

// Start runs the web interface.
func Start(port int) {

	// Init Echo
	e := echo.New()

	e.HideBanner = true

	// Ice database

	if err := ice.Start(); err != nil {
		log.Panic(err)
	}

	defer ice.Close()

	//// Start the database
	// c := database.Config{
	// 	Temp: viper.GetBool("database.temp"),
	// }

	// // If this is not a temporary database.
	// if !c.Temp {
	// 	c.Database = viper.GetString("database.database")
	// 	c.Host = viper.GetString("database.host")
	// 	c.Password = viper.GetString("database.password")
	// 	c.Port = viper.GetInt("database.port")
	// 	c.Postgres = viper.GetBool("database.postgres")
	// 	c.User = viper.GetString("database.user")
	// }

	// var err error

	// // Get the database up and running.
	// db, err = c.Connect()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer db.Close()

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

	si.POST("/", stationUpdate)
	si.GET("/", stationRetrieve)
	si.DELETE("/", stationDelete)

	///////////////////
	// Station Users //
	///////////////////

	user := si.Group("/user")

	user.POST("/", notImplemented)
	user.POST("/:user/", notImplemented)
	user.GET("/:user/", notImplemented)
	user.DELETE("/:user/", notImplemented)

	////////////////////
	// Station Tracks //
	////////////////////

	trk := si.Group("/track")

	trk.POST("/", trackCreate)
	trk.POST("/:track/", trackUpdate)
	trk.GET("/:track/", trackGet)
	trk.DELETE("/:track/", trackDelete)

	trk.POST("/:track/play/", trackPlay)

	////////////////////
	// Station Stream //
	////////////////////

	strm := si.Group("/stream")

	strm.POST("/", streamPost)
	strm.GET("/", streamGet)
	strm.DELETE("/", streamDelete)

	//////////////
	// Playlist //
	//////////////

	p := si.Group("/playlist")

	p.POST("/", playlistCreate)
	p.POST("/:playlist/", playlistUpdate)
	p.GET("/:playlist/", playlistGet)
	p.DELETE("/:playlist/", playlistDelete)

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
