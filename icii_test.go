package main

import (
	"bytes"
	"encoding/json"
	"net/http"
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
	orgCreate := base + "organization/"

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

	// Create a station.

	org := database.Station{Name: "Generic Station"}

	resp, err := testTokenPOST(orgCreate, token, org)
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

func testTokenGET(u string, t string) (*http.Response, error) {

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+t)

	client := http.Client{}

	return client.Do(req)
}

func testTokenPOST(u string, t string, i interface{}) (*http.Response, error) {

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
