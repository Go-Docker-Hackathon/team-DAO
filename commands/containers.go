package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"../utils"

	"github.com/fsouza/go-dockerclient"
)

// Now volrep only supports running locally.
var Endpoint = "unix:///var/run/docker.sock"

// Now volrep only supports using /var/lib/docker/vfs/dir as a default path.
// FIXME: make it compiltile with all kinds of situations wherever the root dir path is.
var VolumePath = "/var/lib/docker/vfs/dir"

const storage_path = "/root/.volrep"

func GetContainer(name string) (*docker.Container, error) {
	// 1.check the source container:
	//   1.1 whether source is legal.
	//   1.2 whether it matches a container in Docker Daemon.
	//   1.3 whether it has a data volume.
	if err := utils.ValidateName(name); err != nil {
		fmt.Println("Name of the source container is illegal.")
		os.Exit(1)
	}

	// Construct a docker client
	client, _ := docker.NewClient(Endpoint)

	container, err := isContainerExist(name, client)

	if err != nil {
		fmt.Errorf("error !")
		os.Exit(1)
	}

	containerID := container.ID
	containerImage := container.Config.Image
	//containerVolume := container.Volumes
	containerName := container.Name

	fmt.Println("Below is container's details.")
	fmt.Println("Container ID :      " + containerID)
	fmt.Println("Container Name :    " + containerName)
	fmt.Println("Container Image :   " + containerImage)
	//fmt.Print("Container Volumes : ")
	//fmt.Println(containerVolume)

	return container, nil
}

func StartContainer(name string) error {
	client, _ := docker.NewClient(Endpoint)
	err := client.StartContainer(name, nil)
	return err
}

func StopContainer(name string, timeout uint) error {
	client, _ := docker.NewClient(Endpoint)
	err := client.StopContainer(name, timeout)
	return err
}

// If the volumes has the prefix like '/var/lib/docker/vfs/dir',
// We can regard it as a data volume.
func GetContainerDataVolumes(sourceCon *docker.Container) (map[string]string, error) {
	sourceDataVolumes := map[string]string{}
	for k, v := range sourceCon.Volumes {
		if strings.HasPrefix(v, VolumePath) {
			sourceDataVolumes[k] = v
		}
	}
	if len(sourceDataVolumes) == 0 {
		return sourceDataVolumes, fmt.Errorf("The container has no data volume. \nAborting !")
	}
	return sourceDataVolumes, nil
}

func GetAllCompressedVolumes(sourceCon *docker.Container) (map[string][]string, error) {
	allCompressedVolumes := map[string][]string{}

	// If we use go-dockerclient, container.Name will be like "/iamcontainername"
	containerNameStr := strings.Replace(sourceCon.Name, "/", "_", -1)
	// Combine ID and Name together to identify
	containerPath := sourceCon.ID + containerNameStr
	// constrcut container's absolute path
	CompressedVolumesPath := path.Join(storage_path, containerPath)

	//fmt.Println("Container compressed path: " + CompressedVolumesPath)

	files, err := ioutil.ReadDir(CompressedVolumesPath)
	if err != nil {
		fmt.Println("Error when reading dirs in " + CompressedVolumesPath)
		os.Exit(1)
	}

	for _, file := range files {
		fileNameStr := file.Name()
		// Only in containerVolumePathStr of compress.go, we use "-"
		index := strings.Index(fileNameStr, "-")
		if index == -1 {
			fmt.Println("Got an error when get index of '-' in compressed file name.")
			continue
		}
		timeStr := fileNameStr[:index]
		//fmt.Println("TimeStr is: ", timeStr)
		// It seems code below does not work.
		// slice append?
		allCompressedVolumes[timeStr] = append(allCompressedVolumes[timeStr], fileNameStr)
	}
	return allCompressedVolumes, nil

}

func isContainerExist(name string, client *docker.Client) (*docker.Container, error) {
	container, err := client.InspectContainer(name)
	return container, err
}
