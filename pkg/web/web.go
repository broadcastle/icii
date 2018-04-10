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
	e.Renderer = &Template{templates: template.Must(template.ParseGlob("public/views/*.html"))}

	// Frontpage
	e.GET("/", displayPage)

	/////////////
	//// API ////
	/////////////

	// API
	av := "/api/v1/"

	songs := av + "song"

	// api/v1/song
	e.PUT(songs, uploadSong)

	// api/v1/song/:id
	songID := songs + "/:id"

	e.GET(songID, getSong)
	e.PUT(songID, updateSong)
	e.DELETE(songID, deleteSong)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))

}

func displayPage(c echo.Context) error {

	w := "World"
	if c.Param("id") != "" {
		w = c.Param("id")
	}

	return c.Render(http.StatusOK, "index", w)
}
