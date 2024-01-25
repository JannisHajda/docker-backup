/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"docker-backup/interfaces"
	"docker-backup/internal/db"
	"docker-backup/internal/db/driver"
	"docker-backup/internal/dockerclient"
	"os"

	"github.com/spf13/cobra"
)

func getDbClient() interfaces.DatabaseClient {
	driver, err := driver.NewPostgresDriver("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	client, err := db.NewDatabaseClient(driver)
	if err != nil {
		panic(err)
	}

	return client
}

func getDockerClient() interfaces.DockerClient {
	client, err := dockerclient.NewDockerClient()
	if err != nil {
		panic(err)
	}

	return client
}

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	_ = getDbClient()
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(projectsCmd)
}
