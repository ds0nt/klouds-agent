package memoryMetrics_test

import (
	"testing"
	"plugins/memoryMetrics"
	"github.com/stretchr/testify/assert"
)

func TestMemoryDataParsing(t *testing.T) {
	fakeDataSource := func(containerId string) []string {
		return []string{"cache 123", "rss 2345", "swap 45678", "hierarchical_memory_limit 9876"}
	}

	result := memoryMetrics.Grab(fakeDataSource, "FOO")

	assert.NotNil(t, result)
	assert.Equal(t, result.Cache, 123)
	assert.Equal(t, result.Rss, 2345)
	assert.Equal(t, result.Swap, 45678)
	assert.Equal(t, result.Total, 9876)
}

func TestMemoryNoData(t *testing.T) {
	fakeDataSource := func(containerId string) []string {
		return []string{}
	}

	result := memoryMetrics.Grab(fakeDataSource, "FOO")

	assert.NotNil(t, result)
	assert.Equal(t, result.Cache, -1)
	assert.Equal(t, result.Rss, -1)
	assert.Equal(t, result.Swap, -1)
	assert.Equal(t, result.Total, -1)
}
