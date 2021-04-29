package main

import (
	"github.com/testground/sdk-go/run"
)

var testcases = map[string]interface{}{
	"firstTry": run.InitializedTestCaseFn(firstTry),
}

func main() {
	run.InvokeMap(testcases)
}
