package web

import (
	"errors"
	"time"

	"broadcastle.co/code/icii/pkg/database"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// User info
type User struct {
	database.User
}

// GetUserFromContext returns a user from a JWT token.
func GetUserFromContext(c echo.Context) (User, error) {

	i := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64)

	var user User

	err := db.First(&user, uint(i)).Error

	return user, err

}

// Create a user.
func (user *User) Create() error {

	var found database.User

	if err := db.Where("email = ?", user.Email).First(&found).Error; err == nil {
		return errors.New("user with that email exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	return db.Create(&user).Error

}

// Login as the user.
func (user *User) Login() (string, error) {

	wrong := errors.New("email and/or password was incorrect")

	if user.Email == "" && user.Password == "" {
		return "", wrong
	}

	var found User

	if err := db.Where("email = ?", user.Email).First(&found).Error; err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password)); err != nil {
		return "", err
	}

	*user = found

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(viper.GetString("icii.jwt")))
	if err != nil {
		return "", err
	}

	return t, nil
}

// Update user with the data from info
func (user *User) Update(info User) error {

	if err := db.First(&user, user.ID).Error; err != nil {
		return err
	}

	// Hash a updated password.
	if info.Password != "" {

		hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		info.Password = string(hash)

	}

	return db.Model(&user).Updates(info).Error
}

// Delete the user
func (user *User) Delete() error {
	return db.Delete(&user).Error
}

// Get a user
func (user *User) Get() error {

	return db.Where(&user).First(&user).Error

}

// CreateStation allows for user to create a station from info.
func (user *User) CreateStation(info Station) error {

	if err := info.Create(); err != nil {
		return err
	}

	permissions := database.UserPermission{
		StationID:     info.ID,
		UserID:        user.ID,
		TrackRead:     true,
		TrackWrite:    true,
		UserRead:      true,
		UserWrite:     true,
		StreamRead:    true,
		StreamWrite:   true,
		StationRead:   true,
		StationWrite:  true,
		ScheduleRead:  true,
		ScheduleWrite: true,
		PlaylistRead:  true,
		PlaylistWrite: true,
	}

	if err := db.Create(&permissions).Error; err != nil {
		return err
	}

	return nil

}
