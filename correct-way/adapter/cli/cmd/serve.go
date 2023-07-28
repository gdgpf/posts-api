/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"

	"github.com/gdgpf/posts-api/adapter/http/rest"
	"github.com/gdgpf/posts-api/adapter/postgres"
	"github.com/gdgpf/posts-api/adapter/validate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Initialize the REST api",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		database := postgres.GetConnection(ctx)
		defer database.Close()

		postgres.InitializeTranslations()
		postgres.RunMigrations()

		configHTTPServerPort := os.Getenv("PORT")

		if configHTTPServerPort == "" {
			configHTTPServerPort = viper.GetString("server.port")
		}

		validator := validate.InitializeValidator()
		rest.InitializeRest(configHTTPServerPort, database, validator)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
