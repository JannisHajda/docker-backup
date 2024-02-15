/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"docker-backup/internal/worker"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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

var backupContainerCmd = &cobra.Command{
	// ./docker-backup backup --container test-container --passphrase test --output backups
	Use:   "backup",
	Short: "Backup a container",
	Long:  `Backup a container to a borg repository`,
	Run: func(cmd *cobra.Command, args []string) {
		container, _ := cmd.Flags().GetString("container")
		passphrase, _ := cmd.Flags().GetString("passphrase")
		output, _ := cmd.Flags().GetString("output")

		if container == "" || passphrase == "" || output == "" {
			fmt.Printf("container, passphrase, and output are required\n")
			return
		}

		w, err := worker.NewWorker(container, output, passphrase)
		if err != nil {
			panic(err)
		}

		defer w.Stop()

		err = w.Backup()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(backupContainerCmd)
	backupContainerCmd.Flags().StringP("container", "c", "", "Container to backup")
	backupContainerCmd.Flags().StringP("passphrase", "p", "", "Passphrase to use for encryption")
	backupContainerCmd.Flags().StringP("output", "o", "", "Output volume to use for backup")
}
