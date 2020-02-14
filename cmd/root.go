// Package cmd copyright © 2020 Ivan Marquez
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mdb",
	Short: "A CLI for better management of MongoDB commands",
	Long: `
	mdb is a CLI that empowers MongoDB commands.

	This application allows you to run the most
	common commands for importing and exporting
	data from MongoDB. It can be configured with
	environment variables so that you don't have
	to spend time searching for urls and credentials.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".env", "config file")
	rootCmd.PersistentFlags().StringP("env", "e", "", "runtime environment (dev, staging, prod, etc.)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func validateEnv(cmd *cobra.Command, args []string) error {
	env, _ := cmd.Flags().GetString("env")

	if env == "" {
		return fmt.Errorf("specify which environment will be used to run commands")
	} else if e := viper.GetString(env); e == "" {
		return fmt.Errorf("unrecognized environment: %s", env)
	}

	return nil
}
