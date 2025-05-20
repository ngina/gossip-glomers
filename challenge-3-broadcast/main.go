package main

import (
	"encoding/json"
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type server struct {
	n      *maelstrom.Node
	values []any
	neighbors []any
}

func main() {
	n := maelstrom.NewNode()
	s := &server{n: n}
	n.Handle("broadcast", s.broadcast)
	n.Handle("read", s.read)
	n.Handle("topology", s.topology)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *server) broadcast(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	s.values = append(s.values, body["message"])
	msg_id := body["msg_id"]

	// Asynchronously replicate message to neighbor nodes
	for _, neighbor := range s.neighbors {
		if msg.Src != neighbor {
			sendMsgBody := map[string]any{
				"type": "broadcast",
				"message": body["message"],
				"msg_id":  msg_id,
			}
			err := s.n.Send(neighbor.(string), sendMsgBody)
			if err != nil {
				return err
			}
		}
	}

	// Respond to request
	responseBody := map[string]any{
		"type":        "broadcast_ok",
		"msg_id":      msg_id,
		"in_reply_to": msg_id,
	}
	return s.n.Reply(msg, responseBody)
}

func (s *server) read(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	responseBody := map[string]any{
		"type":        "read_ok",
		"messages":    s.values,
		"msg_id":      body["msg_id"],
		"in_reply_to": body["msg_id"],
	}
	return s.n.Reply(msg, responseBody)
}

func (s *server) topology(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	topologyRaw, ok := body["topology"].(map[string]any)
	if !ok {
		log.Println("topology is not a map[string]interface{}")
		return fmt.Errorf("topology is not a map[string]interface{}")
	}
	neighborsRaw, ok := topologyRaw[s.n.ID()].([]any)
	if !ok {
		log.Printf("topology[%s] is not a []interface{}", s.n.ID())
		return fmt.Errorf("topology[%s] is not a []interface{}", s.n.ID())
	}
	s.neighbors = append(s.neighbors, neighborsRaw...)

	responseBody := map[string]any{
		"type":        "topology_ok",
		"msg_id":      body["msg_id"],
		"in_reply_to": body["msg_id"],
	}
	return s.n.Reply(msg, responseBody)
}
