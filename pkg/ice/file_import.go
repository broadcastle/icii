package ice

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	filetype "gopkg.in/h2non/filetype.v1"
)

// FormImportFile imports file into a temporary folder and checks it.
func FormImportFile(file *multipart.FileHeader) (string, error) {

	logrus.Debug("running ice.FormImportFile()")

	src, err := file.Open()
	if err != nil {
		logrus.Warn(err)
		return "", err
	}

	defer src.Close()

	head := make([]byte, 261)
	if _, err := src.Read(head); err != nil {
		logrus.Warn(err)
		return "", err
	}

	if !filetype.IsMIME(head, "audio/mpeg") {
		logrus.Warn("invalid filetype")
		return "", errors.New("invalid filetype")
	}

	if _, err := src.Seek(0, 0); err != nil {
		logrus.Warn(err)
		return "", err
	}

	u := uuid.NewV4()
	ext := path.Ext(file.Filename)
	tmp := path.Join(viper.GetString("files.location"), u.String()+ext)

	dst, err := os.Create(tmp)
	if err != nil {
		logrus.Warn(err)
		return "", err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		logrus.Warn(err)
		return "", err
	}

	logrus.Debug("completed ice.FormImportFile()")

	return tmp, nil
}
