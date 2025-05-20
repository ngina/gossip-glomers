package main

import (
	"encoding/json"
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type server struct {
    n *maelstrom.Node
    id int
}

func main() {
	n := maelstrom.NewNode()
    s := &server{n: n}
    n.Handle("generate", s.generate)

    if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *server) generate_unique_id() string {
	s.id = s.id + 1
	return fmt.Sprintf("%v%v", s.n.ID(), s.id)
}

func (s *server) generate(msg maelstrom.Message) error {
    var body map[string]any
    if err := json.Unmarshal(msg.Body, &body); err != nil {
        return err
    }

    id := s.generate_unique_id()
    body["type"] = "generate_ok"
    body["id"] = id

    return s.n.Reply(msg, body)
}

