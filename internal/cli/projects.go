package cli

import (
	"os"

	"github.com/JannisHajda/docker-backup/internal/borgClient"
	"github.com/JannisHajda/docker-backup/internal/db"
	"github.com/JannisHajda/docker-backup/internal/db/drivers"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  `A command to manage projects. Use 'projects init' to initialize a new project.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Long:  `Create a new project with the specified name.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		bc, err := borgClient.NewBorgClient()

		if err != nil {
			println(err.Error())
			return
		}

		err = bc.InitializeRepo(name, "./backup")

		if err != nil {
			println(err.Error())
			return
		}

		driver := drivers.PostgresDriver{
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
			Host:     os.Getenv("PG_HOST"),
			Port:     os.Getenv("PG_PORT"),
			Database: os.Getenv("PG_DATABASE"),
			Sslmode:  os.Getenv("PG_SSLMODE"),
		}

		db, err := db.Connect(driver)
		defer db.Close()

		if err != nil {
			println(err.Error())
			return
		}

		err = db.AddProject(name)

		if err != nil {
			println(err.Error())
			return
		}

		project := Project{Name: name}
		println("Project " + project.Name + " initialized")
	},
}

func init() {
	initCmd.Flags().StringP("name", "n", "", "Name of the project")
	initCmd.MarkFlagRequired("name") // Make the name flag required
	projectsCmd.AddCommand(initCmd)
}

type Project struct {
	Name string
}
