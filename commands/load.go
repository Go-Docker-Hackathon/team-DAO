package commands

import (
	"fmt"
	"os"
	"strconv"

	"../utils"

	"github.com/codegangsta/cli"
)

func CmdLoad(c *cli.Context) {
	if len(c.Args()) != 2 {
		fmt.Println("Command `volrep load` needs exact 2 arguments. Please check again.")
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

	fmt.Println("\nThe target container is " + target)
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
	_, err = GetContainerDataVolumes(sourceCon)
	if err != nil {
		fmt.Println("The source container has no data volume!\n Aborting!")
		os.Exit(1)
	}

	targetDataVolumes, err := GetContainerDataVolumes(targetCon)
	if err != nil {
		fmt.Println("The target container has no data volume!\n Aborting!")
		os.Exit(1)
	}

	//fmt.Println(sourceDataVolumes)
	//fmt.Println(targetDataVolumes)

	// If their images have same repo name and same tag name.
	// Check the state of source container.
	sourceState := sourceCon.State.Running
	if sourceState == true {
		fmt.Println("The source container is running and data volumes duplication will stop the running source container.")
		fmt.Print("Attention: Stop the source container <" + sourceCon.Name + "> ? (yes/no) ")

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

	allCompressedVolumes, err := GetAllCompressedVolumes(sourceCon)

	//fmt.Println("The source container <" + sourceCon.Name + "> has " + strconv.Itoa(len(allCompressedVolumes)) + " volumes.")

	if err != nil {
		fmt.Println("Error when getting all compressed volumes.")
		os.Exit(1)
	}

	VolumesNum := len(allCompressedVolumes)

	fmt.Println("Here is " + strconv.Itoa(VolumesNum) + " replicated volumes, like below:")

	timeStrArray := []string{}
	var i int = 1
	for key, value := range allCompressedVolumes {
		fmt.Println(strconv.Itoa(i) + ".Replicate Time: " + key)
		timeStrArray = append(timeStrArray, key)
		for _, item := range value {
			fmt.Println("    Compressed file name: " + item)
		}
		i++
	}

	fmt.Print("Please choose which one to load: ")

	pickedNum := utils.AskForNumberPick(i - 1)
	if pickedNum == -1 {
		fmt.Println("Invalid input index.")
		os.Exit(1)
	}

	sourceDataVolumesArray := allCompressedVolumes[timeStrArray[pickedNum-1]]
	//fmt.Println(sourceDataVolumesArray)

	targetState := targetCon.State.Running
	if targetState == true {
		fmt.Println("The target container is running and data volumes duplication will stop the running source container.")
		fmt.Print("Attention: Stop the target container <" + targetCon.Name + "> ? (yes/no) ")

		// Ask for user's confirmation of stopping running container.
		isConfirmed := utils.AskForConfirmation()
		if isConfirmed != true {
			fmt.Println("Stopping target container is not confirmed. Aborting...")
			os.Exit(1)
		}

		// Container is running, there are two policies: stop or pause container.
		// 1.Try to stop container.
		fmt.Print("Confirmed. \n1.Stop container...   ")

		if err := StopContainer(targetCon.ID, 10); err != nil {
			fmt.Println("Got an error when stopping a running container. Error is %s", err)
			os.Exit(1)
		}
		fmt.Println("OK")
	}

	// Start to duplicate the specific volumes from source to target.
	fmt.Println("2.Start to load data volumes from source container...   ")
	err = utils.LoadVolumesForTargetContainer(sourceCon, sourceDataVolumesArray, targetDataVolumes)
	if err != nil {
		fmt.Println("Failed to load compressed data volumes for target containers.")
		os.Exit(1)
	} else {
		fmt.Println("Loading compressed data volumes finished successfully.")
		fmt.Println("OK")
		//fmt.Println("Start source container and target container.")
	}

	fmt.Print("3.Start source container...   ")
	err = StartContainer(sourceCon.ID)
	if err != nil {
		fmt.Println("Fail to start the source container.")
		os.Exit(1)
	}
	fmt.Println("OK")

	fmt.Print("4.Start target container...   ")
	err = StartContainer(targetCon.ID)
	if err != nil {
		fmt.Println("Fail to start the target container.")
		os.Exit(1)
	}
	fmt.Println("OK")

	fmt.Println("Well Done!")
}
