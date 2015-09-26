package containers_test

import (
	"testing"
	"plugins/containers"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestListNoContainers(t *testing.T) {
	fakeDataSource := func() []byte {
		return []byte{}
	}

	containers := containers.List(fakeDataSource)

	assert.Empty(t, containers)
}

func TestListOneUpContainer(t *testing.T) {
	fakeDataSource := func() []byte {
		response := `[
			{
				"Command":	"da command",
				"Created":	1398153384,
				"Id":		"da id",
				"Image":	"da image",
				"Names": [ "first name", "second name" ],
				"Ports": [
					{ "IP": "0.0.0.0", "PrivatePort": 27017, "PublicPort": 27017, "Type": "tcp" },
					{ "PublicPort": 28017, "Type": "tcp" }
				],
				"Status":	"Up 1 seconds"
			}
			]`
		return []byte(response)
	}

	containers := containers.List(fakeDataSource)

	assert.NotEmpty(t, containers)
	assert.Equal(t, len(containers), 1)
	assert.Equal(t, containers["da id"], []string{"first name", "second name"})
}

func TestListTwoUpContainers(t *testing.T) {
	fakeDataSource := func() []byte {
		response := `[
			{
				"Command":	"da command",
				"Created":	1398153384,
				"Id":		"da id",
				"Image":	"da image",
				"Names": [ "/first name", "/second name/first name" ],
				"Ports": [
					{ "IP": "0.0.0.0", "PrivatePort": 27017, "PublicPort": 27017, "Type": "tcp" },
					{ "PublicPort":	28017, "Type": "tcp" }
				],
				"Status":	"Up 1 seconds"
			},
			{
				"Command":	"da other command",
				"Created":	1398153484,
				"Id":		"da other id",
				"Image":	"da other image",
				"Names": [ "/another first name", "/another second name/another first name" ],
				"Ports": [ { "IP": "0.0.0.0", "PrivatePort": 123, "PublicPort": 234, "Type": "tcp" } ],
				"Status":	"Up 10 seconds"
			}
			]`
		return []byte(response)
	}

	containers := containers.List(fakeDataSource)

	assert.NotEmpty(t, containers)
	assert.Equal(t, len(containers), 2)
	assert.Equal(t, containers["da id"], []string{"/first name", "/second name/first name"})
	assert.Equal(t, containers["da other id"], []string{"/another first name", "/another second name/another first name"})

}

func TestGrabOneContainer(t *testing.T) {
	fakeDataSource := func(containerId string) []byte {
		response := `{
		   "ID" : "2a1f22565121fdaf1f43da2afc2e5495773953c57b01fd58f9be63caaa30f477",
		   "Created" : "2014-04-22T07:56:24.239219361Z",
		   "Image" : "724de0535b4f4b087dd78ad502744d5872629d78f7bea3a212bf702355bb7b3d",
		   "VolumesRW" : {},
		   "Volumes" : {},
		   "ExecDriver" : "native-0.1",
		   "Driver" : "aufs",
		   "Args" : [],
		   "HostnamePath" : "/home/docker/containers/2a1f22565121fdaf1f43da2afc2e5495773953c57b01fd58f9be63caaa30f477/hostname",
		   "NetworkSettings" : {
			  "PortMapping" : null,
			  "IPPrefixLen" : 16,
			  "Ports" : {
				 "28017/tcp" : null,
				 "27017/tcp" : [ { "HostIp" : "0.0.0.0", "HostPort" : "27017" } ]
			  },
			  "Bridge" : "docker0",
			  "Gateway" : "172.17.42.1",
			  "IPAddress" : "172.17.0.2"
		   },
		   "HostsPath" : "/home/docker/containers/2a1f22565121fdaf1f43da2afc2e5495773953c57b01fd58f9be63caaa30f477/hosts",
		   "HostConfig" : {
			  "PortBindings" : {
				 "28017/tcp" : null,
				 "27017/tcp" : [ { "HostPort" : "27017", "HostIp" : "0.0.0.0" } ]
			  },
			  "NetworkMode" : "",
			  "PublishAllPorts" : false,
			  "Privileged" : false,
			  "Links" : null,
			  "VolumesFrom" : null,
			  "Dns" : null,
			  "DnsSearch" : null,
			  "ContainerIDFile" : "",
			  "LxcConf" : [],
			  "Binds" : null
		   },
		   "Name" : "/mongodb",
		   "Path" : "mongod",
		   "ProcessLabel" : "",
		   "State" : {
			  "Pid" : 11680,
			  "StartedAt" : "2014-06-02T19:24:52.022935883Z",
			  "ExitCode" : 0,
			  "Running" : true,
			  "FinishedAt" : "2014-06-02T07:49:05.376035878Z"
		   },
		   "MountLabel" : "",
		   "Config" : {
			  "Domainname" : "",
			  "NetworkDisabled" : false,
			  "Hostname" : "2a1f22565121",
			  "OpenStdin" : false,
			  "Tty" : false,
			  "CpuShares" : 0,
			  "MemorySwap" : 0,
			  "StdinOnce" : false,
			  "Env" : null,
			  "Volumes" : null,
			  "Image" : "bacongobbler/mongodb",
			  "AttachStderr" : false,
			  "Memory" : 0,
			  "OnBuild" : null,
			  "PortSpecs" : null,
			  "User" : "",
			  "WorkingDir" : "",
			  "ExposedPorts" : {
				 "28017/tcp" : {},
				 "27017/tcp" : {}
			  },
			  "AttachStdout" : false,
			  "Entrypoint" : null,
			  "Cmd" : [ "mongod" ],
			  "AttachStdin" : false
		   },
		   "ResolvConfPath" : "/etc/resolv.conf"
		}`
		return []byte(response)
	}

	container := containers.Grab(fakeDataSource, "FOO")

	assert.True(t, container.Running)
	assert.Equal(t, container.Pid, 11680)
	refDate, _ := time.Parse(time.RFC3339, "2014-06-02T19:24:52.022935883Z")
	refUptime := int64(time.Since(refDate).Seconds())
	assert.Equal(t, container.Uptime, refUptime)
	assert.Equal(t, container.Image, "724de0535b4f4b087dd78ad502744d5872629d78f7bea3a212bf702355bb7b3d")
	assert.Equal(t, container.Name, "mongodb")
	assert.Equal(t, container.IPAddress, "172.17.0.2")
	assert.NotEmpty(t, container.Ports)
	assert.Equal(t, len(container.Ports), 2)
	firstPort := container.Ports[0]
	assert.Equal(t, firstPort.ContainerPort, "28017/tcp")
	assert.Equal(t, firstPort.HostIp, "")
	assert.Equal(t, firstPort.HostPort, "")
	secondPort := container.Ports[1]
	assert.Equal(t, secondPort.ContainerPort, "27017/tcp")
	assert.Equal(t, secondPort.HostIp, "0.0.0.0")
	assert.Equal(t, secondPort.HostPort, "27017")
}
