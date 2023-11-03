package cli

import (
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

type Project struct {
	Name string
}
