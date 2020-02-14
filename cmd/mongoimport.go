// Package cmd copyright © 2020 Ivan Marquez <js.ivan.marquez@gmail.com>
package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ivan-marquez/mdb/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mongoimportCmd represents the mongoimport command
var mongoimportCmd = &cobra.Command{
	Use:   "mongoimport",
	Short: "Import data to a MongoDB collection",
	Long: `
	mongoimport tool imports content from 
	an Extended JSON, CSV, or TSV export created 
	by mongoexport, or potentially, another 
	third-party export tool.
	`,
	PreRunE: validateEnv,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, _ := cmd.Flags().GetString("env")
		collection, _ := cmd.Flags().GetString("collection")
		file, _ := cmd.Flags().GetString("file")

		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("The specified file path does not exist: %s", file)
		}

		conn, _ := url.Parse(viper.GetString(env))
		pw, _ := conn.User.Password()
		URL := []string{
			conn.Query().Get("replicaSet"),
			conn.Hostname() + ":" + conn.Port(),
		}

		c := []string{
			"mongoimport",
			"--host",
			strings.Join(URL, "/"),
			"--collection",
			collection,
			"--db",
			path.Base(conn.Path),
			"--ssl",
			"--authenticationDatabase",
			"admin",
			"--username",
			conn.User.Username(),
			"--password",
			pw,
			"--file",
			"/tmp/" + path.Base(file),
		}

		internal.InitializeMDBContainer(c, []string{file})

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mongoimportCmd)

	mongoimportCmd.Flags().StringP("collection", "c", "", "specifies the collection to import")
	mongoimportCmd.Flags().StringP("file", "f", "", "specifies location and name of a file containing the data to import")

	mongoimportCmd.MarkFlagRequired("collection")
	mongoimportCmd.MarkFlagRequired("file")
}
