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
	//fmt.Println("Parameters: ", index, len(timeStrArray))
	if index == len(timeStrArray)+1 {
		// if we need to remove all
		// the index user input is len(timeStrArray)+1
		fmt.Println(path.Join(storage_path, containerFileName))
		err := os.RemoveAll(path.Join(storage_path, containerFileName))
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

func LoadVolumesForTargetContainer(container *docker.Container, volumes []string, targetDataVolumes map[string]string) error {

	containerIdStr := container.ID
	containerNameStr := strings.Replace(container.Name, "/", "_", -1)
	containerFileName := containerIdStr + containerNameStr
	//fmt.Println("volumes:", volumes)
	//fmt.Println("target Volumes:", targetDataVolumes)
	for _, volumeFilename := range volumes {
		//fmt.Println("VolumeFilename:", volumeFilename)

		// here volumeFilename is something like this :
		// 2015_06_05_22_41_08-var-lib-mysql.tar
		// we first get the first "-" and "."
		// Then replace "-" with "/".
		// result is volume view of container itself.
		first := strings.Index(volumeFilename, "-")
		last := strings.LastIndex(volumeFilename, ".")
		volumePath := volumeFilename[first:last]
		volumePath = strings.Replace(volumePath, "-", "/", -1)

		//fmt.Println("volumePath: " + volumePath)

		for key, value := range targetDataVolumes {
			//fmt.Println("key:", key, "volumePath:", volumePath)
			if volumePath == key {
				// We got the matched compressed data volume and
				// actual data volume path of target container
				// 1.delete data volume of target contaier.
				// 2.untar compressed data volume to actual data volume path of container

				//fmt.Println("Remove the path of " + value)
				// it will remove everything including the file name
				// So we need to mkdir value first when executing `tar -xf `
				err := os.RemoveAll(value)
				if err != nil {
					fmt.Println("Failed to remove all details in data volumes.")
					os.Exit(1)
				}

				err = os.MkdirAll(value, 744)
				if err != nil {
					fmt.Println("Failed to mkdir -p " + value)
					os.Exit(1)
				}

				dataVolumeAbsPath := path.Join(storage_path, containerFileName, volumeFilename)
				cmd := "cd " + value + "&& tar -xf " + dataVolumeAbsPath

				//fmt.Println("cmd is: " + cmd)

				_, err = ExecShell(cmd)

				if err != nil {
					fmt.Println("Command tar ran into an Error.")
					return err
				}
			}
		}
	}

	return nil
}
