package memoryMetrics

import (
	"strings"
	"math/big"
)

type MemoryStruct struct {
	Total	*big.Int
	Rss		*big.Int
	Cache	*big.Int
	Swap	*big.Int
}

func Grab(dataSource func(string) []string, containerId string) MemoryStruct {
	lines := dataSource(containerId)
	total := big.NewInt(-1)
	rss := big.NewInt(-1)
	cache := big.NewInt(-1)
	swap := big.NewInt(-1)
	for i := 0; i < len(lines); i++ {
		elements := strings.Split(lines[i], " ")
		switch elements[0] {
		case "cache" :
			cache.SetString(elements[1], 10)
		case "rss" :
			rss.SetString(elements[1], 10)
		case "swap" :
			swap.SetString(elements[1], 10)
		case "hierarchical_memory_limit" :
			total.SetString(elements[1], 10)
		}
	}
	return MemoryStruct{total, rss, cache, swap}
}
