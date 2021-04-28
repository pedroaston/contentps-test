// Welcome, testground plan writer!
// If you are seeing this for the first time, check out our documentation!
// https://app.gitbook.com/@protocol-labs/s/testground/

package main

import (
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

func main() {
	run.Invoke(runf)
}

func runf(runenv *runtime.RunEnv) error {
	runenv.RecordMessage("Hello, Testground!")
	return nil
}
