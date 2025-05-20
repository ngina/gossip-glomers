package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type server struct {
	n      *maelstrom.Node
	values []any
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

	responseBody := map[string]any{
		"type":        "broadcast_ok",
		"msg_id":      body["msg_id"],
		"in_reply_to": body["msg_id"],
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

	responseBody := map[string]any{
		"type":        "topology_ok",
		"msg_id":      body["msg_id"],
		"in_reply_to": body["msg_id"],
	}
	return s.n.Reply(msg, responseBody)
}
