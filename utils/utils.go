package utils

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// FIXME
// By dafault, we assume the name is legal.
func ValidateName(name string) error {
	return nil
}

// Duplicate data volumes from source container to target container.
func DuplicateDataVolumes(sVols map[string]string, tVols map[string]string) error {
	for sInPath, sOutPath := range sVols {
		// Traverse source container's data volumes
		// And get the corresponding target container's data volume path
		if tOutPath, ok := tVols[sInPath]; ok {
			// Start to duplicate source data volume's detais to the target
			if err := copyVolume(sOutPath, tOutPath); err != nil {
				fmt.Println("Got an error when duplicating data volume %s to %s", sOutPath, tOutPath)
				return err
			}
		}
	}
	return nil
}

func RemoveCompressedVolumes(container *docker.Container, timeStrArray []string, index int) error {
	containerIdStr := container.ID
	containerNameStr := strings.Replace(container.Name, "/", "_", -1)
	containerFileName := containerIdStr + containerNameStr
	fmt.Println("Parameters: ", index, len(timeStrArray))
	if index == len(timeStrArray)+1 {
		// if we need to remove all
		// the index user input is len(timeStrArray)+1
		fmt.Println(path.Join(storage_path, containerFileName))
		err := os.Remove(path.Join(storage_path, containerFileName))
		if err != nil {
			fmt.Println("Failed to remove the container directory.")
			os.Exit(1)
		}
		return nil
	} else if index >= 1 && index <= len(timeStrArray) {
		containerPath := path.Join(storage_path, containerFileName)
		files, err := ioutil.ReadDir(containerPath)
		if err != nil {
			fmt.Println("Error when get all dirs in container path.")
			os.Exit(1)
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), timeStrArray[index-1]) {
				fmt.Println("This is a corresponding compressed data volume. Remove it.")
				err := os.Remove(path.Join(storage_path, containerFileName, file.Name()))
				if err != nil {
					fmt.Println("Error in removing one specified compressed data volume.")
					os.Exit(1)
				}
			}
		}
		return nil
	} else {
		fmt.Println("Invalid input index.")
		os.Exit(1)
	}

	return nil
}
