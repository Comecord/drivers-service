package config

import (
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"path/filepath"
	"strings"
)

type GlonassApi struct {
	Domain  string `mapstructure:"domain"`
	Url     string `mapstructure:"url"`
	Monitor struct {
		Vehicle struct {
			Uri         string `mapstructure:"uri"`
			Description string `mapstructure:"description"`
		} `mapstructure:"vehicle"`
	} `mapstructure:"monitor"`
	V3 struct {
		Agents struct {
			Uri string `mapstructure:"uri"`
		} `mapstructure:"agents"`
		Auth struct {
			Uri string `mapstructure:"uri"`
		} `mapstructure:"auth"`
	} `mapstructure:"v3"`
}

type ApiParams struct {
	PathParams  []string
	QueryParams map[string]string
}

var glonassApi GlonassApi

func LoadGlonassApi() (GlonassApi, error) {

	filename, _ := filepath.Abs("config")

	// Настройка Viper
	viper.SetConfigName("glonass")      // имя файла без расширения
	viper.SetConfigType("json")         // формат файла
	viper.AddConfigPath("$HOME/config") // директория с файлом
	viper.AddConfigPath(".")            // директория с файлом
	viper.AddConfigPath(filename)       // директория с файлом

	// Чтение конфигурации
	if err := viper.ReadInConfig(); err != nil {
		return glonassApi, err
	}

	// Преобразование в структуру
	if err := viper.Unmarshal(&glonassApi); err != nil {
		return glonassApi, err
	}
	return glonassApi, nil
}

// RequestApiUrl Данная функция собирает URL из конфигурационного файла json
func RequestApiUrl(uri string, params *ApiParams) string {
	// Формируем базовый URL
	api := glonassApi.Url + uri

	// Если есть параметры пути, заменяем их в URI
	if params.PathParams != nil {
		if len(params.PathParams) > 0 {
			path := strings.Join(params.PathParams, "/")
			api = fmt.Sprintf("%s/%s", api, path)
		}
	}

	// Формируем параметры запроса
	if params.QueryParams != nil {
		query := url.Values{}
		for key, value := range params.QueryParams {
			query.Add(key, value)
		}
		api += "?" + query.Encode()
	}

	return api
}
