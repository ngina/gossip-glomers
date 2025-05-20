## Gossip Glomers Challenge
My solutions to the Gossip Glomers Challenge: a series of distributed systems challenges.

## Challenge #1: Echo 
This challenge mainly includes setting up the repo.

## Challenge #2: Unique ID Generator
TODO
 
## Challenge #3: Broadcast

### Challenge #3b: Multi-Node Broadcast

### Challenge #3b: Multi-Node Broadcast

#### What is the high-level goal of this challenge?
Make sure all nodes eventually see all broadcasted values, using reliable propagation across a cluster without network partitions.

#### My approach:
1. Save the list of neighbors from the topology message: Each node receives a "topology" messsage as the initial message. I parsed the neighbors of the given node and saved it to a list.

2. Broadcasting to all neighbors except the sender of the message: When a node receives a "broadcast" message, I asynchronously broadcast using Send(dest, responseBody) to all neighbors except the node that matches the sender, if it exists.

#### Results:
Tests passed but I notice repeated values in the node state ie `node.values`. This suggests messages are being sent and re-sent to neighboring nodes. I'm thinking some dedupe logic is needed here but yet to figure out how to fix it.

Viewing Network Requests/Responses: Each test run will generate a svg of all the network requests sent. You can view this at `store/latest/messages.svg`. Found it very helpful to see the network requests and responses to nodes, which helped me detect bugs much easier.

#### Next Steps:
1. Fix repeated node values.


### References
* DDIA: Ch 9. Consistency and Consensus
* Ordering using [Lamport Timestamps & Total Order Broadcast](https://youtu.be/yIvft09RTAg?si=1eY4InG_y6SKnDxJhttps://youtu.be/yIvft09RTAg?si=1eY4InG_y6SKnDxJ)