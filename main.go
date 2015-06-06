package main

import (
	"fmt"
	"os"
	"path"

	"./commands"

	"github.com/codegangsta/cli"
)

var (
	flDockerRoot = cli.StringFlag{
		Name:  "dockerdir",
		Value: "/var/lib/docker",
		Usage: "It is the root path of docker.",
	}
	flRemoveAll = cli.BoolFlag{
		Name:  "all",
		Value: false,
		Usage: "It represents wether to remove all compressed packages.",
	}
)

func main() {
	arg_num := len(os.Args)

	// If the command parameter number indeed smaller than 2 ??
	if arg_num <= 1 {
		fmt.Errorf("The num of input args is %s, too short.\n", arg_num-1)
		return
	}

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "Allen Sun"
	app.Email = "allen.sun@daocloud.io"

	app.Commands = []cli.Command{
		{
			Name:        "dup",
			Usage:       "Duplicate a container's data volumes to another container.",
			Description: "FORMAT : volrep dup srcCon destCon.",
			Flags: []cli.Flag{
				flDockerRoot,
			},
			Action: cmdDup,
		},
		{
			Name:        "save",
			Usage:       "save a container's data volumes as compressed volumes.",
			Description: "The save command is used like this, volrep save containerA.",
			Action:      cmdSave,
		},
		{
			Name:        "load",
			Usage:       "load a container's data volumes as compressed volumes.",
			Description: "FORMAT: volrep load srcCon destCon.",
			Action:      cmdLoad,
		},
		{
			Name:        "list",
			Usage:       "list a container's all compressed volume.",
			Description: "FORMAT: volrep list srcCon.",
			Action:      cmdList,
		},
		{
			Name:        "rm",
			Usage:       "rm a container's compressed volume.",
			Description: "FORMAT: volrep rm srcCon.",
			Flags: []cli.Flag{
				flAll,
			},
			Action: cmdRm,
		},
	}

	app.Run(os.Args)
}
