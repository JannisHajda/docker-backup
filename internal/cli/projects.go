package cli

import (
	"os"

	"github.com/JannisHajda/docker-backup/internal/borgClient"
	"github.com/JannisHajda/docker-backup/internal/db"
	"github.com/JannisHajda/docker-backup/internal/db/drivers"
	"github.com/JannisHajda/docker-backup/internal/dockerClient"
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

		bc, err := borgClient.NewBorgClient("./backup")

		if err != nil {
			println(err.Error())
			return
		}

		err = bc.InitializeRepo(name)

		if err != nil {
			println(err.Error())
			return
		}

		driver := getDriver()

		var errs []error

		database, err := db.NewDatabase(driver)
		if err != nil {
			errs = append(errs, err)
			return
		}

		project, err := database.AddProject(name)
		if err != nil {
			errs = append(errs, err)
			return
		}

		if len(errs) > 0 {
			for _, err := range errs {
				println(err.Error())
			}
			return
		}

		println("Project " + project.Name + " initialized")
	},
}

var listProjects = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		driver := getDriver()
		database, err := db.NewDatabase(driver)

		if err != nil {
			println(err.Error())
			return
		}

		projects, err := database.GetAllProjects()

		if err != nil {
			println(err.Error())
			return
		}

		for _, project := range projects {
			println(project.Name)
		}
	},
}

var addContainer = &cobra.Command{
	Use:   "add-container",
	Short: "Add a container to a project",
	Long:  `Add a container to a project.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		containerID, _ := cmd.Flags().GetString("container")

		driver := getDriver()

		database, err := db.NewDatabase(driver)

		if err != nil {
			println(err.Error())
			return
		}

		project, err := database.GetProjectByName(projectName)

		if err != nil {
			println(err.Error())
			return
		}

		docker, err := dockerClient.NewDockerClient()
		if err != nil {
			println(err.Error())
			return
		}

		dockerContainer, err := docker.GetContainerByID(containerID)
		if err != nil {
			println(err.Error())
			return
		}

		container, err := database.GetOrAddContainer(dockerContainer.ID, dockerContainer.Name)

		if err != nil {
			println(err.Error())
			return
		}

		err = project.AddContainer(container.ID)

		if err != nil {
			if database.IsUniqueViolationError(err) {
				println("Container already added to project")
				return
			}

			println(err.Error())
			return
		}

		println("Container " + container.Name + " (" + container.ID + ") added to project " + project.Name)
	},
}

func init() {
	initProject.Flags().StringP("name", "n", "", "Name of the project")
	initProject.MarkFlagRequired("name") // Make the name flag required

	addContainer.Flags().StringP("project", "p", "", "Name of the project")
	addContainer.Flags().StringP("container", "c", "", "ID of the container")
	addContainer.MarkFlagRequired("project")
	addContainer.MarkFlagRequired("container")

	projectsCmd.AddCommand(initProject)
	projectsCmd.AddCommand(addContainer)
	projectsCmd.AddCommand(listProjects)
}

type Project struct {
	Name string
}
