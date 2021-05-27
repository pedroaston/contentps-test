package test

import (
	"context"
	"time"

	"github.com/libp2p/test-plans/dht/utils"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/testground/sdk-go/runtime"
	"github.com/testground/sdk-go/sync"

	pubsub "github.com/pedroaston/contentpubsub"
)

func FaultScoutTest(runenv *runtime.RunEnv) error {
	commonOpts := GetCommonOpts(runenv)

	ctx, cancel := context.WithTimeout(context.Background(), commonOpts.Timeout)
	defer cancel()

	ri, err := Base(ctx, runenv, commonOpts)
	if err != nil {
		return err
	}

	if err := TestFaultScout(ctx, ri); err != nil {
		return err
	}
	Teardown(ctx, ri.RunInfo)

	return nil
}

func TestFaultScout(ctx context.Context, ri *DHTRunInfo) error {

	runenv := ri.RunEnv
	readyState := sync.State("ready")
	createdState := sync.State("created")
	subbedState := sync.State("subbed")
	crashedState := sync.State("crashed")
	finishedState := sync.State("finished")
	recordedState := sync.State("recorded")

	stager := utils.NewBatchStager(ctx, ri.Node.info.Seq, runenv.TestInstanceCount, "peer-records", ri.RunInfo)

	if err := stager.Begin(); err != nil {
		return err
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

	variant := "RR"
	ps := pubsub.NewPubSub(ri.Node.dht, Region(ri.Node.info.Seq%3).String())

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

	// Crash Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-6":
		ps.TerminateService()
	}

	time.Sleep(100 * time.Millisecond)
	ri.Client.MustSignalEntry(ctx, crashedState)
	err3rdStop := <-ri.Client.MustBarrier(ctx, crashedState, runenv.TestInstanceCount).C
	if err3rdStop != nil {
		return err3rdStop
	}

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("Publishing via ipfs is lit!", "ipfs T")
	case "pub-2":
		ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
	case "pub-3":
		ps.MyPublish("Publishing via ipfs is sublime!", "ipfs T")
	case "pub-4":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	}

	time.Sleep(time.Second)
	ri.Client.MustSignalEntry(ctx, finishedState)
	err4thStop := <-ri.Client.MustBarrier(ctx, finishedState, runenv.TestInstanceCount).C
	if err4thStop != nil {
		return err4thStop
	}

	finalMem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	finalCpu, err := cpu.Times(false)
	if err != nil {
		return err
	}

	nEScout, _, latScout, _ := ps.ReturnReceivedEventsStats()
	runenv.R().RecordPoint("Number of peers - ScoutSubs fault"+variant, float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
	runenv.R().RecordPoint("Events received - ScoutSubs fault"+variant, float64(nEScout))
	runenv.R().RecordPoint("Avg event latency - ScoutSubs fault"+variant, float64(latScout))
	runenv.R().RecordPoint("Avg time to sub - ScoutSubs fault"+variant, float64(ps.ReturnSubStats()))
	runenv.R().RecordPoint("CPU used - ScoutSubs fault"+variant, finalCpu[0].User-initCpu[0].User)
	runenv.R().RecordPoint("Memory used - ScoutSubs fault"+variant, float64(finalMem.Used)-float64(initMem.Used))

	ri.Client.MustSignalEntry(ctx, recordedState)
	err5thStop := <-ri.Client.MustBarrier(ctx, recordedState, runenv.TestInstanceCount).C
	if err5thStop != nil {
		return err5thStop
	}

	if err := stager.End(); err != nil {
		return err
	}

	return nil
}