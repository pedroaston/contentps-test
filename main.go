package main

import (
	test "github.com/libp2p/test-plans/dht/test"
	"github.com/testground/sdk-go/run"
)

var testCases = map[string]interface{}{
	"normal-scout-test":     test.NormalScoutTest,
	"subburst-scout-test":   test.SubBurstScoutTest,
	"eventburst-scout-test": test.EventBurstScoutTest,
	"fault-scout-test":      test.FaultScoutTest,
	"longrun-scout-test":    test.LongRunScoutTest,
	"fast-test":             test.FastTest,
}

func main() {
	run.InvokeMap(testCases)
}
