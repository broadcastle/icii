package web

import (
	"net/http"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/labstack/echo"
	slugify "github.com/mozillazg/go-slugify"
)

func organizationCreate(c echo.Context) error {

	var org database.Organization

	c.Bind(&org)

	if org.Name == "" {
		return c.JSON(http.StatusMethodNotAllowed, "need a organization name")
	}

	if org.Slug == "" {
		org.Slug = slugify.Slugify(org.Name)
	}

	db.Create(&org)

	userID := getJwtID(c)

	permissions := database.UserPermission{
		OrganizationID: org.ID,
		UserID:         userID,
		TrackEdit:      true,
		TrackAdd:       true,
		TrackRemove:    true,
		UserEdit:       true,
		UserAdd:        true,
		UserRemove:     true,
		StreamEdit:     true,
		StreamAdd:      true,
		StreamRemove:   true,
		OrgEdit:        true,
		OrgAdd:         true,
		OrgRemove:      true,
	}

	db.Create(&permissions)

	return c.JSON(http.StatusOK, org)
}
