package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"

	"github.com/ajugalushkin/url-shortener-version2/cmd"
)

type AppConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FlagLogLevel    string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDsn     string `env:"DATABASE_DSN"`
	SecretKey       string `env:"SECRET_KEY"`
}

func init() {
	err := godotenv.Load("/docker/.env")
	if err != nil {
		log.Debug("Error loading .env file", "error", err)
	}

	viper.SetDefault("Server_Address", "localhost:8080")
	viper.SetDefault("Base_URL", "http://localhost:8080")
	viper.SetDefault("Log_Level", "Debug")
	viper.SetDefault("File_Storage_PATH", "")
	viper.SetDefault("DataBase_Dsn", "")
	viper.SetDefault("Secret_Key", "")
}

func bindToEnv() {
	_ = viper.BindEnv("Server_Address")
	_ = viper.BindEnv("Base_URL")
	_ = viper.BindEnv("Log_Level")
	_ = viper.BindEnv("File_Storage_PATH")
	_ = viper.BindEnv("DataBase_Dsn")
	_ = viper.BindEnv("Secret_Key")
}

func ReadConfig() *AppConfig {
	bindToEnv()
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}

	result := &AppConfig{
		ServerAddress:   viper.GetString("Server_Address"),
		BaseURL:         viper.GetString("Base_URL"),
		FlagLogLevel:    viper.GetString("Log_Level"),
		FileStoragePath: viper.GetString("File_Storage_PATH"),
		DataBaseDsn:     viper.GetString("DataBase_Dsn"),
		SecretKey:       viper.GetString("Secret_Key"),
	}
	return result
}

type ctxConfig struct{}

func ContextWithFlags(ctx context.Context, config *AppConfig) context.Context {
	return context.WithValue(ctx, ctxConfig{}, config)
}

func FlagsFromContext(ctx context.Context) *AppConfig {
	if config, ok := ctx.Value(ctxConfig{}).(*AppConfig); ok {
		return config
	}
	return &AppConfig{}
}
