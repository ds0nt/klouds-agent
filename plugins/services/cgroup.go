package services

import (
	"io/ioutil"
	"strings"
	"os"
	"os/exec"
	"strconv"
)

type CgroupService struct {
	DockerPrefix	string
}

func getLines(filePath string) []string {
	data, _ := ioutil.ReadFile(filePath)
	return strings.Split(string(data), "\n")
}

func (s *CgroupService) GetMemLines(containerId string) []string {
	return getLines("/sys/fs/cgroup/memory/" + s.DockerPrefix + containerId + "/memory.stat")
}

func (s *CgroupService) GetCpuLines(containerId string) []string {
	return getLines("/sys/fs/cgroup/cpu,cpuacct/" + s.DockerPrefix + containerId + "/cpuacct.stat")
}

func (s *CgroupService) GetNetworkLines(containerId string) map[string]uint64 {
	if _, err := os.Stat("/var/run/netns"); os.IsNotExist(err) {
		os.Mkdir("/var/run/netns", os.ModeDir)
	}
	data, err := ioutil.ReadFile("/sys/fs/cgroup/devices/" + s.DockerPrefix + containerId + "/tasks")
	if err != nil {
		panic(err)
	}
	pid := strings.Split(string(data), "\n")[0]
	os.Symlink("/proc/" + pid + "/ns/net", "/var/run/netns/" + containerId)
	cmd := exec.Command("/bin/ip", "netns", "exec", containerId, "netstat", "-i")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	os.Remove("/var/run/netns/" + containerId)
	lines := strings.Split(string(out), "\n")
	result := make(map[string]uint64, 0)
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "eth0") {
			elements := make([]string, 0)
			parts := strings.Split(lines[i], " ")
			for j := range parts {
				if strings.Trim(parts[j], " ") != "" {
					elements = append(elements, parts[j])
				}
			}
			result["rx-ok"], err = strconv.ParseUint(elements[3], 10, 64)
			if err != nil {
				panic(err)
			}
			result["rx-err"], err = strconv.ParseUint(elements[4], 10, 64)
			if err != nil {
				panic(err)
			}
			result["tx-ok"], err = strconv.ParseUint(elements[7], 10, 64)
			if err != nil {
				panic(err)
			}
			result["tx-err"], err = strconv.ParseUint(elements[8], 10, 64)
			if err != nil {
				panic(err)
			}
		}
	}
	return result
}

