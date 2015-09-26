package services

import (
	"strings"
	"net/http"
	"bytes"
	"net"
	"io/ioutil"
)

type DockerService struct {
	Url string
}

func (s *DockerService) getDockerJson(path string) []byte {
	var client *http.Client
	if strings.HasPrefix(s.Url, "http") {
		path = s.Url + path
		client = &http.Client{}
	} else {
		path = "http://docker" + path
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(network string, addr string) (net.Conn, error) {
					return net.Dial("unix", s.Url)
				},
			},
		}
	}
	req, _ := http.NewRequest("GET", path, bytes.NewBuffer(nil))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	dockerJson, _ := ioutil.ReadAll(resp.Body)
	return dockerJson
}

func (s *DockerService) List() []byte {
	return s.getDockerJson("/containers/json?all=1")
}

func (s *DockerService) Inspect(containerId string) []byte {
	return s.getDockerJson("/containers/" + containerId + "/json")
}
