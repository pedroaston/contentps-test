# Content Pub-Sub : Testground Plan
This repo has a test-plan to evaluate the performance of a pubsub middleware over ipfs kad-dht. This repository was adapted from [libp2p/test-plans](https://github.com/libp2p/test-plans/tree/master/dht).For all details about the project check [here](https://github.com/pedroaston/smartpubsub-ipfs).

## Run Test Environment
Use the files at composition.
'''go
cd composition
testground run composition -f 30-normal-test.toml
'''

## Test Scenarios

NOTE: Due to the time limitations of my work I've chosen to concentrate in preparing the system to a non-byzantine fault scenario, focusing primarily in the system's efficiency, scalability, reliability, fault tolerance and flexibility. 

### Normal Scenario
- Test Case: normal-scout-test
- Goal: Understand the how the system behaves when a small group of publishers are publishing events on the system.

### Subscription Burst Scenario
- Test Case: subburst-scout-test
- Goal: Similar to the normal scenario but this time when publishers forward the events, system will be flooded with subscriptions requests, that are not relevant for the delivery of those events but will load the system anyway.

### Event Burst Scenario
- Test Case: eventburst-scout-test
- Goal: In this scenario the subscription routine is similar to the ones above, but this time all subscribers are also publishing. This means there are 10 times (in the 30 nodes scenario) more events being published. This allow us to understand how the system copes with high event's load.

### Fault Tolerance Scenario
- Test Case: fault-scout-test
- Goal: In this scenario we only demonstrate that our system can tolerate up to f localized faults. This means that the system only fails if f+1 consecutive nodes by ID fail between heartbeat cycles. Both f and heartbeat period are configurable.

### Long Run Scenario
- Test Case: longrun-scout-test
- Goal: Here we demonstarte the refreshing protocol of our system. We have a initial round like in the normal scenario, and then some subs unsubscribe. After a complete refresing cycle occurs, new events are published and only forwarded to the peers that didn't unsubscribe.