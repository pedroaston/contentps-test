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

	variant := "RR"
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
		expectedE = []string{"Portugal has the world's best waves!"}
	case "sub-group-2":
		expectedE = []string{"Portugal won the world cup!"}
	case "sub-group-3":
		expectedE = []string{"Surf trip to bali for 1050, just today!"}
	case "sub-group-4":
		expectedE = []string{"Surf trip to bali for 1050, just today!"}
	case "sub-group-5":
		expectedE = []string{"Hawai surf vavations for 1500!"}
	case "sub-group-6":
		expectedE = []string{"ipfs has just hit 100 millions of unique user!"}
	case "sub-group-7":
		expectedE = []string{"Tesla stocks drop to 10 in the lastest crash!"}
	case "sub-group-8":
		expectedE = []string{"Benfica won against bayern 10-0!"}
	case "sub-group-9":
		expectedE = []string{"Bruno Mars will come to Portugal in June"}
	case "sub-group-10":
		expectedE = []string{"Filecoin is now the most used crypto!"}
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
	cfg.ConcurrentProcessingFactor = 1000
	cfg.RPCTimeout = 20 * time.Second
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
		ps.MySubscribe("surf T/trip T/price R 1200 2000")
	case "sub-group-6":
		ps.MySubscribe("ipfs T")
	case "sub-group-7":
		ps.MySubscribe("tesla T/stock T/price R 0 200")
	case "sub-group-8":
		ps.MySubscribe("soccer T/benfica T")
	case "sub-group-9":
		ps.MySubscribe("bruno mars T")
	case "sub-group-10":
		ps.MySubscribe("filecoin T")
	}

	time.Sleep(5 * time.Second)
	Sync(ctx, ri.RunInfo, subbedState)

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("Publishing via ipfs is lit!", "sporting T")
	case "pub-2":
		ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
	case "pub-3":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	case "pub-4":
		ps.MyPublish("Portugal won the world cup!", "portugal T/soccer T")
	case "pub-5":
		ps.MyPublish("Hawai surf vavations for 1500!", "surf T/hawai T/trip T/price R 1500 1500")
	case "pub-6":
		ps.MyPublish("Tesla stocks drop to 10 in the lastest crash!", "tesla T/stock T/price R 10 10")
	case "pub-7":
		ps.MyPublish("Bruno Mars will come to Portugal in June", "portugal T/bruno mars T")
	case "pub-8":
		ps.MyPublish("Filecoin is now the most used crypto!", "filecoin T")
	case "pub-9":
		ps.MyPublish("Benfica won against bayern 10-0!", "benfica T/soccer T")
	case "pub-10":
		ps.MyPublish("ipfs has just hit 100 millions of unique user!", "ipfs T")
	}

	time.Sleep(30 * time.Second)
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
		ps.MySubscribe("surf T/bali T")
	case "sub-group-4":
		ps.MySubscribe("bitcoin T/price R 10000 15000")
	case "sub-group-5":
		ps.MySubscribe("temperature R 30 40")
	case "sub-group-6":
		ps.MySubscribe("ronaldo T/sporting T")
	case "sub-group-7":
		ps.MySubscribe("dogecoin T")
	case "sub-group-8":
		ps.MySubscribe("gta T")
	case "sub-group-9":
		ps.MySubscribe("assassin creed T")
	case "sub-group-10":
		ps.MySubscribe("pc T/ram R 32 32")
	}

	// Publishing Routine
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "pub-1":
		ps.MyPublish("Publishing via ipfs is lit!", "sporting T")
	case "pub-2":
		ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
	case "pub-3":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	case "pub-4":
		ps.MyPublish("Portugal won the world cup!", "portugal T/soccer T")
	case "pub-5":
		ps.MyPublish("Hawai surf vavations for 1500!", "surf T/hawai T/trip T/price R 1500 1500")
	case "pub-6":
		ps.MyPublish("Tesla stocks drop to 10 in the lastest crash!", "tesla T/stock T/price R 10 10")
	case "pub-7":
		ps.MyPublish("Bruno Mars will come to Portugal in June", "portugal T/bruno mars T")
	case "pub-8":
		ps.MyPublish("Filecoin is now the most used crypto!", "filecoin T")
	case "pub-9":
		ps.MyPublish("Benfica won against bayern 10-0!", "benfica T/soccer T")
	case "pub-10":
		ps.MyPublish("ipfs has just hit 100 millions of unique user!", "ipfs T")
	}

	time.Sleep(30 * time.Second)
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
	time.Sleep(time.Second)

	groupSMan := 9
	// Expected events
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		prefixes := []string{"I already surfed 1"}
		sufixes := []string{" portuguese beaches!"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
	case "sub-group-2":
		prefixes := []string{"Portugal will score 1"}
		sufixes := []string{" goals at the world cup!"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
	case "sub-group-3":
		prefixes := []string{"Surf trip to bali for 11"}
		sufixes := []string{", just today!"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
	case "sub-group-4":
		prefixes := []string{"Surf trip to bali for 11"}
		sufixes := []string{", just today!"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
	case "sub-group-5":
		prefixes := []string{"Surf trip to hawai for 16"}
		sufixes := []string{", just today!"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
	case "sub-group-6":
		prefixes := []string{"Using IPFS, is 1"}
		sufixes := []string{" times cooler than flying-cars!"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
	case "sub-group-7":
		prefixes := []string{"Guys i want to sell tesla stocks for 1"}
		sufixes := []string{" !"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
	case "sub-group-8":
		prefixes := []string{"I love benfica 1"}
		sufixes := []string{" more than my kids!"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
	case "sub-group-9":
		prefixes := []string{"Im seling bruno mars tickets for 1"}
		sufixes := []string{", only accept filecoins!"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
	case "sub-group-10":
		prefixes := []string{"Im giving 1"}
		sufixes := []string{" filecoins to who can finish my thesis!"}
		groupSizes := []int{groupSMan}
		expectedE = ExpectedEvents(prefixes, sufixes, groupSizes)
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
		event3 := fmt.Sprintf("Sporting stadium is 1%d times uglier than benfica's", ri.Node.info.GroupSeq)
		ps.MyPublish(event3, "sporting T")
	case "sub-group-4":
		event4 := fmt.Sprintf("Portugal will score 1%d goals at the world cup!", ri.Node.info.GroupSeq)
		ps.MyPublish(event4, "portugal T/soccer T")
	case "sub-group-5":
		event2 := fmt.Sprintf("Using IPFS, is 1%d times cooler than flying-cars!", ri.Node.info.GroupSeq)
		ps.MyPublish(event2, "ipfs T")
	case "sub-group-6":
		event1 := fmt.Sprintf("I already surfed 1%d portuguese beaches!", ri.Node.info.GroupSeq)
		ps.MyPublish(event1, "portugal T/surf T")
	case "sub-group-7":
		event1 := fmt.Sprintf("I love benfica %d more than my kids!", ri.Node.info.GroupSeq)
		ps.MyPublish(event1, "soccer T/benfica T")
	case "sub-group-8":
		event5 := fmt.Sprintf("Guys i want to sell tesla stocks for 1%d !", ri.Node.info.GroupSeq)
		pred5 := fmt.Sprintf("tesla T/stock T/price R 1%d 1%d", ri.Node.info.GroupSeq, ri.Node.info.GroupSeq)
		ps.MyPublish(event5, pred5)
	case "sub-group-9":
		event1 := fmt.Sprintf("Im giving 1%d filecoins to who can finish my thesis!", ri.Node.info.GroupSeq)
		ps.MyPublish(event1, "filecoin T")
	case "sub-group-10":
		event1 := fmt.Sprintf("Im seling bruno mars tickets for 1%d, only accept filecoins!", ri.Node.info.GroupSeq)
		ps.MyPublish(event1, "bruno mars T")
	}

	time.Sleep(40 * time.Second)
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
	time.Sleep(time.Second)

	// Expected events
	switch ri.RunInfo.RunEnv.RunParams.TestGroupID {
	case "sub-group-1":
		expectedE = []string{"Portugal has the world's best waves!"}
	case "sub-group-2":
		expectedE = []string{"Portugal won the world cup!"}
	case "sub-group-3":
		expectedE = []string{"Surf trip to bali for 1050, just today!"}
	case "sub-group-4":
		expectedE = []string{"Surf trip to bali for 1050, just today!"}
	case "sub-group-5":
		expectedE = []string{"Hawai surf vavations for 1500!"}
	case "sub-group-6":
		expectedE = []string{"ipfs has just hit 100 millions of unique user!"}
	case "sub-group-7":
		expectedE = []string{"Tesla stocks drop to 10 in the lastest crash!"}
	case "sub-group-8":
		expectedE = []string{"Benfica won against bayern 10-0!"}
	case "sub-group-9":
		expectedE = []string{"Bruno Mars will come to Portugal in June"}
	case "sub-group-10":
		expectedE = []string{"Filecoin is now the most used crypto!"}
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
		ps.MyPublish("Publishing via ipfs is lit!", "sporting T")
	case "pub-2":
		ps.MyPublish("Portugal has the world's best waves!", "portugal T/surf T")
	case "pub-3":
		ps.MyPublish("Surf trip to bali for 1050, just today!", "surf T/bali T/trip T/price R 1050 1050")
	case "pub-4":
		ps.MyPublish("Portugal won the world cup!", "portugal T/soccer T")
	case "pub-5":
		ps.MyPublish("Hawai surf vavations for 1500!", "surf T/hawai T/trip T/price R 1500 1500")
	case "pub-6":
		ps.MyPublish("Tesla stocks drop to 10 in the lastest crash!", "tesla T/stock T/price R 10 10")
	case "pub-7":
		ps.MyPublish("Bruno Mars will come to Portugal in June", "portugal T/bruno mars T")
	case "pub-8":
		ps.MyPublish("Filecoin is now the most used crypto!", "filecoin T")
	case "pub-9":
		ps.MyPublish("Benfica won against bayern 10-0!", "benfica T/soccer T")
	case "pub-10":
		ps.MyPublish("ipfs has just hit 100 millions of unique user!", "ipfs T")
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
