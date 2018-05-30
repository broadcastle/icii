package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
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
	userCreate := base + "user/"
	userLogin := userCreate + "login/"
	userEdit := userCreate + "edit/"
	station := base + "station/"
	track := base + "track/"

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

	// Create the first user.
	if _, err := testPOST(userCreate, user); err != nil {
		t.Error(err)
	}

	// Check if user is only able to use one email.
	if r, _ := testPOST(userCreate, user); r.StatusCode == http.StatusOK {
		t.Error("user was able to use the same email twice")
	}

	// Create the second user.
	if _, err := testPOST(userCreate, user2); err != nil {
		t.Error(err)
	}

	// Login as the first user and get a token.
	token, err := login(userLogin, user)
	if err != nil {
		t.Error(err)
	}

	// Test a login with the wrong password.
	if _, err := login(userLogin, database.User{Email: "user@name.com", Password: "wrong password"}); err == nil {
		t.Error("able to log in with incorrect password")
	}

	// Change users name.
	if err := token.post(userEdit, database.User{Name: "Updated Name"}); err != nil {
		t.Error(err)
	}

	// Create a station.
	if err := token.post(station, database.Station{Name: "Generic Station"}); err != nil {
		t.Error(err)
	}

	// Upload a file.

	p := map[string]string{
		"title":   "loping string",
		"station": "1",
		"genre":   "testing",
		"year":    "2018",
	}

	if err := token.upload(track, p, "./test_audio/loping_sting.mp3"); err != nil {
		t.Error(err)
	}

	// Upload the file again.
	if err := token.upload(track, p, "./test_audio/loping_sting.mp3"); err != nil {
		t.Error(err)
	}

	// Change the title of the first track
	if err := token.post(track+"1/", database.Track{Title: "Loping String 2"}); err != nil {
		t.Error(err)
	}

	// Delete the duplicate.
	if err := token.delete(track + "2/"); err != nil {
		t.Error(err)
	}

	return

}

type Token struct {
	String string
}

func login(u string, user database.User) (Token, error) {

	var t Token

	rul, err := testPOST(u, user)
	if err != nil {
		return t, err
	}

	defer rul.Body.Close()

	if rul.StatusCode != http.StatusOK {
		return t, errors.New(string(rul.StatusCode) + " instead of 200 ")
	}

	msg := web.JSONResponse{}

	if err := json.NewDecoder(rul.Body).Decode(&msg); err != nil {
		return t, err
	}

	t.String = msg.Msg

	return t, nil

}

func (t Token) post(u string, i interface{}) error {

	j, err := json.Marshal(&i)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u, bytes.NewReader(j))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+t.String)
	req.Header.Add("content-type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := statusError(u, resp); err != nil {
		return err
	}

	return nil

}

func statusError(link string, resp *http.Response) error {

	if resp.StatusCode != http.StatusOK {

		d, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return err
		}

		log.Printf("unable to post to %v\nreason: %v\n%s", link, resp.StatusCode, string(d))
	}

	return nil
}

func (t Token) get(u string) error {

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+t.String)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	log.Printf("unable to retrieve %v\nreason:\n%v\n%v", u, resp.StatusCode, resp.Body)
	// }
	if err := statusError(u, resp); err != nil {
		return err
	}

	return nil

}

func (t Token) delete(u string) error {

	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+t.String)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	log.Printf("Unable to delete %v: %v\n%v", u, resp.StatusCode, resp.Body)
	// }
	if err := statusError(u, resp); err != nil {
		return err
	}

	return nil

}

func (t Token) upload(u string, params map[string]string, path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio", filepath.Base(path))
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, file); err != nil {
		return err
	}

	for key, val := range params {
		writer.WriteField(key, val)
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+t.String)
	req.Header.Set("content-type", writer.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	log.Println(resp.Body)
	// }
	if err := statusError(u, resp); err != nil {
		return err
	}

	return nil

}

func testPOST(u string, i interface{}) (*http.Response, error) {

	j, err := json.Marshal(&i)
	if err != nil {
		return nil, err
	}

	return http.Post(u, "application/json", bytes.NewReader(j))

}
