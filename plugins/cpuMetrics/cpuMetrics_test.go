package cpuMetrics_test

import (
	"testing"
	"plugins/cpuMetrics"
	"github.com/stretchr/testify/assert"
)

func TestCpuDataParsing(t *testing.T) {
	fakeDataSource := func(containerId string) []string {
		return []string{"user 24", "system 35"}
	}

	result := cpuMetrics.Grab(fakeDataSource, "FOO")

	assert.NotNil(t, result)
	assert.Equal(t, result.User, 24)
	assert.Equal(t, result.System, 35)
}

func TestCpuNoData(t *testing.T) {
	fakeDataSource := func(containerId string) []string {
		return []string{}
	}

	result := cpuMetrics.Grab(fakeDataSource, "FOO")

	assert.NotNil(t, result)
	assert.Equal(t, result.User, -1)
	assert.Equal(t, result.System, -1)
}
