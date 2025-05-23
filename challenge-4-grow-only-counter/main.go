package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const counterKey string = "counter"

type server struct {
	n   *maelstrom.Node
	kv  *maelstrom.KV
	ctx context.Context
}

func main() {
	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)

	s := &server{
		n:   n,
		kv:  kv,
		ctx: context.Background(),
	}
	log.Printf("running node: %v", s.n)
	s.n.Handle("add", s.incrementHandler)
	s.n.Handle("read", s.readHandler)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *server) incrementHandler(req maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}
	incrBy := body["delta"].(float64)
	for {
		old, err := s.kv.ReadInt(s.ctx, counterKey)
		if err != nil {
			if maelstrom.ErrorCode(err) == maelstrom.KeyDoesNotExist {
				old = 0
			} else {
				return err
			}
		}

		newVal := old + int(incrBy)
		err = s.kv.CompareAndSwap(s.ctx, counterKey, old, newVal, true)
		if err == nil { // successfully updated value
			break
		} else if maelstrom.ErrorCode(err) != maelstrom.PreconditionFailed { // failure not related to precondition
			return err
		}
		time.Sleep(30 * time.Millisecond)
	}

	resBody := map[string]any{
		"type": "add_ok",
	}
	return s.n.Reply(req, resBody)
}

func (s *server) readHandler(req maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}

	v, err := s.kv.ReadInt(s.ctx, counterKey)
	if err != nil {
		log.Printf("%s", err)
		if maelstrom.ErrorCode(err) != maelstrom.KeyDoesNotExist {
			return err
		}
	}

	resBody := map[string]any{
		"type":        "read_ok",
		"value":       int(v),
		"msg_id":      body["msg_id"],
		"in_reply_to": body["msg_id"],
	}
	return s.n.Reply(req, resBody)
}
