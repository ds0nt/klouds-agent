package cpuMetrics

import (
	"strings"
	"strconv"
)

type CpuStruct struct {
	User	int64
	System	int64
}

func Grab(dataSource func(string) []string, containerId string) CpuStruct {
	lines := dataSource(containerId)
	user := int64(-1)
	system := int64(-1)
	for i := 0; i < len(lines); i++ {
		elements := strings.Split(lines[i], " ")
		switch elements[0] {
		case "user" :
			intVal, err := strconv.ParseInt(elements[1], 10, 64)
			if err != nil {
				panic(err)
			}
			user = intVal
		case "system" :
			intVal, err := strconv.ParseInt(elements[1], 10, 64)
			if err != nil {
				panic(err)
			}
			system = intVal
		}
	}

	return CpuStruct{User: user, System: system}
}
