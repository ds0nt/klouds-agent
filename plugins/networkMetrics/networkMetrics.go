package networkMetrics

type NetworkStruct struct {
	RxOk	uint64
	RxErr	uint64
	TxOk	uint64
	TxErr	uint64
}

func Grab(dataSource func(string) map[string]uint64, containerId string) NetworkStruct {
	lines := dataSource(containerId)

	return NetworkStruct{RxOk: lines["rx-ok"], RxErr: lines["rx-err"], TxOk: lines["tx-ok"], TxErr: lines["tx-err"]}
}
