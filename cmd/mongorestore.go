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

// mongorestoreCmd represents the mongorestore command
var mongorestoreCmd = &cobra.Command{
	Use:   "mongorestore",
	Short: "import data from a MongoDB dump",
	Long: `
	mongorestore program loads data from either 
	a binary database dump created by mongodump or 
	the standard input into a mongod or mongos instance.
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
			"mongorestore",
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
			"/tmp/" + path.Base(file),
		}

		internal.InitializeMDBContainer(c, []string{file})

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mongorestoreCmd)

	mongorestoreCmd.Flags().StringP("collection", "c", "", "specifies the name of the destination collection for mongorestore to restore data into when restoring from a BSON file")
	mongorestoreCmd.Flags().StringP("file", "f", "", "specifies location and name of the database dump")

	mongorestoreCmd.MarkFlagRequired("collection")
	mongorestoreCmd.MarkFlagRequired("file")
}
