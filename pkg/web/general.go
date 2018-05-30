package web

import (
	"strconv"

	"broadcastle.co/code/icii/pkg/ice"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// JSONResponse is used to create a json response.
type JSONResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

// This is a fast way to get the user id.
func getJwtID(c echo.Context) (uint, error) {

	i := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64)

	id := uint(i)

	var user ice.User
	user.ID = id

	if err := user.Get(); err != nil {
		return 0, err
	}

	return id, nil

}

func getIDfromParam(name string, c echo.Context) (uint, error) {

	i := c.Param(name)

	id, err := strconv.Atoi(i)
	if err != nil {
		return 0, err
	}

	return uint(id), nil

}

func getStationID(i string) (uint, error) {

	stationID, err := strconv.Atoi(i)
	if err != nil {
		return 0, err
	}

	return uint(stationID), nil

}

func msg(status int, t interface{}) (int, interface{}) {

	switch t.(type) {
	case string:
		return status, &JSONResponse{Status: status, Msg: t.(string)}
	default:
		return status, t
	}
}
