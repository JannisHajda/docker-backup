package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docker-backup",
	Short: "Docker Backup is a tool to backup and restore docker containers",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example: `,
	Run: func(cmd *cobra.Command, args []string) {
		println("Hello World")
	},
}

func Start() {
	rootCmd.AddCommand(projectsCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}
