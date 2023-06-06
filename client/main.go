package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func login(userName string, password string) (*string, error) {
	jsonBody, _ := json.Marshal(map[string]interface{}{
		"username": userName,
		"password": password,
	})

	request, err := http.NewRequest("POST", "http://localhost:8395/auth/login", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	// Set the request headers
	request.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	token, exist := responseBody["token"]
	if !exist {
		return nil, errors.New("error!")
	}
	res := token.(string)
	return &res, nil
}

func getMyInfo(token string) (*string, error) {
	req, err := http.NewRequest("GET", "http://localhost:8395/api/user/me", nil)
	if err != nil {
		return nil, err
	}
	auth := "Bearer " + token
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := string(body)
	return &res, nil
}

func main() {
	token, err := login("admin", "admin12!@")
	if err != nil {
		panic(err)
	}
	fmt.Println("서버에서 받은 토큰: ", *token)
	res, err := getMyInfo(*token)
	if err != nil {
		panic(err)
	}
	fmt.Println("서버에서 받은 사용자이름: ", *res)
}
