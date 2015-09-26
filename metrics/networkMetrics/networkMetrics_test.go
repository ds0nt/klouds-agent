package networkMetrics_test

import (
	"testing"
	"plugins/networkMetrics"
	"github.com/stretchr/testify/assert"
)

func TestNetworkDataParsing(t *testing.T) {
	fakeDataSource := func(containerId string) map[string]uint64 {
		fakeResult := make(map[string] uint64, 0)
		fakeResult["rx-ok"] = uint64(123)
		fakeResult["rx-err"] = uint64(34)
		fakeResult["tx-ok"] = uint64(4567)
		fakeResult["tx-err"] = uint64(78)
		return fakeResult
	}

	result := networkMetrics.Grab(fakeDataSource, "FOO")

	assert.Equal(t, result.RxOk, uint64(123))
	assert.Equal(t, result.RxErr, uint64(34))
	assert.Equal(t, result.TxOk, uint64(4567))
	assert.Equal(t, result.TxErr, uint64(78))
}

func TestNetwordNoData(t *testing.T) {
	fakeDataSource := func(containerId string) map[string]uint64 {
		fakeResult := make(map[string] uint64, 0)
		return fakeResult
	}

	result := networkMetrics.Grab(fakeDataSource, "FOO")

	assert.Equal(t, result.RxOk, uint64(0))
	assert.Equal(t, result.RxErr, uint64(0))
	assert.Equal(t, result.TxOk, uint64(0))
	assert.Equal(t, result.TxErr, uint64(0))

}
