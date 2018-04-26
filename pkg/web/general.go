package web

import (
	"broadcastle.co/code/icii/pkg/database"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// JSONResponse is used to create a json response.
type JSONResponse struct {
	Msg string `json:"msg"`
}

// This is a fast way to get the user id.
func getJwtID(c echo.Context) (uint, error) {

	i := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64)

	var user database.User

	if err := db.First(&user, uint(i)).Error; err != nil {
		return 0, err
	}

	return uint(i), nil

}

func msg(t interface{}) interface{} {

	switch t.(type) {
	case string:
		return &JSONResponse{Msg: t.(string)}
	default:
		return t
	}
}
