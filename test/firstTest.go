package test

import (
	"context"
	"time"

	"github.com/libp2p/test-plans/dht/utils"

	"github.com/testground/sdk-go/runtime"
	"github.com/testground/sdk-go/sync"

	pubsub "github.com/pedroaston/contentpubsub"
)

func FirstTest(runenv *runtime.RunEnv) error {
	commonOpts := GetCommonOpts(runenv)

	ctx, cancel := context.WithTimeout(context.Background(), commonOpts.Timeout)
	defer cancel()

	ri, err := Base(ctx, runenv, commonOpts)
	if err != nil {
		return err
	}

	if err := TestSomething(ctx, ri); err != nil {
		return err
	}
	Teardown(ctx, ri.RunInfo)

	return nil
}

func TestSomething(ctx context.Context, ri *DHTRunInfo) error {

	runenv := ri.RunEnv
	readyState := sync.State("ready")
	createdState := sync.State("created")
	subbedState := sync.State("subbed")
	finishedState := sync.State("finished")

	ri.Client.MustSignalEntry(ctx, readyState)
	err := <-ri.Client.MustBarrier(ctx, readyState, runenv.TestInstanceCount).C
	if err != nil {
		return err
	}

	ps := pubsub.NewPubSub(ri.Node.dht, "PT", "EU")

	ri.Client.MustSignalEntry(ctx, createdState)
	err1stStop := <-ri.Client.MustBarrier(ctx, createdState, runenv.TestInstanceCount).C
	if err1stStop != nil {
		return err1stStop
	}

	// Subscribing Routine
	switch ri.Node.info.Properties.BootstrapStrategy {
	/*
		case 1:
			ps.MySubscribe("portugal T/surf T")
			ps.MySubscribe("ipfs T")
		case 2:
			ps.MySubscribe("ipfs T")
			ps.MySubscribe("portugal T/soccer T")
		case 3:
			ps.MySubscribe("ipfs T")
			ps.MySubscribe("surf T/bali T")
	*/
	case 4:
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/bali T/trip T/price R 1000 1500")
	case 5:
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/trip T/price R 1000 2000")
	case 6:
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/trip T/price R 1000 1400")
	case 7:
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("soccer T/slb T")

	}

	stager := utils.NewBatchStager(ctx, ri.Node.info.Seq, runenv.TestInstanceCount, "peer-records", ri.RunInfo)

	if err := stager.Begin(); err != nil {
		return err
	}

	time.Sleep(time.Second)
	ri.Client.MustSignalEntry(ctx, subbedState)

	err2ndStop := <-ri.Client.MustBarrier(ctx, subbedState, runenv.TestInstanceCount).C
	if err2ndStop != nil {
		return err2ndStop
	}

	/*
		panic: runtime error: invalid memory address or nil pointer dereference
		ERROR [signal SIGSEGV: segmentation violation code=0x1 addr=0x10 pc=0xd2dd32]
	*/

	// Publishing Routine
	switch ri.Node.info.Properties.BootstrapStrategy {
	/*
		case 1:
			ps.MyPublish("Publishing via ipfs is lit!", "ipfs T")
		case 2:
			ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
		case 3:
			ps.MyPublish("Publishing via ipfs is sublime!", "ipfs T")
		case 4:
			ps.MyPublish("Bali some good waves!", "surf T/bali T")
		case 5:
			ps.MyPublish("Benfica is the best football club of the world!", "soccer T/slb T")
	*/
	case 6:
		ps.MyPublish("Publishing via ipfs is exciting!", "ipfs T")
	case 7:
		ps.MyPublish("surf trip for 1200 euros! Promo valid for a week!", "surf T/trip T/bali T/price R 1200 1200")
	}

	time.Sleep(time.Second)
	ri.Client.MustSignalEntry(ctx, finishedState)
	err3rdStop := <-ri.Client.MustBarrier(ctx, finishedState, runenv.TestInstanceCount).C
	if err3rdStop != nil {
		return err3rdStop
	}

	nEScout, nEFast, latScout, latFast := ps.ReturnReceivedEventsStats()
	runenv.R().RecordPoint("Number of events received via ScoutSubs", float64(nEScout))
	runenv.R().RecordPoint("Number of events received via FastDelivery", float64(nEFast))
	runenv.R().RecordPoint("Avg latency of events received via ScoutSubs", float64(latScout))
	runenv.R().RecordPoint("Avg latency of events received via ScoutSubs", float64(latFast))

	if err := stager.End(); err != nil {
		return err
	}

	return nil
}
