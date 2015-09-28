package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	// https://github.com/codegangsta/cli
	app := cli.NewApp()
	app.Name = "klouds-agent"
	app.Usage = "lets docker and it's containers know what's up. "
	app.Commands = []cli.Command{
		{
			Name:    "start-api",
			Aliases: []string{"api"},
			Usage:   "http interface",
			Action: func(c *cli.Context) {
				Serve()
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list running containers",
			Action: func(c *cli.Context) {
				docker := NewDockerClient()
				List(docker)
			},
		},
		{
			Name:  "inspect",
			Usage: "inspect a container",
			Action: func(c *cli.Context) {
				containerId := c.Args().First()
				docker := NewDockerClient()
				Inspect(docker, containerId)
			},
		},
		{
			Name:  "start",
			Usage: "start a container",
			Action: func(c *cli.Context) {
				containerId := c.Args().First()
				docker := NewDockerClient()
				Start(docker, containerId)
			},
		},
		{
			Name:  "stop",
			Usage: "stop a container",
			Action: func(c *cli.Context) {
				containerId := c.Args().First()
				docker := NewDockerClient()
				Stop(docker, containerId)
			},
		},
		{
			Name:  "create",
			Usage: "create a container",
			Action: func(c *cli.Context) {
				name := c.Args().First()
				image := c.Args()[1]
				docker := NewDockerClient()
				Create(docker, name, image)
			},
		},
		{
			Name:  "build",
			Usage: "build a container",
			Action: func(c *cli.Context) {
				repoName := c.Args().First()
				context := c.Args()[1]
				docker := NewDockerClient()
				Build(docker, repoName, context)
			},
		},
	}
	app.Run(os.Args)
}
