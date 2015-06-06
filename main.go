package main

import (
	"fmt"
	"os"
	"path"

	"./commands"

	"github.com/codegangsta/cli"
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
	app.Email = " allen.sun@daocloud.io "

	app.Commands = []cli.Command{
		{
			Name:        "dup",
			Usage:       "Duplicate a container's data volumes to another one.",
			Description: "FORMAT : volrep dup srcCon destCon.",
			Flags: []cli.Flag{
				flDockerRoot,
			},
			Action: commands.CmdDup,
		},
		{
			Name:        "save",
			Usage:       "Save a container's data volumes as compressed volumes.",
			Description: "FORMAT: volrep save containerA.",
			Action:      commands.CmdSave,
		},
		{
			Name:        "load",
			Usage:       "Load a container's data volumes from sepcified container's compressed volumes.",
			Description: "FORMAT: volrep load srcCon destCon.",
			Action:      commands.CmdLoad,
		},
		{
			Name:        "list",
			Usage:       "List a container's all compressed volume.",
			Description: "FORMAT: volrep list srcCon.",
			Action:      commands.CmdList,
		},
		{
			Name:        "rm",
			Usage:       "Remove a container's compressed volume.",
			Description: "FORMAT: volrep rm srcCon.",
			Flags: []cli.Flag{
				flRemoveAll,
			},
			Action: commands.CmdRm,
		},
	}

	app.Run(os.Args)
}
