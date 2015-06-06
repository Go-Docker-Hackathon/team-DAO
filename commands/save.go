package commands

import (
	"fmt"
	"os"

	"../utils"

	"github.com/codegangsta/cli"
)

func CmdSave(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println("Command `volrep save` needs exact 1 arguments. Please check again.")
		os.Exit(1)
	}

	target := c.Args()[0]

	fmt.Println("Container <" + target + "> is the one whose volumes need compressed.")

	tarCon, err := GetContainer(target)
	if err != nil {
		fmt.Println("Container " + tarCon.Name + " does not exist.")
		fmt.Println("Please check container name again.")
		os.Exit(1)
	}

	// Check if source container has a data volume.
	// Only select data volume of source container into sourceDataVolumes.
	// A container can have serveral data volumes.
	dataVolumes, err := GetContainerDataVolumes(tarCon)
	//fmt.Println(dataVolumes)

	if err != nil {
		fmt.Println("The source container has no data volume!")
		fmt.Println("Nothing to compress. ")
		fmt.Println("Aborting!")
		os.Exit(1)
	}

	// If their images have same repo name and same tag name.
	// Check the state of source container.
	conState := tarCon.State.Running
	if conState == true {
		fmt.Println("The container is running and saving data volumes will stop the running source container.")
		fmt.Print("Attention: Stop container <" + tarCon.Name + "> ? (yes/no) ")

		// Ask for user's confirmation of stopping running container.
		isConfirmed := utils.AskForConfirmation()
		if isConfirmed != true {
			fmt.Println("Stopping container is not confirmed.")
			fmt.Println("Aborting...")
			os.Exit(1)
		}

		// Container is running, there are two policies: stop or pause container.
		// 1.Try to stop container.
		fmt.Println("Confirmed.")
		fmt.Print("1.Stop container...   ")
		if err := StopContainer(tarCon.ID, 10); err != nil {
			fmt.Println("Got an error when stopping a running container. Error is %s", err)
			os.Exit(1)
		}
		fmt.Println("OK")
	}

	// Start to save all data volumes into tars.
	fmt.Print("2.Start to compress data volumes...   ")
	err = utils.CompressDataVolumes(tarCon, dataVolumes)
	if err != nil {
		fmt.Println("Got an error when compressing the data volumes. Error is %s", err)
	}
	fmt.Println("OK")

	fmt.Print("3.Start container...   ")
	err = StartContainer(tarCon.ID)
	if err != nil {
		fmt.Println("Fail to start the container.")
		os.Exit(1)
	}
	fmt.Println("OK")

	fmt.Println("Well Done!")

}
