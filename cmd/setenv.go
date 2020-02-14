// Package cmd copyright © 2020 Ivan Marquez <js.ivan.marquez@gmail.com>
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setenvCmd represents the setenv command
var setenvCmd = &cobra.Command{
	Use:   "setenv",
	Short: "set the database environment URL",
	Long: `
	setenv command allows you to setup multiple
	environments for multiple connection URLs.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := cmd.Flags().GetString("name")
		u, _ := cmd.Flags().GetString("url")

		viper.Set(n, u)
		viper.WriteConfig()

		fmt.Println("environment url set:")
		fmt.Printf("%s: %s\n", n, u)
	},
}

func init() {
	rootCmd.AddCommand(setenvCmd)

	setenvCmd.Flags().StringP("name", "n", "", "environment name")
	setenvCmd.Flags().StringP("url", "v", "", "environment url")

	setenvCmd.MarkFlagRequired("name")
	setenvCmd.MarkFlagRequired("url")
}
