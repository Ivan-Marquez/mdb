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

// mongoexportCmd represents the mongoexport command
var mongoexportCmd = &cobra.Command{
	Use:   "mongoexport",
	Short: "export data from a MongoDB collection",
	Long: `
	mongoexport is a command-line tool that produces 
	a JSON or CSV export of data stored in a MongoDB 
	instance.
	`,
	PreRunE: validateEnv,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, _ := cmd.Flags().GetString("env")
		collection, _ := cmd.Flags().GetString("collection")

		conn, _ := url.Parse(viper.GetString(env))
		pw, _ := conn.User.Password()
		URL := []string{
			conn.Query().Get("replicaSet"),
			conn.Hostname() + ":" + conn.Port(),
		}

		c := []string{
			"mongoexport",
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
			"-o",
			"output.json",
		}

		internal.InitializeMDBContainer(c, nil)

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		fmt.Printf("\noutput path: %s/output.json\n", dir)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mongoexportCmd)

	mongoexportCmd.Flags().StringP("collection", "c", "", "collection to export")

	mongoexportCmd.MarkFlagRequired("collection")
	mongoexportCmd.MarkFlagRequired("env")
}
