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
	rootCmd.PersistentFlags().BoolP("EnableHTTPS", "s", false, "Enable HTTPS in the web server")
	rootCmd.PersistentFlags().BoolP("Config", "c", false, "set JSON config file")
	rootCmd.PersistentFlags().StringP("TrustedSubnet", "t", "", "set CIDR")

	_ = viper.BindPFlag("Server_Address", rootCmd.PersistentFlags().Lookup("ServerAddress"))
	_ = viper.BindPFlag("Base_URL", rootCmd.PersistentFlags().Lookup("BaseURL"))
	_ = viper.BindPFlag("Log_Level", rootCmd.PersistentFlags().Lookup("LogLevel"))
	_ = viper.BindPFlag("File_Storage_PATH", rootCmd.PersistentFlags().Lookup("FileStoragePATH"))
	_ = viper.BindPFlag("DataBase_Dsn", rootCmd.PersistentFlags().Lookup("DataBaseDsn"))
	_ = viper.BindPFlag("Enable_HTTPS", rootCmd.PersistentFlags().Lookup("EnableHTTPS"))
	_ = viper.BindPFlag("Config", rootCmd.PersistentFlags().Lookup("Config"))
	_ = viper.BindPFlag("Trusted_Subnet", rootCmd.PersistentFlags().Lookup("TrustedSubnet"))
}
