package harvester

import (
	"bytes"
	"context"
	"drivers-service/config"
	"drivers-service/data/cache"
	"drivers-service/pkg/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AuthPost struct {
	AuthId        string `json:"AuthId"`
	UserId        string `json:"UserId"`
	User          string `json:"User"`
	Notifications bool   `json:"Notifications"`
}

var conf = config.GetConfig()
var ctx = context.Background()

var logger = logging.NewLogger(config.GetConfig())
var authPostData = &AuthPost{}

var rdbInit = cache.InitRedis(conf, ctx)
var rdb = cache.GetRedis()

// Login Получение данных логина из Glonass и запись в redis на 10 минут
func Login() *AuthPost {
	posturl := "https://hosting.glonasssoft.ru/api/v3/auth/login"

	body := []byte(`{
		"login": "demo",
  		"password": "$$demo$$"
	}`)

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected status code: %d", res.StatusCode))
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Response Body:", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, authPostData)
	if err != nil {
		fmt.Println(err)
	}
	logger.Debugf("Auth: %s, UserId: %s, Username: %s", authPostData.AuthId, authPostData.UserId, authPostData.User)

	cache.Set(ctx, rdb, "auth", authPostData.AuthId, 10*time.Minute)

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

func authInterceptor() string {
	authKey, _ := cache.Get(ctx, rdb, "auth")
	if authKey == nil {
		Login()
	}
	return ""
}
