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
	recordedState := sync.State("recorded")

	seq := ri.Client.MustSignalEntry(ctx, readyState)
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
	switch seq {
	case 1:
		ps.MySubscribe("portugal T/surf T")
		ps.MySubscribe("ipfs T")
	case 3:
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("portugal T/soccer T")
	case 6:
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/bali T")
	case 10:
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/bali T/trip T/price R 1000 1500")
	case 12:
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/trip T/price R 1000 2000")
	case 14:
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/trip T/price R 1000 1400")
	default:
		ps.MySubscribe("ipfs T")
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

	// Publishing Routine
	switch seq {
	case 2:
		ps.MyPublish("Publishing via ipfs is lit!", "ipfs T")
		//case 3:
		//ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
		//case 3:
		//ps.MyPublish("Publishing via ipfs is sublime!", "ipfs T")
		//case 4:
		//ps.MyPublish("Bali some good waves!", "surf T/bali T")
		//case 5:
		//ps.MyPublish("Publishing via ipfs is exciting!", "ipfs T")
		//ps.MyPublish("Benfica is the best football club of the world!", "soccer T/slb T")
		//case 6:
		//ps.MyPublish("Publishing via ipfs is exciting!", "ipfs T")
	}

	time.Sleep(time.Second)
	ri.Client.MustSignalEntry(ctx, finishedState)
	err3rdStop := <-ri.Client.MustBarrier(ctx, finishedState, runenv.TestInstanceCount).C
	if err3rdStop != nil {
		return err3rdStop
	}

	nEScout, nEFast, latScout, latFast := ps.ReturnReceivedEventsStats()
	runenv.R().RecordPoint("Number of peers", float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.R().RecordPoint("Number of events received via ScoutSubs", float64(nEScout))
	runenv.R().RecordPoint("Number of events received via FastDelivery", float64(nEFast))
	runenv.R().RecordPoint("Avg latency of events received via ScoutSubs", float64(latScout))
	runenv.R().RecordPoint("Avg latency of events received via FastDelivery", float64(latFast))

	ri.Client.MustSignalEntry(ctx, recordedState)
	err4thStop := <-ri.Client.MustBarrier(ctx, recordedState, runenv.TestInstanceCount).C
	if err4thStop != nil {
		return err4thStop
	}

	if err := stager.End(); err != nil {
		return err
	}

	return nil
}
