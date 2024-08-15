package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"

	"github.com/ajugalushkin/url-shortener-version2/cmd"
)

// AppConfig структура параметров заауска.
type AppConfig struct {
	ServerAddress     string `env:"SERVER_ADDRESS"`
	ServerAddressGrpc string `env:"SERVER_ADDRESS_GRPC"`
	BaseURL           string `env:"BASE_URL"`
	FlagLogLevel      string `env:"LOG_LEVEL"`
	FileStoragePath   string `env:"FILE_STORAGE_PATH"`
	DataBaseDsn       string `env:"DATABASE_DSN"`
	SecretKey         string `env:"SECRET_KEY"`
	EnableHTTPS       bool   `env:"ENABLE_HTTPS"`
	Config            string `env:"CONFIG"`
	TrustedSubnet     string `env:"TRUSTED_SUBNET"`
}

// init функция инициализации начальных значений для параметров запуска.
func init() {
	err := godotenv.Load("/docker/.env")
	if err != nil {
		log.Debug("Error loading .env file", "error", err)
	}

	viper.SetDefault("Server_Address", "localhost:8080")
	viper.SetDefault("Server_Address_Grpc", "localhost:3200")
	viper.SetDefault("Base_URL", "http://localhost:8080")
	viper.SetDefault("Log_Level", "Debug")
	viper.SetDefault("File_Storage_PATH", "")
	viper.SetDefault("DataBase_Dsn", "")
	viper.SetDefault("Secret_Key", "")
	viper.SetDefault("Enable_HTTPS", false)
	viper.SetDefault("Config", "")
	viper.SetDefault("Trusted_Subnet", "")
}

// bindToEnv функция для маппинга полей из ENV с полями структуры.
func bindToEnv() {
	_ = viper.BindEnv("Server_Address")
	_ = viper.BindEnv("Server_Address_Grpc")
	_ = viper.BindEnv("Base_URL")
	_ = viper.BindEnv("Log_Level")
	_ = viper.BindEnv("File_Storage_PATH")
	_ = viper.BindEnv("DataBase_Dsn")
	_ = viper.BindEnv("Secret_Key")
	_ = viper.BindEnv("Enable_HTTPS")
	_ = viper.BindEnv("Config")
	_ = viper.BindEnv("Trusted_Subnet")
}

// loadConfiguration function for read json config
func loadConfiguration(file string) AppConfig {
	var config AppConfig
	configFile, err := os.Open(file)
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Warn("Error closing file JSON config", "error", err)
		}
	}(configFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return AppConfig{}
	}
	return config
}

// readConfig функция для чтения конфига.
func readConfig() *AppConfig {
	bindToEnv()
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}

	result := &AppConfig{
		ServerAddress:     viper.GetString("Server_Address"),
		ServerAddressGrpc: viper.GetString("Server_Address_Grpc"),
		BaseURL:           viper.GetString("Base_URL"),
		FlagLogLevel:      viper.GetString("Log_Level"),
		FileStoragePath:   viper.GetString("File_Storage_PATH"),
		DataBaseDsn:       viper.GetString("DataBase_Dsn"),
		SecretKey:         viper.GetString("Secret_Key"),
		EnableHTTPS:       viper.GetBool("Enable_HTTPS"),
		Config:            viper.GetString("Config"),
		TrustedSubnet:     viper.GetString("Trusted_Subnet"),
	}

	if result.Config != "" {
		configJSON := loadConfiguration(result.Config)

		if result.ServerAddress == "" {
			result.ServerAddress = configJSON.ServerAddress
		}
		if result.ServerAddressGrpc == "" {
			result.ServerAddressGrpc = configJSON.ServerAddressGrpc
		}
		if result.BaseURL == "" {
			result.BaseURL = configJSON.BaseURL
		}
		if result.FlagLogLevel == "" {
			result.FlagLogLevel = configJSON.FlagLogLevel
		}
		if result.FileStoragePath == "" {
			result.FileStoragePath = configJSON.FileStoragePath
		}
		if result.DataBaseDsn == "" {
			result.DataBaseDsn = configJSON.DataBaseDsn
		}
		if result.SecretKey == "" {
			result.SecretKey = configJSON.SecretKey
		}
		if result.EnableHTTPS {
			result.EnableHTTPS = configJSON.EnableHTTPS
		}
		if result.TrustedSubnet == "" {
			result.TrustedSubnet = configJSON.TrustedSubnet
		}
	}

	return result
}

// переменные для генерации инстанции
var (
	cfg  *AppConfig
	once sync.Once
)

// GetConfig получение инстанции
func GetConfig() *AppConfig {
	once.Do(
		func() {
			// инициализируем объект
			cfg = readConfig()
		})

	return cfg
}
