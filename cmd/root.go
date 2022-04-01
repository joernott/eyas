package cmd

import (
	"fmt"
	"os"

	"github.com/joernott/eyas/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "eyas",
	Short: "EYAml Server is a web frontend for encrypting passwords for puppet/hiera",
	Long: `A web frontend to eyaml which allows a user to encrypt passwords for hiera
	without the need to install eyaml on their machine and getting the public keys. Eyas
	allows batch encryption and automatic password generation.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Server(viper.GetBool("ssl"), viper.GetInt("port"), viper.GetString("logfile"), viper.GetString("loglevel"), viper.GetString("keydir"))
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "eyas.yaml", "Config file")
	rootCmd.PersistentFlags().BoolP("ssl", "s", true, "Whether to use SSL or not")
	rootCmd.PersistentFlags().IntP("port", "p", 8443, "Port to run the server on")
	rootCmd.PersistentFlags().StringP("keydir", "k", "keys", "Folder containing directories with PKCS7 public keys")
	rootCmd.PersistentFlags().StringP("logfile", "l", "", "Where to log, defaults to stdout")
	rootCmd.PersistentFlags().StringP("loglevel", "L", "info", "Loglevel can be one of panic, fatal, error, warn, info, debug, trace")
	viper.BindPFlag("ssl", rootCmd.PersistentFlags().Lookup("ssl"))
	viper.SetDefault("ssl", true)
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("port", 8443)
	viper.BindPFlag("keydir", rootCmd.PersistentFlags().Lookup("keydir"))
	viper.SetDefault("keydir", "keys")
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	viper.SetDefault("logfile", "")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.SetDefault("loglevel", "info")
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file ", viper.ConfigFileUsed())
		os.Exit(1)
	}
}
