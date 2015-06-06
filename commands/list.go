package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	//"../utils"

	"github.com/codegangsta/cli"
)

func CmdList(c *cli.Context) {

	if len(c.Args()) != 1 {
		fmt.Println("Command `volrep load` needs exact 1 arguments. Please check again.")
		os.Exit(1)
	}

	source := c.Args()[0]

	fmt.Println("The source container is " + source + ".")

	sourceCon, err := GetContainer(source)
	if err != nil {
		fmt.Println("Can not get docker container corresponding the name you provided.")
		fmt.Println("Please check your specified source container name again.")
		os.Exit(1)
	}

	allCompressedVolumes, err := GetAllCompressedVolumes(sourceCon)

	if err != nil {
		fmt.Println("Error when getting all compressed volumes.", allCompressedVolumes)
		os.Exit(1)
	}
	fmt.Println("Container <" + source + "> has " + strconv.Itoa(len(allCompressedVolumes)) + " compressed data volumes.")

	i := 1
	for _, value := range allCompressedVolumes {
		index := strings.Index(value[0], "-")
		if index == -1 {
			fmt.Println("Got an error when get index of '-' in compressed file name.")
			continue
		}
		timeStr := value[0][:index]

		fmt.Println(strconv.Itoa(i) + ". Replication time: " + strings.Replace(timeStr[:10], "_", "-", -1) + " " + strings.Replace(timeStr[11:], "_", ":", -1))
		for _, item := range value {
			fmt.Println("   " + item)
		}
		i++
	}
}
