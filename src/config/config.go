package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"log"
	"os"
	"time"
)

type Config struct {
	Server struct {
		IPort   int    `yaml:"internalPort"`
		EPort   int    `yaml:"externalPort"`
		RunMode string `yaml:"runMode"`
	}
	SMTP struct {
		EmailFrom string `yaml:"emailFrom"`
		Host      string `yaml:"smtpHost"`
		Pass      string `yaml:"smtpPass"`
		Port      string `yaml:"smtpPort"`
		User      string `yaml:"smtpUser"`
		Auth      string `yaml:"smtpAuth"`
		Security  string `yaml:"smtpSecurity"`
	}
	Logger struct {
		FilePath string `yaml:"filePath"`
		Encoding string `yaml:"encoding"`
		Level    string `yaml:"level"`
		Logger   string `yaml:"logcod"`
	}
	Cors struct {
		AllowOrigins string `yaml:"allowOrigins"`
	}
	MongoX struct {
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		Database    string `yaml:"database"`
		ReplicaName string `yaml:"replicaName"`
		Replication string `yaml:"replication"`
		AuthSource  string `yaml:"authSource"`
	}
	Redis struct {
		Host               string        `yaml:"host"`
		Port               int           `yaml:"port"`
		Password           string        `yaml:"password"`
		Db                 int           `yaml:"db"`
		DialTimeout        time.Duration `json:"dialTimeout"`
		ReadTimeout        time.Duration `json:"readTimeout"`
		WriteTimeout       time.Duration `json:"writeTimeout"`
		IdleCheckFrequency time.Duration `json:"idleCheckFrequency"`
		PoolSize           int           `json:"poolSize"`
		PoolTimeout        time.Duration `json:"poolTimeout"`
	}
	Password struct {
		IncludeChars     bool `yaml:"includeChars"`
		IncludeDigits    bool `yaml:"includeDigits"`
		MinLength        int  `yaml:"minLength"`
		MaxLength        int  `yaml:"maxLength"`
		IncludeUppercase bool `yaml:"includeUppercase"`
		IncludeLowercase bool `yaml:"includeLowercase"`
	}
	Otp struct {
		ExpireTime time.Duration `yaml:"expireTime"`
		Digits     int           `yaml:"digits"`
		Limiter    time.Duration `yaml:"limiter"`
	}
	Jwt struct {
		Secret                     string        `yaml:"secret"`
		RefreshSecret              string        `yaml:"refreshSecret"`
		AccessTokenExpireDuration  time.Duration `yaml:"accessTokenExpireDuration"`
		RefreshTokenExpireDuration time.Duration `yaml:"refreshTokenExpireDuration"`
	}
}

func GetConfig() *Config {
	cfgPath := getConfigPath(os.Getenv("APP_ENV"))
	b, err := LoadConfig(cfgPath, "yml")
	if err != nil {
		log.Fatalf("Error in load config %v", err)
	}

	cfg, err := ParseConfig(b)
	if err != nil {
		log.Fatalf("Error in parse config %v", err)
	}

	return cfg
}

func ParseConfig(b []byte) (*Config, error) {
	var cnf Config
	err := yaml.Unmarshal(b, &cnf)
	if err != nil {
		fmt.Printf("Erro in parse Config: %v", err)
	}
	return &cnf, nil
}

func LoadConfig(filename string, fileType string) ([]byte, error) {
	yamlFile, err := os.ReadFile(filename + "." + fileType)
	if err != nil {
		return nil, err
	}
	return yamlFile, nil
}

func getConfigPath(env string) string {
	if env == "docker" {
		return "/app/config/config-docker"
	} else if env == "production" {
		return "docker/config-production"
	} else {
		return "docker/config-development"
	}
}
