package commands

import (
	"fmt"
	"os"
	"strconv"

	"../utils"

	"github.com/codegangsta/cli"
)

func CmdRm(c *cli.Context) {
	isAll := c.Bool("all")

	if isAll == false && len(c.Args()) != 1 {
		fmt.Println("Command `volrep load` needs exact 1 arguments. Please check again.")
		os.Exit(1)
	}

	source := c.Args()[0]

	fmt.Println("The source container is " + source)

	sourceCon, err := GetContainer(source)
	if err != nil {
		fmt.Println("Can not get docker container corresponding the name you provided.")
		fmt.Println("Please check your specified source container name again.")
		os.Exit(1)
	}

	allCompressedVolumes, err := GetAllCompressedVolumes(sourceCon)
	if err != nil {
		fmt.Println("Error when getting all compressed volumes.")
		os.Exit(1)
	}

	var i int = 1
	for key, value := range allCompressedVolumes {
		fmt.Println(strconv.Itoa(i) + ":" + key)
		for _, item := range value {
			fmt.Println("    " + item)
		}
		i++
	}

	var index int

	if isAll == true {
		index = -1
	} else {
		index = utils.AskForNumberPick(len(allCompressedVolumes))
		if index == -1 {
			fmt.Println("The number you pick is invalid, please check it again.")
			os.Exit(1)
		}
	}

	fmt.Println("Prepare to rm the specified compressed data volumes.")

	err = utils.RemoveCompressedVolumes(sourceCon.ID, index-1)
	if err != nil {
		fmt.Println("Something happened when removing compressed volumes, error:", err)
		os.Exit(1)
	}

	fmt.Println("Successfully remove compressed volumes you specified.")
}
