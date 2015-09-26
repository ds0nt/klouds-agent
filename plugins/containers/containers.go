package containers

import (
	"encoding/json"
	"time"
	"strings"
)

type DockerListStruct struct {
	Id		string
	Names	[]string
}

func List(dataSource func() []byte) map[string][]string {
	var dockerApiData []DockerListStruct
	dockerJson := dataSource()
	json.Unmarshal(dockerJson, &dockerApiData)
	var containers = make(map[string][]string, 0)
	for i := 0; i < len(dockerApiData); i++ {
		containers[dockerApiData[i].Id] = dockerApiData[i].Names
	}
	return containers
}

type DockerContainerStruct struct {
	Name			string
	Image			string
	NetworkSettings	NetworkStruct
	State			StateStruct
}

type StateStruct struct {
	Running		bool
	Pid			int64
	StartedAt	time.Time
	Uptime		int64
}

type NetworkStruct struct {
	IPAddress	string
	Ports		map[string][]PortStruct
}

type PortStruct struct {
	ContainerPort	string
	HostPort		string
	HostIp			string
}

type ContainerStruct struct {
	Name		string
	Image		string
	Running		bool
	Pid			int64
	Uptime		int64
	IPAddress	string
	Ports		[]PortStruct
	Parents		[]string
	Children	[]string
}

func Grab(dataSource func(containerId string) []byte, containerId string) ContainerStruct {
	var dockerApiData DockerContainerStruct
	dockerJson := dataSource(containerId)
	json.Unmarshal(dockerJson, &dockerApiData)
	result := ContainerStruct{}
	if dockerApiData.Image != "" {
		result.Name = strings.Split(dockerApiData.Name, "/")[1]
		result.Image = dockerApiData.Image
		result.Running = dockerApiData.State.Running
		result.Pid = dockerApiData.State.Pid
		result.IPAddress = dockerApiData.NetworkSettings.IPAddress
		result.Uptime = int64(time.Since(dockerApiData.State.StartedAt).Seconds())
		for key, value := range dockerApiData.NetworkSettings.Ports {
			portStruct := PortStruct{ContainerPort: key}
			if len(value) > 0 {
				portStruct.HostPort = value[0].HostPort
				portStruct.HostIp = value[0].HostIp
			}
			result.Ports = append(result.Ports, portStruct)
		}
	}
	return result
}
