package cli

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"fmt"

	"github.com/spf13/cobra"
)

var keyfile string
var passphrase string

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  `Manage projects`,
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long:  `List projects`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client := getDbClient()
		projects, err := client.GetProjects()
		if err != nil {
			fmt.Printf("Could not get projects: %s\n", err.Error())
		}

		for _, project := range projects {
			fmt.Printf("Project: %s\n", project.GetName())
		}

		return
	},
}

var addProjectCmd = &cobra.Command{
	Use:   "add",
	Short: "Add project",
	Long:  `Add project`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if keyfile == "" {
			return
		}

		fmt.Print("Adding with keyfile: ", keyfile, "\n")

		if passphrase == "" {
			return
		}

		fmt.Print("Adding with passphrase: ", passphrase, "\n")

		client := getDbClient()
		project, err := client.AddProject(args[0])
		if err != nil {
			fmt.Printf("Could not add project: %s\n", err.Error())
			return
		}

		fmt.Printf("Added project: %s\n", project.GetName())
	},
}

func AddContainerToProject(projectName string, containerName string) error {
	db := getDbClient()

	project, err := db.GetProjectByName(projectName)
	if err != nil {
		return err
	}

	fmt.Printf("Got project: %s\n", project.GetName())

	docker := getDockerClient()
	dockerContainer, err := docker.GetContainer(containerName)
	if err != nil {
		return err
	}

	fmt.Printf("Got container: %s\n", dockerContainer.GetName())

	dbContainer, err := db.GetOrAddContainer(dockerContainer.GetID(), dockerContainer.GetName())
	if err != nil {
		return err
	}

	err = project.AddContainer(dbContainer)
	if err != nil {
		if _, ok := err.(*errors.ContainerAlreadyInProjectError); ok {
			fmt.Printf("Container already in project\n")
		} else {
			return err
		}
	}

	var errs []error
	for _, dockerVolume := range dockerContainer.GetVolumes() {
		dbVolume, err := db.GetOrAddVolume(dockerVolume.GetName())
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = dbContainer.AddVolume(dbVolume)
		if err != nil {
			if _, ok := err.(*errors.VolumeAlreadyInContainerError); ok {
				fmt.Printf("Volume already in container\n")
			} else {
				errs = append(errs, err)
				continue
			}
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

var addContainerToProjectCmd = &cobra.Command{
	Use:   "add-container",
	Short: "Add container to project",
	Long:  `Add container to project`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		containerName := args[1]

		err := AddContainerToProject(projectName, containerName)
		if err != nil {
			fmt.Printf("Could not add container to project: %s\n", err.Error())
			return
		}
	},
}

var listContainersCmd = &cobra.Command{
	Use:   "list-containers",
	Short: "List containers",
	Long:  `List containers`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		db := getDbClient()

		project, err := db.GetProjectByName(projectName)
		if err != nil {
			fmt.Printf("Could not get project: %s\n", err.Error())
			return
		}

		containers, err := project.GetContainers()
		if err != nil {
			fmt.Printf("Could not get containers: %s\n", err.Error())
			return
		}

		for _, container := range containers {
			fmt.Printf("Container: %s\n", container.GetName())
		}

		return
	},
}

func BackupProjectCmd(projectName string) error {
	db := getDbClient()
	docker := getDockerClient()

	project, err := db.GetProjectByName(projectName)
	if err != nil {
		return err
	}

	containers, err := project.GetContainers()
	if err != nil {
		return err
	}

	var dockerContainers []interfaces.DockerContainer
	var errs []error
	for _, container := range containers {
		dockerContainer, err := docker.GetContainer(container.GetID())
		if err != nil {
			errs = append(errs, err)
			continue
		}

		dockerContainers = append(dockerContainers, dockerContainer)
	}

	volumesMap := make(map[string]interfaces.DockerVolume)
	for _, dockerContainer := range dockerContainers {
		for _, dockerVolume := range dockerContainer.GetVolumes() {
			volumesMap[dockerVolume.GetName()] = dockerVolume
		}
	}

	var volumes []interfaces.DockerVolume
	for _, dockerVolume := range volumesMap {
		dockerVolume.SetMountPoint("/input/" + dockerVolume.GetName())
		volumes = append(volumes, dockerVolume)
	}

	// Pre-Backup (ends with docker stop)
	errs = []error{}
	for _, dockerContainer := range dockerContainers {
		err = dockerContainer.Stop()
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	// Backup (ends with docker start)
	worker, err := docker.CreateContainer("worker", volumes)
	if err != nil {
		return err
	}

	err = worker.Start()
	if err != nil {
		return err
	}

	err = worker.StopAndRemove()
	if err != nil {
		return err
	}

	errs = []error{}
	// Post-Backup (starts with docker start)
	for _, dockerContainer := range dockerContainers {
		err = dockerContainer.Start()
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

var backupProjectCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup project",
	Long:  `Backup project`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		err := BackupProjectCmd(projectName)
		if err != nil {
			fmt.Printf("Could not backup project: %s\n", err.Error())
			return
		}
	},
}

func init() {
	projectsCmd.AddCommand(listProjectsCmd)
	projectsCmd.AddCommand(addProjectCmd)
	projectsCmd.AddCommand(addContainerToProjectCmd)
	projectsCmd.AddCommand(listContainersCmd)
	projectsCmd.AddCommand(backupProjectCmd)

	addProjectCmd.Flags().StringVarP(&keyfile, "keyfile", "k", "", "Path to the keyfile")
	addProjectCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase for the project")
}
