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

func FinalScoutTest(runenv *runtime.RunEnv) error {
	commonOpts := GetCommonOpts(runenv)

	ctx, cancel := context.WithTimeout(context.Background(), commonOpts.Timeout)
	defer cancel()

	ri, err := Base(ctx, runenv, commonOpts)
	if err != nil {
		return err
	}

	if err := TestFinalScout(ctx, ri); err != nil {
		return err
	}
	Teardown(ctx, ri.RunInfo)

	return nil
}

func TestFinalScout(ctx context.Context, ri *DHTRunInfo) error {

	replicationFactor := 0
	rFactor := "0"
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

	Sync(ctx, ri.RunInfo, readyState)
	// Begining 1st-Sub

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
	cfg.FaultToleranceFactor = replicationFactor
	ps := pubsub.NewPubSub(ri.Node.dht, cfg)
	ps.SetHasOldPeer()

	Sync(ctx, ri.RunInfo, createdState)

	// Subscribing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		ps.MySubscribe("portugal T/surf T")
	case "sub-group-2":
		ps.MySubscribe("portugal T/soccer T")
	case "sub-group-3":
		ps.MySubscribe("surf T/bali T")
	case "sub-group-4":
		ps.MySubscribe("surf T/bali T/trip T/price R 1000 1500")
	case "sub-group-5":
		ps.MySubscribe("surf T/trip T/price R 1000 2000")
	case "sub-group-6":
		ps.MySubscribe("surf T/trip T/price R 1000 1400")
	case "sub-group-7":
		ps.MySubscribe("benfica T/soccer")
	case "sub-group-8":
		ps.MySubscribe("trip T/price R 1000 2000")
	case "sub-group-9":
		ps.MySubscribe("surf T")
	case "sub-group-10":
		ps.MySubscribe("portugal T/museum T/art T")
	}

	time.Sleep(2 * time.Second)
	Sync(ctx, ri.RunInfo, subbedState)

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("Bali is full of Left handers! Goofie paradise :D", "bali T/surf T")
	case "pub-2":
		ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
	case "pub-3":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	case "pub-4":
		ps.MyPublish("Portugal won the world cup!", "portugal T/soccer T")
	case "pub-5":
		ps.MyPublish("Visit or surf trip website at narlytrips.com with some travels only for 1050$", "surf T/trip T/price R 1050 1050")
	case "pub-6":
		ps.MyPublish("Benfica just signed Messi!", "benfica T/soccer T")
	case "pub-7":
		ps.MyPublish("Cruise trips for only 1570$, only tommorow", "trip T/price R 1570 1570")
	case "pub-8":
		ps.MyPublish("Expo of impessionist art in CCB, next week ", "portugal T/museum T/art T")
	}

	time.Sleep(3 * time.Second)
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
	ps.ReturnCorrectnessStats(nil)
	subs := ps.ReturnSubStats()
	runenv.R().RecordPoint("# Peers - ScoutSubs"+rFactor+"1st", float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
	runenv.R().RecordPoint("CPU used - ScoutSubs"+rFactor+"1st", finalCpu[0].User-initCpu[0].User)
	runenv.R().RecordPoint("Memory used - ScoutSubs"+rFactor+"1st", float64(finalMem.Used-initMem.Used)/(1024*1024))

	for _, ev := range events {
		runenv.R().RecordPoint("Event Latency - ScoutSubs"+rFactor+"1st", float64(ev))
	}
	for _, sb := range subs {
		runenv.R().RecordPoint("Sub Latency - ScoutSubs"+rFactor+"1st", float64(sb))
	}

	Sync(ctx, ri.RunInfo, recordedState)

	// Begining 2nd-Sub
	time.Sleep(time.Second)

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
	case "sub-group-2":
		ps.MySubscribe("surf T/azores T")
	case "sub-group-3":
		ps.MySubscribe("tesla T/stock T/price R 0 400")
	case "sub-group-4":
		ps.MySubscribe("bitcoin T/price R 10000 15000")
	case "sub-group-5":
		ps.MySubscribe("temperature R 30 40")
	case "sub-group-6":
		ps.MySubscribe("ronaldo T/sporting T")
	case "sub-group-7":
		ps.MySubscribe("portugal T/temperature R 30 40")
	case "sub-group-8":
		ps.MySubscribe("italy T/temperature R 30 35")
	case "sub-group-9":
		ps.MySubscribe("laptop T/ram R 32 64/ssd 512 512/price R 1000 1500")
	case "sub-group-10":
		ps.MySubscribe("ronaldo T/sporting T")
	}

	time.Sleep(3 * time.Second)
	Sync(ctx, ri.RunInfo, subbedState)

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("Tesla stock plummeted to 100 bucks after Musk tweet", "tesla T/stock T/price R 100 100")
	case "pub-2":
		ps.MyPublish("Portugal in the next few weeks will have nice waves and temperature might reach 38ºC", "portugal T/surf T/temperature R 38 38")
	case "pub-3":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	case "pub-4":
		ps.MyPublish("Portugal arquipelagos have uncrowded world class waves. Come visit!", "madeira T/azores T/surf T")
	case "pub-5":
		ps.MyPublish("Visit or surf trip website at narlytrips.com with some travels only for 1050$", "surf T/trip T/price R 1050 1050")
	case "pub-6":
		ps.MyPublish("Italy in the next few weeks will have temperatures around 34ºC", "italy T/temperature R 34 34")
	case "pub-7":
		ps.MyPublish("Top gear laptops on 75 per cent discont", "laptop T/ram R 64 64/ssd R 512 512/price R 1100 1100")
	case "pub-8":
		ps.MyPublish("Ronaldo refuses sporting, for benfica is the best portuguese club", "ronaldo T/sporting T")
	}

	time.Sleep(5 * time.Second)
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
	ps.ReturnCorrectnessStats(nil)
	subs = ps.ReturnSubStats()
	runenv.R().RecordPoint("# Peers - ScoutSubs"+rFactor+"2nd", float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
	runenv.R().RecordPoint("CPU used - ScoutSubs"+rFactor+"2nd", finalCpu[0].User-initCpu[0].User)
	runenv.R().RecordPoint("Memory used - ScoutSubs"+rFactor+"2nd", float64(finalMem.Used-initMem.Used)/(1024*1024))

	for _, ev := range events {
		runenv.R().RecordPoint("Event Latency - ScoutSubs"+rFactor+"2nd", float64(ev))
	}
	for _, sb := range subs {
		runenv.R().RecordPoint("Sub Latency - ScoutSubs"+rFactor+"2nd", float64(sb))
	}

	Sync(ctx, ri.RunInfo, recordedState)

	// Begining 3rd-Sub
	time.Sleep(time.Second)

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
		ps.MySubscribe("ipfs T")
	case "sub-group-2":
		ps.MySubscribe("ipfs T")
	case "sub-group-3":
		ps.MySubscribe("ipfs T")
	case "sub-group-4":
		ps.MySubscribe("ipfs T")
	case "sub-group-5":
		ps.MySubscribe("ipfs T")
	case "sub-group-6":
		ps.MySubscribe("ipfs T")
	case "sub-group-7":
		ps.MySubscribe("ipfs T")
	case "sub-group-8":
		ps.MySubscribe("ipfs T")
	case "sub-group-9":
		ps.MySubscribe("ipfs T")
	case "sub-group-10":
		ps.MySubscribe("ipfs T")
	}

	time.Sleep(3 * time.Second)
	Sync(ctx, ri.RunInfo, subbedState)

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("The new content-based IPFS pubsub is already working!", "ipfs T")
	case "pub-2":
		ps.MyPublish("Today the weather in portugal is sunny, with temperatures reaching 32ºC", "portugal T/temperature R 32 32")
	case "pub-3":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	case "pub-4":
		ps.MyPublish("Today a new version of Kademlia DHT was lanched for IPFS!", "ipfs T")
	case "pub-5":
		ps.MyPublish("Benfica sign ronaldo to team up with messi!", "benfica T/soccer T")
	case "pub-6":
		ps.MyPublish("scoutsubs is the best working pub-sub on the web", "ipfs T")
	case "pub-7":
		ps.MyPublish("Big riders charge madeira's sweel of the year at Jardim do mar", "surf T/madeira T")
	case "pub-8":
		ps.MyPublish("Visit narlytrips.com for cheap and wonderfull surf trips!", "surf T/trip T")
	}

	time.Sleep(5 * time.Second)
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
	ps.ReturnCorrectnessStats(nil)
	subs = ps.ReturnSubStats()
	runenv.R().RecordPoint("# Peers - ScoutSubs"+rFactor+"3rd", float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
	runenv.R().RecordPoint("CPU used - ScoutSubs"+rFactor+"3rd", finalCpu[0].User-initCpu[0].User)
	runenv.R().RecordPoint("Memory used - ScoutSubs"+rFactor+"3rd", float64(finalMem.Used-initMem.Used)/(1024*1024))

	for _, ev := range events {
		runenv.R().RecordPoint("Event Latency - ScoutSubs"+rFactor+"3rd", float64(ev))
	}
	for _, sb := range subs {
		runenv.R().RecordPoint("Sub Latency - ScoutSubs"+rFactor+"3rd", float64(sb))
	}

	Sync(ctx, ri.RunInfo, recordedState)

	// Begining Fault scenario
	time.Sleep(time.Second)

	initMem, err = mem.VirtualMemory()
	if err != nil {
		return err
	}
	initCpu, err = cpu.Times(false)
	if err != nil {
		return err
	}

	Sync(ctx, ri.RunInfo, readyState)

	// Subscribe Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		ps.MySubscribe("IST T")
	case "sub-group-2":
		ps.MySubscribe("IST T")
	case "sub-group-3":
		ps.MySubscribe("IST T")
	case "sub-group-4":
		ps.MySubscribe("IST T")
	case "sub-group-5":
		ps.MySubscribe("IST T")
	case "sub-group-6":
		ps.MySubscribe("IST T")
	case "sub-group-7":
		ps.MySubscribe("FCUL T")
	case "sub-group-8":
		ps.MySubscribe("MIT T")
	case "sub-group-9":
		ps.MySubscribe("MIT T")
	case "sub-group-10":
		ps.MySubscribe("Caltech T")
	}

	time.Sleep(3 * time.Second)
	Sync(ctx, ri.RunInfo, subbedState)

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("IPFS has reached 10 million weekly unique users!", "ipfs T")
	case "pub-2":
		ps.MyPublish("Portuguese league surpaced the german league in FIFA ratings", "portugal T/soccer T")
	case "pub-3":
		ps.MyPublish("FCUL is lanching a new research grant for Computer Science Phd Students", "FCUL T")
	case "pub-4":
		ps.MyPublish("IST received the best portuguese freshmen", "IST T")
	case "pub-5":
		ps.MyPublish("Visit narlytrips.com for cheap and wonderfull surf trips!", "surf T/trip T")
	case "pub-6":
		ps.MyPublish("MIT start colaboration with IST!", "IST T/MIT T")
	case "pub-7":
		ps.MyPublish("Caltech shows the world the first quantic comuter!", "Caltech T")
	case "pub-8":
		ps.MyPublish("Benfica sells expendables for 250 M€!", "benfica T/soccer T")
	}

	time.Sleep(5 * time.Second)
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
	ps.ReturnCorrectnessStats(nil)
	subs = ps.ReturnSubStats()
	runenv.R().RecordPoint("# Peers - ScoutSubs"+rFactor+"4th", float64(len(ri.Node.dht.RoutingTable().GetPeerInfos())))
	runenv.RecordMessage("GroupID >> " + ri.RunInfo.RunEnv.RunParams.TestGroupID)
	runenv.R().RecordPoint("CPU used - ScoutSubs"+rFactor+"4th", finalCpu[0].User-initCpu[0].User)
	runenv.R().RecordPoint("Memory used - ScoutSubs"+rFactor+"4th", float64(finalMem.Used-initMem.Used)/(1024*1024))

	for _, ev := range events {
		runenv.R().RecordPoint("Event Latency - ScoutSubs"+rFactor+"4th", float64(ev))
	}
	for _, sb := range subs {
		runenv.R().RecordPoint("Sub Latency - ScoutSubs"+rFactor+"4th", float64(sb))
	}

	Sync(ctx, ri.RunInfo, recordedState)

	if err := stager.End(); err != nil {
		return err
	}

	return nil
}
