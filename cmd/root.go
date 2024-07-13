package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// структура для root уровня.
var (
	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "Application URL Shortener",
		Long:  `Application URL Shortener`,
	}
)

// Execute позволяет вызывать root.Execute из другого пакета.
func Execute() error {
	return rootCmd.Execute()
}

// init функция позволяет считать параметры запуска из флагов,
// для чтения используется cobra + viper.
func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringP("ServerAddress", "a", "", "address and port to run server")
	rootCmd.PersistentFlags().StringP("BaseURL", "b", "", "Base URL for POST request")
	rootCmd.PersistentFlags().StringP("LogLevel", "l", "", "Log level")
	rootCmd.PersistentFlags().StringP("FileStoragePATH", "f", "", "full name of the file where data in JSON format is saved")
	rootCmd.PersistentFlags().StringP("DataBaseDsn", "d", "", "DB path for connect")

	_ = viper.BindPFlag("Server_Address", rootCmd.PersistentFlags().Lookup("ServerAddress"))
	_ = viper.BindPFlag("Base_URL", rootCmd.PersistentFlags().Lookup("BaseURL"))
	_ = viper.BindPFlag("Log_Level", rootCmd.PersistentFlags().Lookup("LogLevel"))
	_ = viper.BindPFlag("File_Storage_PATH", rootCmd.PersistentFlags().Lookup("FileStoragePATH"))
	_ = viper.BindPFlag("DataBase_Dsn", rootCmd.PersistentFlags().Lookup("DataBaseDsn"))
}
