package main

import (
	"log"
	"os"

	"github.com/samalba/dockerclient"
)

// NewDockerClient creates a docker client
func NewDockerClient() dockerclient.Client {
	docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	return docker
}

// List the containers
func List(docker dockerclient.Client) []dockerclient.Container {
	// Get only running containers
	containers, err := docker.ListContainers(true, false, "")

	if err != nil {
		log.Fatal(err)
	}

	return containers
}

// Inspect a container
func Inspect(docker dockerclient.Client, id string) *dockerclient.ContainerInfo {
	info, _ := docker.InspectContainer(id)
	log.Println(info)
	return info
}

// Create a container
func Create(docker dockerclient.Client, name string, image string) string {
	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image: image,
		// Cmd:   []string{"bash"},
		AttachStdin: true,
		Tty:         true}
	id, err := docker.CreateContainer(containerConfig, name)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

// Start a container
func Start(docker dockerclient.Client, id string) string {
	// Start the container
	hostConfig := &dockerclient.HostConfig{}
	err := docker.StartContainer(id, hostConfig)
	if err != nil {
		log.Fatal(err)
	}
	return "OK"
}

// Stop a container
func Stop(docker dockerclient.Client, id string) string {
	// Stop the container (with 5 seconds timeout)
	err := docker.StopContainer(id, 5)
	if err != nil {
		log.Fatal(err)
	}
	return "OK"
}

// Build a container
func Build(docker dockerclient.Client, repoName string, context string) {
	// Build a docker image
	// some.tar contains the build context (Dockerfile any any files it needs to add/copy)
	dockerBuildContext, err := os.Open(context)
	defer dockerBuildContext.Close()

	buildImageConfig := &dockerclient.BuildImage{
		Context:        dockerBuildContext,
		RepoName:       repoName,
		SuppressOutput: false,
	}
	reader, err := docker.BuildImage(buildImageConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(reader)
}
