package main

import (
	"context"

	"github.com/testground/sdk-go/network"
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
	"github.com/testground/sdk-go/sync"
)

func firstTry(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	var (
		enrolledState = sync.State("enrolled")
		releasedState = sync.State("released")

		ctx = context.Background()
	)

	// instantiate a sync service client, binding it to the RunEnv.
	client := sync.MustBoundClient(ctx, runenv)
	defer client.Close()

	// instantiate a network client; see 'Traffic shaping' in the docs.
	netclient := network.NewClient(client, runenv)
	runenv.RecordMessage("waiting for network initialization")

	// wait for the network to initialize; this should be pretty fast.
	netclient.MustWaitNetworkInitialized(ctx)
	runenv.RecordMessage("network initilization complete")

	// signal entry in the 'enrolled' state, and obtain a sequence number.
	seq := client.MustSignalEntry(ctx, enrolledState)
	runenv.RecordMessage("my sequence ID: %d", seq)

	// wait for all to signal to be at enroledState
	err := <-client.MustBarrier(ctx, enrolledState, runenv.TestInstanceCount).C
	if err != nil {
		return err
	}

	client.MustSignalEntry(ctx, releasedState)
	runenv.RecordMessage("I'm %d, and I'm free!", seq)

	return nil
}
