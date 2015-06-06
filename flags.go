package main

import (
	"github.com/codegangsta/cli"
)

var (
	flDockerRoot = cli.StringFlag{
		Name:  "dockerdir",
		Value: "/var/lib/docker",
		Usage: "It is the root path of docker.",
	}
	flRemoveAll = cli.BoolFlag{
		Name: "all",
		//Value: false,
		Usage: "It represents wether to remove all compressed packages.",
	}
)
