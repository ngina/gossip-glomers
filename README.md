## Gossip Glomers Challenge

My solutions to the Gossip Glomers Challenge: a series of distributed systems challenges.

---

## Challenge #1: Echo

This challenge mainly includes setting up the repo.

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

- **Tests pass without duplicate values**
Previously, tests passed but nodes contained repeated broadcast messages, suggesting messages were being re-broadcast without deduplication. Fixed this by introducing **deduplication logic**, using the message as the deduplication key to ensure each value is stored only once.

~~Tests passed but I noticed repeated values in the node state (i.e., `node.values`). This suggested that messages were being re-sent to neighboring nodes. I suspected deduplication logic was needed but hadn’t implemented it yet.~~

#### Next Steps
- [ ] Make `message_id` unique: Currently the same `message_id` is reused across broadcasts from different nodes.
- ~~[x] Fix repeated node values.~~

---

## Debugging Network Requests
**Viewing Network Requests/Responses:** Each test run generates an SVG of all network requests sent. You can view it at `store/latest/messages.svg`.This visualization helped me detect bugs more easily and understand propagation patterns.

---

### References

- *Designing Data-Intensive Applications* (DDIA), Ch. 9: Consistency and Consensus  
- [Ordering with Lamport Timestamps & Total Order Broadcast](https://youtu.be/yIvft09RTAg?si=1eY4InG_y6SKnDxJ)

---