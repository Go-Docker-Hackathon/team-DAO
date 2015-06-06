package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
)

const storage_path = "/root/.volrep"
const layout = "2006-01-02 15:04:05"

func CompressDataVolumes(container *docker.Container, volumes map[string]string) error {
	// Check if the storage path exists.
	_, err := os.Stat(storage_path)
	if err != nil {
		os.Mkdir(storage_path, 0)
	}

	containerNameStr := strings.Replace(container.Name, "/", "_", -1)
	containerPath := container.ID + containerNameStr

	if _, err := os.Stat(path.Join(storage_path, containerPath)); err != nil {
		fmt.Println(path.Join(storage_path, containerPath) + " does not exist yet.")
		fmt.Println("Try to make this directory ...")

		err := os.MkdirAll(path.Join(storage_path, containerPath), 700)
		if err != nil {
			fmt.Println("Mkdir Error: ", path.Join(storage_path, container.ID+"_"+containerNameStr))
			os.Exit(1)
		}
	}

	for key, value := range volumes {
		// Convert current time into string
		// in timeStr we use "_"
		timeStr := time.Now().Format(layout)
		timeStr = strings.Replace(timeStr, "-", "_", -1)
		timeStr = strings.Replace(timeStr, " ", "_", -1)
		timeStr = strings.Replace(timeStr, ":", "_", -1)

		// unlike timeStr, here in containerVolumePathStr we use "-"
		// it will be used in GetAllCompressedVolumes.
		containerVolumePathStr := strings.Replace(key, "/", "-", -1)
		archiveName := timeStr + containerVolumePathStr + ".tar"

		destPathName := path.Join(storage_path, containerPath, archiveName)

		//fmt.Println(archiveName)
		//fmt.Println(destPathName)

		// Enter volume path and compress all thing into ~/.volrep
		// whether to mkdir -p first?
		cmd := "cd " + value + " && tar -zcvf " + destPathName + " ."

		_, err = ExecShell(cmd)

		if err != nil {
			fmt.Println("Command tar ran into an Error.")
			return err
		}
	}
	return nil
}

func ExecShell(cmd string) (string, error) {
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	result := string(out[:len(out)])
	return result, nil
}
