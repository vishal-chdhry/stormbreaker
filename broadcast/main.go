package main

import (
	"encoding/json"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	store := make([]int, 0)

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		body := struct {
			T       string `json:"type"`
			Message int    `json:"message"`
		}{}
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		store = append(store, body.Message)
		response := map[string]interface{}{
			"type": "broadcast_ok",
		}
		return n.Reply(msg, response)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "read_ok"
		body["messages"] = store
		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// TODO: store topology
		body["type"] = "topology_ok"
		delete(body, "topology")
		return n.Reply(msg, body)
	})
	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
