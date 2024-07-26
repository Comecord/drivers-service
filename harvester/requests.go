package harvester

import (
	"bytes"
	"drivers-service/config"
	"drivers-service/pkg/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthPost struct {
	AuthId        string `json:"AuthId"`
	UserId        string `json:"UserId"`
	User          string `json:"User"`
	Notifications bool   `json:"Notifications"`
}

var logger = logging.NewLogger(config.GetConfig())
var authPostData = &AuthPost{}

func Login() *AuthPost {
	posturl := "https://hosting.glonasssoft.ru/api/v3/auth/login"

	body := []byte(`{
		"login": "demo",
  		"password": "$$demo$$"
	}`)

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected status code: %d", res.StatusCode))
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response Body:", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, authPostData)
	if err != nil {
		panic(err)
	}
	logger.Debugf("Auth: %s, UserId: %s, Username: %s", authPostData.AuthId, authPostData.UserId, authPostData.User)
	return authPostData
}

func GetVehicleList(authData *AuthPost) any {
	fmt.Println(authPostData.AuthId, authPostData.UserId)
	posturl := "https://hosting.glonasssoft.ru/api/monitoringVehicles"
	r, err := http.NewRequest("GET", posturl, nil)
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Auth", authData.AuthId)
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response Body:", string(bodyBytes))

	var dataResponse map[string]interface{}

	err = json.Unmarshal(bodyBytes, &dataResponse)
	if err != nil {
		panic(err)
	}
	return dataResponse
}
