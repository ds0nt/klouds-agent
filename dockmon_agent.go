package main

import (
	"plugins/containers"
	"plugins/memoryMetrics"
	"encoding/json"
	"plugins/cpuMetrics"
	"plugins/networkMetrics"
	"flag"
	"io/ioutil"
	"services"
	"strings"
	"net/http"
	"os"
	"fmt"
)


type dataStruct struct {
	Memory		memoryMetrics.MemoryStruct
	Cpu			cpuMetrics.CpuStruct
	Info		containers.ContainerStruct
	Network		networkMetrics.NetworkStruct
}

type agentConf struct {
	DockerUrl			string
	CgroupDockerPrefix	string
	CollectorUrl		string
	AccountId			string
}

func readConf() agentConf {
	var confFile string
	flag.StringVar(&confFile, "c", "/etc/dockmon_agent.conf", "Configuration file location (absolute path)")
	flag.Parse()
	var conf agentConf
	confJson, err := ioutil.ReadFile(confFile)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(confJson, &conf)
	return conf
}

func getJsonData(conf agentConf) map[string] *dataStruct {
	dockerService := &services.DockerService{Url: conf.DockerUrl}
	cgroupService := &services.CgroupService{DockerPrefix: conf.CgroupDockerPrefix}
	containersMap := containers.List(dockerService.List)
	name2Id := make(map[string]string, 0)
	data := make(map[string]*dataStruct)
	for id, names := range containersMap {
		data[id] = &dataStruct{}
		data[id].Info = containers.Grab(dockerService.Inspect, id)
		name2Id[data[id].Info.Name] = id
		for i := 0; i < len(names); i++ {
			name := names[i]
			if strings.LastIndex(name, "/") > 0 {
				parent := strings.Split(name, "/")
				data[id].Info.Parents = append(data[id].Info.Parents, parent[1])
			}
		}
		if data[id].Info.Running {
			data[id].Memory = memoryMetrics.Grab(cgroupService.GetMemLines, id)
			data[id].Cpu = cpuMetrics.Grab(cgroupService.GetCpuLines, id)
			data[id].Network = networkMetrics.Grab(cgroupService.GetNetworkLines, id)
		}
	}
	for _, value := range data {
		if value.Info.Parents != nil {
			for i := 0; i < len(value.Info.Parents); i++ {
				parentId := name2Id[value.Info.Parents[i]]
				data[parentId].Info.Children = append(data[parentId].Info.Children, value.Info.Name)
			}
		}
	}
	return data
}


func main() {
	conf := readConf()
	jsonData, err := json.Marshal(getJsonData(conf))
	if err != nil {
		panic(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", conf.CollectorUrl + "/server/" + hostname, strings.NewReader(string(jsonData)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-user", conf.AccountId)
	response, err := client.Do(req)
//	response, err := http.Post(conf.CollectorUrl + "/server/" + hostname, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		panic(err)
	}

	fmt.Println(conf.CollectorUrl + "/server/" + hostname)
	fmt.Println(string(jsonData))

	if response.StatusCode > 399 {
		fmt.Errorf("status code received : " + string(response.StatusCode))
	}
}
