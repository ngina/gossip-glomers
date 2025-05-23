## Gossip Glomers Challenge

My solutions to the [Gossip Glomers](https://fly.io/dist-sys/) Challenge: a series of distributed systems challenges.

---

## Challenge #1: Echo

This challenge mainly includes setting up the repo.

### Setup 
- Create and init go directory, and create `main.go` file
```
mkdir challenge-1-echo
cd challenge-1-echo && go mod init maelstrom-echo
touch main.go
```

- Add implementation logic to `main.go`. Remember to import maelstrom package, if needed.
```
import (
   ...

   maelstrom github.com/jepsen-io/maelstrom/demo/go // refers to maelstrom package
)
```

- To compile our program, fetch the Maelstrom library and install. This will build the maelstrom-echo binary and place it in your $GOBIN path which is typically ~/go/bin.
```
go get github.com/jepsen-io/maelstrom/demo/go
go install .
```

- Run node in Maelstrom: Start up Maelstrom and pass path to binary. This command instructs maelstrom to run the "echo" workload against our binary. It runs a single node and it will send "echo" commands for 10 seconds.
```
./maelstrom test -w echo --bin ~/go/bin/maelstrom-echo --node-count 1 --time-limit 10

```

---

## Challenge #2: Unique ID Generator

TODO

---

## Challenge #3: Broadcast

### Challenge #3b: Single-Node Broadcast

#### High-Level Goal

Implement all three handlers: `broadcast`, `read`, and `topology`.

---

### Challenge #3b: Multi-Node Broadcast

#### High-Level Goal

Make sure all nodes eventually see all broadcasted values, using reliable propagation across a cluster without network partitions.

#### My Approach

1. **Save neighbors from the topology message**  
   Each node receives a `"topology"` message as the initial message. I parsed the neighbors of the current node and stored them in a map (which is used as a set).

2. **Broadcast to all neighbors except the sender**  
   When a node receives a `"broadcast"` message, it sends the message to all its neighbors using `Send(dest, responseBody)` to all neighbors except the node that matches the sender, if it exists.

#### Results

- **Tests pass without duplicate values:**
Previously, tests passed but nodes contained repeated broadcast message values, suggesting messages were being re-broadcast without deduplication. Fixed this by introducing deduplication logic, using the message as the deduplication key to ensure each value is stored only once.

~~ - Tests passed but I noticed repeated values in the node state (i.e., `node.values`). This suggested that messages were being re-sent to neighboring nodes. I suspected deduplication logic was needed but hadn’t implemented it yet.~~

#### Next Steps
- [ ] Make `message_id` unique: Currently the same `message_id` is reused across broadcasts from different nodes.
- ~~[x] Fix repeated node values.~~

---

### Challenge #4: Grow-Only Counter

For this challenge, we can leverage the sequentially-consistent key value store to implement the grow only counter.
We can do this because:
```
...there are no concurrent operations in a linearizable datastore: there must be a single timeline along which all operations are totally ordered. There might be several requests waiting to be handled, but the data store ensures that every request is handled atomically at a single point in time, acting on a single copy of the data, along a single timeline, without any concurrency.

~DDIA Ch.9 pg 341
```

#### My Approach

1. **For each add request:** fetch the latest counter value, modify it, then use compare-and-swap to update the counter if the last read value hasn't changed; otherwise, retry until the operation succeeds.

2. **For each read request**, read latest value and return.

---

## Debugging Network Requests
**Viewing Network Requests/Responses:** 
Each test run generates an SVG of all network requests sent. You can view it at `store/latest/messages.svg`.This visualization helped me detect bugs more easily and understand propagation patterns.

---

### References

- *Designing Data-Intensive Applications* (DDIA), Ch. 9: Consistency and Consensus  
- [Ordering with Lamport Timestamps & Total Order Broadcast](https://youtu.be/yIvft09RTAg?si=1eY4InG_y6SKnDxJ)

---