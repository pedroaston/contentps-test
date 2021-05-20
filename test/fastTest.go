package test

import (
	"context"

	"github.com/testground/sdk-go/runtime"
)

func FastTest(runenv *runtime.RunEnv) error {
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

func TestFastDelivery(ctx context.Context, ri *DHTRunInfo) error {
	/*
		runenv := ri.RunEnv
		readyState := sync.State("ready")
		createdState := sync.State("created")
		advertisedState := sync.State("advertised")
		subbedState := sync.State("subbed")
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

		ps := pubsub.NewPubSub(ri.Node.dht, "PT")

		ri.Client.MustSignalEntry(ctx, createdState)
		err1stStop := <-ri.Client.MustBarrier(ctx, createdState, runenv.TestInstanceCount).C
		if err1stStop != nil {
			return err1stStop
		}

		// Create and Advertise Routine
		switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
		case "pub-1":
			ps.CreateMulticastGroup("ipfs T")
		case "pub-2":
			ps.CreateMulticastGroup("portugal T/surf T")
		case "pub-3":
			ps.CreateMulticastGroup("surf T/bali T/trip T/price R 1000 2000")
		}

		time.Sleep(time.Second)

		ri.Client.MustSignalEntry(ctx, advertisedState)
		err2ndStop := <-ri.Client.MustBarrier(ctx, advertisedState, runenv.TestInstanceCount).C
		if err2ndStop != nil {
			return err2ndStop
		}

		// Subscribing Routine
		switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
		case "sub-group-1":
			ps.MySearchAndPremiumSub("surf T")
			ps.MySearchAndPremiumSub("ipfs T")
		case "sub-group-2":
			ps.MySearchAndPremiumSub("ipfs T")
			ps.MySearchAndPremiumSub("portugal T/surf T")
		case "sub-group-3":
			ps.MySearchAndPremiumSub("ipfs T")
			ps.MySearchAndPremiumSub("surf T/bali T/trip T/price R 1000 1500")
		case "sub-group-4":
			ps.MySearchAndPremiumSub("ipfs T")
			ps.MySearchAndPremiumSub("surf T/bali T/trip T/price R 1000 1500")
		case "sub-group-5":
			ps.MySearchAndPremiumSub("ipfs T")
			ps.MySearchAndPremiumSub("surf T/trip T/price R 1500 2000")
		case "sub-group-6":
			ps.MySearchAndPremiumSub("ipfs T")
			ps.MySearchAndPremiumSub("surf T/trip T/price R 1000 1400")
		}

		time.Sleep(time.Second)
		ri.Client.MustSignalEntry(ctx, subbedState)
		err3rdStop := <-ri.Client.MustBarrier(ctx, subbedState, runenv.TestInstanceCount).C
		if err3rdStop != nil {
			return err3rdStop
		}

		// Publishing Routine
		switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
		case "pub-1":
			ps.MyPremiumPublish("ipfs T", "Publishing via ipfs is lit!", "ipfs T")
		case "pub-2":
			ps.MyPremiumPublish("portugal T/surf T", "Portugal has the world's best waves!", "portugal T/surf T")
		case "pub-3":
			ps.MyPremiumPublish("surf T/bali T/trip T/price R 1000 2000", "Surf trip to bali for 1200€, only today!", "surf T/bali T/trip T/price R 1200 1200")
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

		_, nEFast, _, latFast := ps.ReturnReceivedEventsStats()
		runenv.R().RecordPoint("Number of peers", float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
		runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
		runenv.R().RecordPoint("Avg time to subscribe", float64(0))
		runenv.R().RecordPoint("Number of events received via FastDelivery", float64(nEFast))
		runenv.R().RecordPoint("Avg latency of events received via FastDelivery", float64(latFast))
		runenv.R().RecordPoint("Avg time to sub via ScoutSubs", float64(ps.ReturnSubStats()))
		runenv.R().RecordPoint("CPU used - Begin", initCpu[0].User)
		runenv.R().RecordPoint("CPU used - End", finalCpu[0].User)
		runenv.R().RecordPoint("Memory used - Begin", float64(initMem.Used))
		runenv.R().RecordPoint("Memory used - End", float64(finalMem.Used))

		ri.Client.MustSignalEntry(ctx, recordedState)
		err5thStop := <-ri.Client.MustBarrier(ctx, recordedState, runenv.TestInstanceCount).C
		if err5thStop != nil {
			return err5thStop
		}

		if err := stager.End(); err != nil {
			return err
		}
	*/
	return nil
}
