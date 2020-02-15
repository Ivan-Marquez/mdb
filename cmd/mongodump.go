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

// mongodumpCmd represents the mongodump command
var mongodumpCmd = &cobra.Command{
	Use:   "mongodump",
	Short: "create a binary export from MongoDB",
	Long: `
	mongodump is a utility for creating a binary export 
	of the contents of a database. mongodump can export 
	data from either mongod or mongos instances; 
	
	i.e. can export data from standalone, replica set, 
	and sharded cluster deployments.
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
			"mongodump",
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
		}

		internal.InitializeMDBContainer(c, nil)

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		fmt.Printf("\ndump path: %s/dump\n", dir)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mongodumpCmd)

	mongodumpCmd.Flags().StringP("collection", "c", "", "collection to export")

	mongodumpCmd.MarkFlagRequired("collection")
}
