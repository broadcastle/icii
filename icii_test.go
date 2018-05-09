package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"broadcastle.co/code/icii/cmd"
	"broadcastle.co/code/icii/pkg/database"
	"broadcastle.co/code/icii/pkg/web"
)

func TestIcii(t *testing.T) {

	base := "http://localhost:8080/api/v1/"
	userSignUpURL := base + "user/"
	userLoginURL := userSignUpURL + "login/"
	station := base + "station/"
	trackURL := base + "track/"

	cmd.Execute()

	go web.Start(8080)

	time.Sleep(time.Second * 1)

	user := database.User{
		Name:     "User Name",
		Email:    "user@name.com",
		Password: "fake password",
	}

	user2 := database.User{
		Name:     "Another User",
		Email:    "another@user.com",
		Password: "fake password",
	}

	// Create First User
	if _, err := testPOST(userSignUpURL, user); err != nil {
		t.Error(err)
	}

	// Check if user is only able to use one email.
	if r, _ := testPOST(userSignUpURL, user); r.StatusCode == http.StatusOK {
		t.Error("user was able to use the same email twice")
	}

	// Create Second User
	if _, err := testPOST(userSignUpURL, user2); err != nil {
		t.Error(err)
	}

	// Login as the first user and get a token.
	rul, err := testPOST(userLoginURL, user)
	if err != nil {
		t.Error(err)
	}

	defer rul.Body.Close()

	msg := web.JSONResponse{}

	if err := json.NewDecoder(rul.Body).Decode(&msg); err != nil {
		t.Error(err)
	}

	token := msg.Msg

	// Change users name.
	if _, err := testTokenPOST(userSignUpURL+"edit/", token, database.User{Name: "Updated Name"}); err != nil {
		t.Error(err)
	}

	// Create a station.

	org := database.Station{Name: "Generic Station"}

	resp, err := testTokenPOST(station, token, org)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp)
	}

	// Upload a file.

	p := map[string]string{
		"title":   "loping string",
		"station": "1",
		"genre":   "testing",
		"year":    "2018",
	}

	resp, err = uploadFile(trackURL, token, p, "./test_audio/loping_sting.mp3")
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp)
	}

	return

}

func testPOST(u string, i interface{}) (*http.Response, error) {

	j, err := json.Marshal(&i)
	if err != nil {
		return nil, err
	}

	return http.Post(u, "application/json", bytes.NewReader(j))

}

func testTokenGET(u, t string) (*http.Response, error) {

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+t)

	client := http.Client{}

	return client.Do(req)
}

func testTokenPOST(u, t string, i interface{}) (*http.Response, error) {

	j, err := json.Marshal(&i)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+t)
	req.Header.Add("content-type", "application/json")

	client := http.Client{}

	return client.Do(req)
}

func uploadRequest(uri, token string, params map[string]string, paramName, path string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	for key, val := range params {
		writer.WriteField(key, val)
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("content-type", writer.FormDataContentType())

	return req, err

}

func uploadFile(uri, token string, params map[string]string, path string) (*http.Response, error) {

	req, err := uploadRequest(uri, token, params, "audio", path)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)

}
