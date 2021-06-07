package test

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/test-plans/dht/utils"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/testground/sdk-go/runtime"
	"github.com/testground/sdk-go/sync"

	pubsub "github.com/pedroaston/contentpubsub"
)

func EventBurstScoutTest(runenv *runtime.RunEnv) error {
	commonOpts := GetCommonOpts(runenv)

	ctx, cancel := context.WithTimeout(context.Background(), commonOpts.Timeout)
	defer cancel()

	ri, err := Base(ctx, runenv, commonOpts)
	if err != nil {
		return err
	}

	if err := TestEventBurstScout(ctx, ri); err != nil {
		return err
	}
	Teardown(ctx, ri.RunInfo)

	return nil
}

func TestEventBurstScout(ctx context.Context, ri *DHTRunInfo) error {

	runenv := ri.RunEnv
	readyState := sync.State("ready")
	createdState := sync.State("created")
	subbedState := sync.State("subbed")
	finishedState := sync.State("finished")
	recordedState := sync.State("recorded")

	stager := utils.NewBatchStager(ctx, ri.Node.info.Seq, runenv.TestInstanceCount, "peer-records", ri.RunInfo)

	if err := stager.Begin(); err != nil {
		return err
	}

	variant := "BU"
	var expectedE int
	// Expected events
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		expectedE = 9
	case "sub-group-2":
		expectedE = 8
	case "sub-group-3":
		expectedE = 8
	case "sub-group-4":
		expectedE = 8
	case "sub-group-5":
		expectedE = 10
	case "sub-group-6":
		expectedE = 8
	}

	ri.Client.MustSignalEntry(ctx, readyState)
	err := <-ri.Client.MustBarrier(ctx, readyState, runenv.TestInstanceCount).C
	if err != nil {
		return err
	}

	initMem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	initCpu, err := cpu.Times(false)
	if err != nil {
		return err
	}

	ps := pubsub.NewPubSub(ri.Node.dht, pubsub.DefaultConfig("PT"))

	ri.Client.MustSignalEntry(ctx, createdState)
	err1stStop := <-ri.Client.MustBarrier(ctx, createdState, runenv.TestInstanceCount).C
	if err1stStop != nil {
		return err1stStop
	}

	// Subscribing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		ps.MySubscribe("portugal T/surf T")
		ps.MySubscribe("ipfs T")
	case "sub-group-2":
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("portugal T/soccer T")
	case "sub-group-3":
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/bali T")
	case "sub-group-4":
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/bali T/trip T/price R 1000 1500")
	case "sub-group-5":
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/trip T/price R 1000 2000")
	case "sub-group-6":
		ps.MySubscribe("ipfs T")
		ps.MySubscribe("surf T/trip T/price R 1000 1400")
	}

	time.Sleep(time.Second)
	ri.Client.MustSignalEntry(ctx, subbedState)
	err2ndStop := <-ri.Client.MustBarrier(ctx, subbedState, runenv.TestInstanceCount).C
	if err2ndStop != nil {
		return err2ndStop
	}

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		event1 := fmt.Sprintf("I already surfed %d portuguese beaches!", ri.Node.info.Seq%100)
		ps.MyPublish(event1, "portugal T/surf T")
	case "sub-group-2":
		event2 := fmt.Sprintf("Using IPFS, is %d times cooler than flying-cars!", ri.Node.info.Seq%100)
		ps.MyPublish(event2, "ipfs T")
	case "sub-group-3":
		event3 := fmt.Sprintf("Using IPFS, is %d times cooler than flip-flops!", ri.Node.info.Seq%100)
		ps.MyPublish(event3, "ipfs T")
	case "sub-group-4":
		event4 := fmt.Sprintf("Portugal will score 1%d goals at the world cup!", ri.Node.info.Seq%100)
		ps.MyPublish(event4, "portugal T/soccer T")
	case "sub-group-5":
		event5 := fmt.Sprintf("Surf trip to hawai for 15%d, just today!", ri.Node.info.Seq%100)
		pred5 := fmt.Sprintf("surf T/hawai T/trip T/price R 15%d 15%d", ri.Node.info.Seq%100, ri.Node.info.Seq%100)
		ps.MyPublish(event5, pred5)
	case "sub-group-6":
		event6 := fmt.Sprintf("Surf trip to bali for 10%d, just today!", ri.Node.info.Seq%100)
		pred6 := fmt.Sprintf("surf T/bali T/trip T/price R 10%d 10%d", ri.Node.info.Seq%100, ri.Node.info.Seq%100)
		ps.MyPublish(event6, pred6)
	}

	time.Sleep(time.Second)
	ri.Client.MustSignalEntry(ctx, finishedState)
	err3rdStop := <-ri.Client.MustBarrier(ctx, finishedState, runenv.TestInstanceCount).C
	if err3rdStop != nil {
		return err3rdStop
	}

	finalMem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	finalCpu, err := cpu.Times(false)
	if err != nil {
		return err
	}

	events := ps.ReturnEventStats()
	subs := ps.ReturnSubStats()
	nEScout := len(events)
	runenv.R().RecordPoint("# Peers - ScoutSubs eventBurst"+variant, float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
	runenv.R().RecordPoint("CPU used - ScoutSubs eventBurst"+variant, finalCpu[0].User-initCpu[0].User)
	runenv.R().RecordPoint("Memory used - ScoutSubs eventBurst"+variant, float64(finalMem.Used-initMem.Used)/(1024*1024))

	for _, ev := range events {
		runenv.R().RecordPoint("Event Latency - ScoutSubs eventBurst"+variant, float64(ev))
	}
	for _, sb := range subs {
		runenv.R().RecordPoint("Sub Latency - ScoutSubs eventBurst"+variant, float64(sb))
	}

	if expectedE > nEScout {
		runenv.R().RecordPoint("# Events Missing - ScoutSubs eventBurst"+variant, float64(expectedE-nEScout))
		runenv.R().RecordPoint("# Events Duplicated - ScoutSubs eventBurst"+variant, float64(0))
	} else {
		runenv.R().RecordPoint("# Events Missing - ScoutSubs eventBurst"+variant, float64(0))
		runenv.R().RecordPoint("# Events Duplicated - ScoutSubs eventBurst"+variant, float64(nEScout-expectedE))
	}

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
