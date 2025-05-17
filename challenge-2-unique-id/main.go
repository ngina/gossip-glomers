package main

import (
	"encoding/json"
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
    s := &Server{n: n}
    n.Handle("generate", s.generate)

    if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

var id = 0
type Server struct {
    n *maelstrom.Node
}

func (s *Server) generate_unique_id() string {
	id = id + 1
	return fmt.Sprintf("%v%v", s.n.ID(), id)
}

func (s *Server) generate(msg maelstrom.Message) error {
    var body map[string]any
    if err := json.Unmarshal(msg.Body, &body); err != nil {
        return err
    }

    id := s.generate_unique_id()
    body["type"] = "generate_ok"
    body["id"] = id

    return s.n.Reply(msg, body)
}

