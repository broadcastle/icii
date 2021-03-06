package web

import (
	"errors"
	"net/http"

	"broadcastle.co/code/icii/pkg/ice"
	"github.com/labstack/echo"
)

func streamPost(c echo.Context) error {
	return errors.New("nothing in web.streamPost")
}

func streamUpdate(c echo.Context) error {

	// Add auth check.

	var update ice.Stream
	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	stream := ice.InitStream()

	if err := ice.Echo(c, stream); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := ice.Update(stream, update); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, stream)
}

func streamGet(c echo.Context) error {

	// Add auth check.

	stream := ice.InitStream()

	if err := ice.Echo(c, stream); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, stream)

}

func streamDelete(c echo.Context) error {

	stream := ice.InitStream()

	if err := ice.Echo(c, stream); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := ice.Remove(stream); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, "success")

}

func streamPlayTrack(c echo.Context) error {

	stream := ice.InitStream()
	track := ice.InitTrack()

	if err := ice.Echo(c, stream, track); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := stream.(*ice.Stream).Play(*track.(*ice.Track)); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, "success")
}
