package main

import (
	test "github.com/libp2p/test-plans/dht/test"
	"github.com/testground/sdk-go/run"
)

var testCases = map[string]interface{}{
	"1st-test":  test.FirstTest,
	"fast-test": test.FastTest,
}

func main() {
	run.InvokeMap(testCases)
}
