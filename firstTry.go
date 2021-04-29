package main

import (
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

func firstTry(runenv *runtime.RunEnv, initCtx *run.InitContext) error {

	runenv.RecordMessage("Hello, Testground!")

	return nil
}
