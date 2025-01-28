package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env          string `yaml:"env"`
	Port         int    `yaml:"port"`
	MySQLConfig  `yaml:"mysql"`
	SSOConfig    `yaml:"sso"`
	OpenAIConfig `yaml:"openai"`
}

type MySQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"pass"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type SSOConfig struct {
	Secret  string        `yaml:"secret"`
	AppID   int           `yaml:"app_id"`
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type OpenAIConfig struct {
	APIKey string `yaml:"api_key"`
}

func MustLoad() *Config {
	path := getConfigPath()
	return MustLoadByPath(path)
}

func getConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}

func MustLoadByPath(path string) *Config {
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("the file with path %s does not exist", path))
	}
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("error parsing config file")
	}

	config.MySQLConfig.Host = os.Getenv("DB_HOST")
	if config.MySQLConfig.Host == "" {
		config.MySQLConfig.Host = "0.0.0.0"
	}
	fmt.Println("!!! MySQL HOST: " + config.MySQLConfig.Host)

	config.SSOConfig.Host = os.Getenv("SSO_HOST")
	if config.SSOConfig.Host == "" {
		config.SSOConfig.Host = "0.0.0.0"
	}
	fmt.Println("!!! SSO HOST: " + config.SSOConfig.Host)

	return &config
}
