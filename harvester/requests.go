package harvester

import (
	"bytes"
	"context"
	"drivers-service/config"
	"drivers-service/data/cache"
	"drivers-service/pkg/logging"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"strings"
	"time"
)

type AuthPost struct {
	AuthId        string `json:"AuthId"`
	UserId        string `json:"UserId"`
	User          string `json:"User"`
	Notifications bool   `json:"Notifications"`
}

var (
	conf  = config.GetConfig()
	ctx   = context.Background()
	gl, _ = config.LoadGlonassApi()

	logger       = logging.NewLogger(config.GetConfig())
	authPostData = &AuthPost{}

	rdbInit = cache.InitRedis(conf, ctx)
	rdb     = cache.GetRedis()
)

// TODO: Заменить все URL из переменных запросов к Glonass на RequestApiUrl с чтением кофиг json
// TODO: Отлавливать ошибки вместо выдачи паники при запросах

// Login Получение данных логина из Glonass и запись в redis на 10 минут
func Login() *AuthPost {
	params := &config.ApiParams{
		PathParams:  []string{"login"},
		QueryParams: nil,
	}

	posturl := config.RequestApiUrl(gl.V3.Auth.Uri, params)

	body := []byte(fmt.Sprintf(`{
		"login": "%s",
  		"password": "%s"
	}`, conf.Glonass.AuthLogin, conf.Glonass.AuthPassword))

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

	err = cache.Set(ctx, rdb, "auth", authPostData.AuthId, 10*time.Minute)
	if err != nil {
		fmt.Println(err)
	}
	return authPostData
}

// GetVehicleList Получение монитора транспорта
func GetVehicleList() any {
	posturl := config.RequestApiUrl(gl.Monitor.Vehicle.Uri, nil)
	authValue := authInterceptor()

	r, err := http.NewRequest("GET", posturl, nil)
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Auth", authValue)
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
	if rdb == nil {
		fmt.Println("Ошибка: rdb не инициализирован.")
		return ""
	}

	authKey, err := rdb.Get(ctx, "auth").Result()
	if errors.Is(err, redis.Nil) {
		Login()
		authKey, _ = rdb.Get(ctx, "auth").Result()
	}
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.Trim(authKey, "\"")
}
