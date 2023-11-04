package cli

import (
	"os"

	"github.com/JannisHajda/docker-backup/internal/borgClient"
	"github.com/JannisHajda/docker-backup/internal/db"
	"github.com/JannisHajda/docker-backup/internal/db/drivers"
	"github.com/spf13/cobra"
)

func getDriver() drivers.Driver {
	driver := drivers.PostgresDriver{
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		Database: os.Getenv("PG_DATABASE"),
		Sslmode:  os.Getenv("PG_SSLMODE"),
	}
	return driver
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  `A command to manage projects. Use 'projects init' to initialize a new project.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var initProject = &cobra.Command{
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

		driver := getDriver()

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

var listProjects = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		driver := getDriver()

		db, err := db.Connect(driver)
		defer db.Close()

		if err != nil {
			println(err.Error())
			return
		}

		projects, err := db.GetAllProjects()

		if err != nil {
			println(err.Error())
			return
		}

		for _, project := range projects {
			println(project.Name)
		}
	},
}

func init() {
	initProject.Flags().StringP("name", "n", "", "Name of the project")
	initProject.MarkFlagRequired("name") // Make the name flag required

	projectsCmd.AddCommand(initProject)
	projectsCmd.AddCommand(listProjects)
}

type Project struct {
	Name string
}
