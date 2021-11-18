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

	variant := "RR"
	var expectedE []string
	// Expected events
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Portugal has the world's best waves!")
	case "sub-group-2":
		expectedE = append(expectedE, "Publishing via ipfs is lit!")
	case "sub-group-3":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Surf trip to bali for 1050, just today!")
	case "sub-group-4":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Surf trip to bali for 1050, just today!")
	case "sub-group-5":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Surf trip to bali for 1050, just today!")
	case "sub-group-6":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Surf trip to bali for 1050, just today!")
	}

	Sync(ctx, ri.RunInfo, readyState)

	initMem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	initCpu, err := cpu.Times(false)
	if err != nil {
		return err
	}

	cfg := pubsub.DefaultConfig("PT", 10)
	cfg.TestgroundReady = true
	ps := pubsub.NewPubSub(ri.Node.dht, cfg)
	ps.SetHasOldPeer()

	Sync(ctx, ri.RunInfo, createdState)

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
	Sync(ctx, ri.RunInfo, subbedState)

	// Crash Routine
	if ri.RunInfo.RunEnv.RunParams.TestGroupID == "sub-group-6" && (ri.Node.info.GroupSeq == 0 || ri.Node.info.GroupSeq == 1) {
		ps.TerminateService()
	}

	time.Sleep(100 * time.Millisecond)
	Sync(ctx, ri.RunInfo, crashedState)

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("Publishing via ipfs is lit!", "ipfs T")
	case "pub-2":
		ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
	case "pub-3":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	case "pub-4":
		ps.MyPublish("Publishing via ipfs is sublime!", "ipfs T")
	}

	time.Sleep(time.Second)
	Sync(ctx, ri.RunInfo, finishedState)

	finalMem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	finalCpu, err := cpu.Times(false)
	if err != nil {
		return err
	}

	if !(ri.RunInfo.RunEnv.RunParams.TestGroupID == "sub-group-6" && (ri.Node.info.GroupSeq == 0 || ri.Node.info.GroupSeq == 1)) {

		events := ps.ReturnEventStats()
		subs := ps.ReturnSubStats()
		missed, duplicated := ps.ReturnCorrectnessStats(expectedE)
		runenv.R().RecordPoint("# Peers - ScoutSubs fault"+variant, float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
		runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
		runenv.R().RecordPoint("CPU used - ScoutSubs fault"+variant, finalCpu[0].User-initCpu[0].User)
		runenv.R().RecordPoint("Memory used - ScoutSubs fault"+variant, float64(finalMem.Used-initMem.Used)/(1024*1024))
		runenv.R().RecordPoint("# Events Missing - ScoutSubs fault"+variant, float64(missed))
		runenv.R().RecordPoint("# Events Duplicated - ScoutSubs fault"+variant, float64(duplicated))

		for _, ev := range events {
			runenv.R().RecordPoint("Event Latency - ScoutSubs fault"+variant, float64(ev))
		}
		for _, sb := range subs {
			runenv.R().RecordPoint("Sub Latency - ScoutSubs fault"+variant, float64(sb))
		}
	}

	Sync(ctx, ri.RunInfo, recordedState)

	if err := stager.End(); err != nil {
		return err
	}

	return nil
}
