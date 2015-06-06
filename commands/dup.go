package commands

import (
	"fmt"
	"os"

	"../utils"

	"github.com/codegangsta/cli"
)

func CmdDup(c *cli.Context) {
	if len(c.Args()) != 2 {
		fmt.Println("Command `volrep dup` needs exact 2 arguments. Please check again.")
		os.Exit(1)
	}

	source := c.Args()[0]
	target := c.Args()[1]

	fmt.Println("The source container is " + source)

	sourceCon, err := GetContainer(source)
	if err != nil {
		fmt.Println("Can not get docker container corresponding the name you provided.")
		fmt.Println("Please check your specified source container name again.")
		os.Exit(1)
	}

	fmt.Println("The target container is " + target)
	targetCon, err := GetContainer(target)
	if err != nil {
		fmt.Println("Can not get docker container corresponding the name you provided.")
		fmt.Println("Please check your specified target container name again.")
		os.Exit(1)
	}

	// Check if the two containers' images are the same.
	if sourceCon.Config.Image != targetCon.Config.Image {
		fmt.Println("The images of both two container is not the same. Action of duplicating aborted.")
		os.Exit(1)
	}

	// Check if source container has a data volume.
	// Only select data volume of source container into sourceDataVolumes.
	// As the source and target containers use the same image.
	// There is no need to check the target container's data volume.
	// But getting the data volumes is necessary.
	sourceDataVolumes, err := GetContainerDataVolumes(sourceCon)
	if err != nil {
		fmt.Println("The source container has no data volume!\n Aborting!")
		os.Exit(1)
	}

	targetDataVolumes, err := GetContainerDataVolumes(targetCon)
	if err != nil {
		fmt.Println("The target container has no data volume!\n Aborting!")
		os.Exit(1)
	}

	// If their images have same repo name and same tag name.
	// Check the state of source container.
	sourceState := sourceCon.State.Running
	if sourceState == true {
		fmt.Println("The source container is running and data volumes duplication will stop the running source container.")
		fmt.Println("\nAttention: Stop the source container <" + sourceCon.Name + "> ? (yes/no) ")

		// Ask for user's confirmation of stopping running container.
		isConfirmed := utils.AskForConfirmation()
		if isConfirmed != true {
			fmt.Println("Stopping source container is not confirmed. Aborting...")
			os.Exit(1)
		}

		// Container is running, there are two policies: stop or pause container.
		// 1.Try to stop container.
		if err := StopContainer(sourceCon.ID, 10); err != nil {
			fmt.Println("Got an error when stopping a running container. Error is %s", err)
			os.Exit(1)
		}
	}

	targetState := targetCon.State.Running
	if targetState == true {
		fmt.Println("The target container is running and data volumes duplication will stop the running source container.")
		fmt.Println("\nAttention: Stop the target container <" + targetCon.Name + "> ? (yes/no) ")

		// Ask for user's confirmation of stopping running container.
		isConfirmed := utils.AskForConfirmation()
		if isConfirmed != true {
			fmt.Println("Stopping target container is not confirmed. Aborting...")
			os.Exit(1)
		}

		// Container is running, there are two policies: stop or pause container.
		// 1.Try to stop container.
		if err := StopContainer(targetCon.ID, 10); err != nil {
			fmt.Println("Got an error when stopping a running container. Error is %s", err)
			os.Exit(1)
		}
	}

	// Start to duplicate all data volumes from source to target.
	err = utils.DuplicateDataVolumes(sourceDataVolumes, targetDataVolumes)
	if err != nil {
		fmt.Println("Got an error when duplicating the data volumes. Error is %s", err)
	}

	err = StartContainer(sourceCon.ID)
	if err != nil {
		fmt.Println("Fail to start the source container.")
		os.Exit(1)
	}

	err = StartContainer(targetCon.ID)
	if err != nil {
		fmt.Println("Fail to start the target container.")
		os.Exit(1)
	}

	fmt.Println("Well Done!")

}
