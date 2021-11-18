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

func CompleteScoutTest(runenv *runtime.RunEnv) error {
	commonOpts := GetCommonOpts(runenv)

	ctx, cancel := context.WithTimeout(context.Background(), commonOpts.Timeout)
	defer cancel()

	ri, err := Base(ctx, runenv, commonOpts)
	if err != nil {
		return err
	}

	if err := TestCompleteScout(ctx, ri); err != nil {
		return err
	}
	Teardown(ctx, ri.RunInfo)

	return nil
}

func TestCompleteScout(ctx context.Context, ri *DHTRunInfo) error {

	variant := "BU"
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

	var expectedE []string
	// Expected events
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Portugal has the world's best waves!")
	case "sub-group-2":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Portugal won the world cup!")
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
	// Begining normal

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

	time.Sleep(2 * time.Second)
	Sync(ctx, ri.RunInfo, subbedState)

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("Publishing via ipfs is lit!", "ipfs T")
	case "pub-2":
		ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
	case "pub-3":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	case "pub-4":
		ps.MyPublish("Portugal won the world cup!", "portugal T/soccer T")
	}

	time.Sleep(2 * time.Second)
	Sync(ctx, ri.RunInfo, finishedState)

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
	missed, duplicated := ps.ReturnCorrectnessStats(expectedE)
	runenv.R().RecordPoint("# Peers - ScoutSubs normal"+variant, float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
	runenv.R().RecordPoint("CPU used - ScoutSubs normal"+variant, finalCpu[0].User-initCpu[0].User)
	runenv.R().RecordPoint("Memory used - ScoutSubs normal"+variant, float64(finalMem.Used-initMem.Used)/(1024*1024))
	runenv.R().RecordPoint("# Events Missing - ScoutSubs normal"+variant, float64(missed))
	runenv.R().RecordPoint("# Events Duplicated - ScoutSubs normal"+variant, float64(duplicated))

	for _, ev := range events {
		runenv.R().RecordPoint("Event Latency - ScoutSubs normal"+variant, float64(ev))
	}
	for _, sb := range subs {
		runenv.R().RecordPoint("Sub Latency - ScoutSubs normal"+variant, float64(sb))
	}

	Sync(ctx, ri.RunInfo, recordedState)
	// Begining subburst

	initMem, err = mem.VirtualMemory()
	if err != nil {
		return err
	}
	initCpu, err = cpu.Times(false)
	if err != nil {
		return err
	}

	Sync(ctx, ri.RunInfo, readyState)

	// Subscribing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		ps.MySubscribe("surf T/madeira T")
		ps.MySubscribe("dogecoin T")
	case "sub-group-2":
		ps.MySubscribe("dogecoin T")
		ps.MySubscribe("surf T/azores T")
	case "sub-group-3":
		ps.MySubscribe("dogecoin T")
		ps.MySubscribe("surf T/bali T")
	case "sub-group-4":
		ps.MySubscribe("dogecoin T")
		ps.MySubscribe("bitcoin T/price R 10000 15000")
	case "sub-group-5":
		ps.MySubscribe("dogecoin T")
		ps.MySubscribe("temperature R 30 40")
	case "sub-group-6":
		ps.MySubscribe("dogecoin T")
		ps.MySubscribe("ronaldo T/sporting T")
	}

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("Publishing via ipfs is lit!", "ipfs T")
	case "pub-2":
		ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
	case "pub-3":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	case "pub-4":
		ps.MyPublish("Portugal won the world cup!", "portugal T/soccer T")
	}

	time.Sleep(3 * time.Second)
	Sync(ctx, ri.RunInfo, finishedState)

	finalMem, err = mem.VirtualMemory()
	if err != nil {
		return err
	}
	finalCpu, err = cpu.Times(false)
	if err != nil {
		return err
	}

	events = ps.ReturnEventStats()
	subs = ps.ReturnSubStats()
	missed, duplicated = ps.ReturnCorrectnessStats(expectedE)
	runenv.R().RecordPoint("# Peers - ScoutSubs subBurst"+variant, float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
	runenv.R().RecordPoint("CPU used - ScoutSubs subBurst"+variant, finalCpu[0].User-initCpu[0].User)
	runenv.R().RecordPoint("Memory used - ScoutSubs subBurst"+variant, float64(finalMem.Used-initMem.Used)/(1024*1024))
	runenv.R().RecordPoint("# Events Missing - ScoutSubs subBurst"+variant, float64(missed))
	runenv.R().RecordPoint("# Events Duplicated - ScoutSubs subBurst"+variant, float64(duplicated))

	for _, ev := range events {
		runenv.R().RecordPoint("Event Latency - ScoutSubs subBurst"+variant, float64(ev))
	}
	for _, sb := range subs {
		runenv.R().RecordPoint("Sub Latency - ScoutSubs subBurst"+variant, float64(sb))
	}

	Sync(ctx, ri.RunInfo, recordedState)
	// Begining event burst

	expectedE = nil
	// Expected events
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		expectedE = append(expectedE, "Using IPFS, is 10 times cooler than flip-flops!", "Using IPFS, is 11 times cooler than flip-flops!",
			"Using IPFS, is 12 times cooler than flip-flops!", "Using IPFS, is 13 times cooler than flip-flops!",
			"Using IPFS, is 14 times cooler than flip-flops!", "Using IPFS, is 15 times cooler than flip-flops!",
			"Using IPFS, is 16 times cooler than flip-flops!", "Using IPFS, is 10 times cooler than flying-cars!",
			"Using IPFS, is 11 times cooler than flying-cars!", "Using IPFS, is 12 times cooler than flying-cars!",
			"Using IPFS, is 13 times cooler than flying-cars!", "Using IPFS, is 14 times cooler than flying-cars!",
			"Using IPFS, is 15 times cooler than flying-cars!", "Using IPFS, is 16 times cooler than flying-cars!",
			"I already surfed 10 portuguese beaches!", "I already surfed 11 portuguese beaches!",
			"I already surfed 12 portuguese beaches!", "I already surfed 13 portuguese beaches!",
			"I already surfed 14 portuguese beaches!", "I already surfed 15 portuguese beaches!")
	case "sub-group-2":
		expectedE = append(expectedE, "Using IPFS, is 10 times cooler than flip-flops!", "Using IPFS, is 11 times cooler than flip-flops!",
			"Using IPFS, is 12 times cooler than flip-flops!", "Using IPFS, is 13 times cooler than flip-flops!",
			"Using IPFS, is 14 times cooler than flip-flops!", "Using IPFS, is 15 times cooler than flip-flops!",
			"Using IPFS, is 16 times cooler than flip-flops!", "Using IPFS, is 10 times cooler than flying-cars!",
			"Using IPFS, is 11 times cooler than flying-cars!", "Using IPFS, is 12 times cooler than flying-cars!",
			"Using IPFS, is 13 times cooler than flying-cars!", "Using IPFS, is 14 times cooler than flying-cars!",
			"Using IPFS, is 15 times cooler than flying-cars!", "Using IPFS, is 16 times cooler than flying-cars!",
			"Portugal will score 10 goals at the world cup!",
			"Portugal will score 11 goals at the world cup!", "Portugal will score 12 goals at the world cup!",
			"Portugal will score 13 goals at the world cup!", "Portugal will score 14 goals at the world cup!",
			"Portugal will score 15 goals at the world cup!", "Portugal will score 16 goals at the world cup!")
	case "sub-group-3":
		expectedE = append(expectedE, "Surf trip to bali for 1100, just today!", "Surf trip to bali for 1101, just today!",
			"Surf trip to bali for 1102, just today!", "Surf trip to bali for 1103, just today!", "Surf trip to bali for 1104, just today!",
			"Surf trip to bali for 1105, just today!", "Surf trip to bali for 1106, just today!")
	case "sub-group-4":
		expectedE = append(expectedE, "Using IPFS, is 10 times cooler than flip-flops!", "Using IPFS, is 11 times cooler than flip-flops!",
			"Using IPFS, is 12 times cooler than flip-flops!", "Using IPFS, is 13 times cooler than flip-flops!",
			"Using IPFS, is 14 times cooler than flip-flops!", "Using IPFS, is 15 times cooler than flip-flops!",
			"Using IPFS, is 16 times cooler than flip-flops!", "Using IPFS, is 10 times cooler than flying-cars!",
			"Using IPFS, is 11 times cooler than flying-cars!", "Using IPFS, is 12 times cooler than flying-cars!",
			"Using IPFS, is 13 times cooler than flying-cars!", "Using IPFS, is 14 times cooler than flying-cars!",
			"Using IPFS, is 15 times cooler than flying-cars!", "Using IPFS, is 16 times cooler than flying-cars!",
			"Surf trip to bali for 1100, just today!", "Surf trip to bali for 1101, just today!",
			"Surf trip to bali for 1102, just today!")
	case "sub-group-5":
		expectedE = append(expectedE, "Surf trip to bali for 1100, just today!", "Surf trip to bali for 1101, just today!",
			"Surf trip to bali for 1102, just today!", "Surf trip to bali for 1103, just today!", "Surf trip to bali for 1104, just today!",
			"Surf trip to bali for 1105, just today!", "Surf trip to bali for 1106, just today!",
			"Surf trip to hawai for 1600, just today!", "Surf trip to hawai for 1601, just today!", "Surf trip to hawai for 1602, just today!",
			"Surf trip to hawai for 1603, just today!", "Surf trip to hawai for 1604, just today!", "Surf trip to hawai for 1605, just today!",
			"Surf trip to hawai for 1606, just today!")
	case "sub-group-6":
		expectedE = append(expectedE, "Using IPFS, is 10 times cooler than flip-flops!", "Using IPFS, is 11 times cooler than flip-flops!",
			"Using IPFS, is 12 times cooler than flip-flops!", "Using IPFS, is 13 times cooler than flip-flops!",
			"Using IPFS, is 14 times cooler than flip-flops!", "Using IPFS, is 15 times cooler than flip-flops!",
			"Using IPFS, is 16 times cooler than flip-flops!", "Using IPFS, is 10 times cooler than flying-cars!",
			"Using IPFS, is 11 times cooler than flying-cars!", "Using IPFS, is 12 times cooler than flying-cars!",
			"Using IPFS, is 13 times cooler than flying-cars!", "Using IPFS, is 14 times cooler than flying-cars!",
			"Using IPFS, is 15 times cooler than flying-cars!", "Using IPFS, is 16 times cooler than flying-cars!",
			"Surf trip to bali for 1100, just today!", "Surf trip to bali for 1101, just today!",
			"Surf trip to bali for 1102, just today!", "Surf trip to bali for 1103, just today!", "Surf trip to bali for 1104, just today!",
			"Surf trip to bali for 1105, just today!", "Surf trip to bali for 1106, just today!")
	}

	initMem, err = mem.VirtualMemory()
	if err != nil {
		return err
	}
	initCpu, err = cpu.Times(false)
	if err != nil {
		return err
	}

	Sync(ctx, ri.RunInfo, readyState)

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		event6 := fmt.Sprintf("Surf trip to bali for 110%d, just today!", ri.Node.info.GroupSeq)
		pred6 := fmt.Sprintf("surf T/bali T/trip T/price R 110%d 110%d", ri.Node.info.GroupSeq, ri.Node.info.GroupSeq)
		ps.MyPublish(event6, pred6)
	case "sub-group-2":
		event5 := fmt.Sprintf("Surf trip to hawai for 160%d, just today!", ri.Node.info.GroupSeq)
		pred5 := fmt.Sprintf("surf T/hawai T/trip T/price R 160%d 160%d", ri.Node.info.GroupSeq, ri.Node.info.GroupSeq)
		ps.MyPublish(event5, pred5)
	case "sub-group-3":
		event3 := fmt.Sprintf("Using IPFS, is 1%d times cooler than flip-flops!", ri.Node.info.GroupSeq)
		ps.MyPublish(event3, "ipfs T")
	case "sub-group-4":
		event4 := fmt.Sprintf("Portugal will score 1%d goals at the world cup!", ri.Node.info.GroupSeq)
		ps.MyPublish(event4, "portugal T/soccer T")
	case "sub-group-5":
		event2 := fmt.Sprintf("Using IPFS, is 1%d times cooler than flying-cars!", ri.Node.info.GroupSeq)
		ps.MyPublish(event2, "ipfs T")
	case "sub-group-6":
		event1 := fmt.Sprintf("I already surfed 1%d portuguese beaches!", ri.Node.info.GroupSeq)
		ps.MyPublish(event1, "portugal T/surf T")
	}

	time.Sleep(10 * time.Second)
	Sync(ctx, ri.RunInfo, finishedState)

	finalMem, err = mem.VirtualMemory()
	if err != nil {
		return err
	}
	finalCpu, err = cpu.Times(false)
	if err != nil {
		return err
	}

	events = ps.ReturnEventStats()
	subs = ps.ReturnSubStats()
	missed, duplicated = ps.ReturnCorrectnessStats(expectedE)
	runenv.R().RecordPoint("# Peers - ScoutSubs eventBurst"+variant, float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
	runenv.R().RecordPoint("CPU used - ScoutSubs eventBurst"+variant, finalCpu[0].User-initCpu[0].User)
	runenv.R().RecordPoint("Memory used - ScoutSubs eventBurst"+variant, float64(finalMem.Used-initMem.Used)/(1024*1024))
	runenv.R().RecordPoint("# Events Missing - ScoutSubs eventBurst"+variant, float64(missed))
	runenv.R().RecordPoint("# Events Duplicated - ScoutSubs eventBurst"+variant, float64(duplicated))

	for _, ev := range events {
		runenv.R().RecordPoint("Event Latency - ScoutSubs eventBurst"+variant, float64(ev))
	}
	for _, sb := range subs {
		runenv.R().RecordPoint("Sub Latency - ScoutSubs eventBurst"+variant, float64(sb))
	}

	Sync(ctx, ri.RunInfo, recordedState)
	// Begining Fault scenario

	expectedE = nil
	// Expected events
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Portugal has the world's best waves!")
	case "sub-group-2":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Portugal won the world cup!")
	case "sub-group-3":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Surf trip to bali for 1050, just today!")
	case "sub-group-4":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Surf trip to bali for 1050, just today!")
	case "sub-group-5":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Surf trip to bali for 1050, just today!")
	case "sub-group-6":
		expectedE = append(expectedE, "Publishing via ipfs is lit!", "Surf trip to bali for 1050, just today!")
	}

	initMem, err = mem.VirtualMemory()
	if err != nil {
		return err
	}
	initCpu, err = cpu.Times(false)
	if err != nil {
		return err
	}

	Sync(ctx, ri.RunInfo, readyState)

	// Crash Routine
	if ri.RunInfo.RunEnv.RunParams.TestGroupID == "sub-group-6" && (ri.Node.info.GroupSeq == 0 || ri.Node.info.GroupSeq == 1) {
		ps.TerminateService()
	}

	time.Sleep(time.Second)
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
		ps.MyPublish("Portugal won the world cup!", "portugal T/soccer T")
	}

	time.Sleep(3 * time.Second)
	Sync(ctx, ri.RunInfo, finishedState)

	finalMem, err = mem.VirtualMemory()
	if err != nil {
		return err
	}
	finalCpu, err = cpu.Times(false)
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
