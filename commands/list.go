package commands

import (
	"fmt"
	"os"

	//"../utils"

	"github.com/codegangsta/cli"
)

func CmdList(c *cli.Context) {

	if len(c.Args()) != 1 {
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
		fmt.Println("Error when getting all compressed volumes.", allCompressedVolumes)
		os.Exit(1)
	}

}
