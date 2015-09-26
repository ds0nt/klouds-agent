package main

import (
  "github.com/codegangsta/cli"
  "github.com/samalba/dockerclient"
  "log"
  "time"
  "os"
)

func Client() dockerclient.Client {
  docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
  return docker
}

// Callback used to listen to Docker's events
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
    log.Printf("Received event: %#v\n", *event)
}

func MonitorEvents(docker dockerclient.Client) {
  docker.StartMonitorEvents(eventCallback, nil)
  time.Sleep(3600 * time.Second)
}

func List(docker dockerclient.Client)  {
  // Get only running containers
  containers, err := docker.ListContainers(true, false, "")
  if err != nil {
      log.Fatal(err)
  }
  for _, c := range containers {
      log.Println(c.Id)
  }
}

func Inspect(docker dockerclient.Client, containerId string) {
  info, _ := docker.InspectContainer(containerId)
  log.Println(info)
}

func Start(docker dockerclient.Client, containerId string)  {
  // Start the container
  hostConfig := &dockerclient.HostConfig{}
  err := docker.StartContainer(containerId, hostConfig)
  if err != nil {
      log.Fatal(err)
  }
}

func Stop(docker dockerclient.Client, containerId string) {
  // Stop the container (with 5 seconds timeout)
  docker.StopContainer(containerId, 5)
}

  //
  // func Build(docker dockerclient.Client) {
  //   // Build a docker image
  //   // some.tar contains the build context (Dockerfile any any files it needs to add/copy)
  //   dockerBuildContext, err := os.Open("some.tar")
  //   defer dockerBuildContext.Close()
  //   buildImageConfig := &dockerclient.BuildImage{
  //           Context:        dockerBuildContext,
  //           RepoName:       "your_image_name",
  //           SuppressOutput: false,
  //   }
  //   reader, err := docker.BuildImage(buildImageConfig)
  //   if err != nil {
  //       log.Fatal(err)
  //   }
  // }

func Create(docker dockerclient.Client, name string, image string)  {
  // Create a container
  containerConfig := &dockerclient.ContainerConfig{
      Image: image,
      // Cmd:   []string{"bash"},
      AttachStdin: true,
      Tty:   true}
  containerId, err := docker.CreateContainer(containerConfig, name)
  if err != nil {
      log.Fatal(err)
  }
  log.Println(containerId)
}


func main() {
  // https://github.com/codegangsta/cli
  app := cli.NewApp()
  app.Name = "klouds-agent"
  app.Usage = "lets docker and it's containers know what's up. "
  app.Commands = []cli.Command{
    {
      Name:      "list",
      Aliases:     []string{"ls"},
      Usage:     "list running containers",
      Action: func(c *cli.Context) {
        docker := Client()
        List(docker)
      },
    },
    {
      Name:      "inspect",
      Usage:     "inspect a container",
      Action: func(c *cli.Context) {
        containerId := c.Args().First()
        docker := Client()
        Inspect(docker, containerId)
      },
    },
    {
      Name:      "start",
      Usage:     "start a container",
      Action: func(c *cli.Context) {
        containerId := c.Args().First()
        docker := Client()
        Start(docker, containerId)
      },
    },
    {
      Name:      "stop",
      Usage:     "stop a container",
      Action: func(c *cli.Context) {
        containerId := c.Args().First()
        docker := Client()
        Stop(docker, containerId)
      },
    },
    {
      Name:      "create",
      Usage:     "create a container",
      Action: func(c *cli.Context) {
        name := c.Args()[0]
        image := c.Args()[1]
        docker := Client()
        Create(docker, name, image)
      },
    },
      //
      // {
      //   Name:      "list",
      //   Aliases:     []string{"ls"},
      //   Usage:     "list running containers",
      //   Action: func(c *cli.Context) {
      //     docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
      //     name := c.Args().First()

  }
  app.Run(os.Args)
}
