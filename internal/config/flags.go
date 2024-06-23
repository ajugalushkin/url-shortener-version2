package config

import (
	"context"
	"flag"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

	viper.SetDefault("Server_Address", ":8080")
	viper.SetDefault("Base_URL", "http://localhost:8080")
	viper.SetDefault("Log_Level", "Info")
	viper.SetDefault("File_Storage_PATH", "/tmp/")
	viper.SetDefault("DataBase_Dsn", "postgres://praktikum:pass@postgres:5432/shortenurls")
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

func bindToFlag() {
	flag.String("a", "", "address and port to run server")
	flag.String("b", "", "Base URL for POST request")
	flag.String("l", "info", "Log level")
	flag.String("f", "", "full name of the file where data in JSON format is saved")
	flag.String("d", "", "DB path for connect")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}

func ReadConfig() *AppConfig {
	bindToFlag()
	bindToEnv()

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
