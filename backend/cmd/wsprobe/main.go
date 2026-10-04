// Command wsprobe is a terminal WebSocket client for demonstrating live events.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 3 {
		log.Fatal("usage: wsprobe <token> <topic> [topic...]")
	}
	token, topics := os.Args[1], os.Args[2:]

	addr := os.Getenv("WS_URL")
	if addr == "" {
		addr = "ws://localhost:8090/ws"
	}

	conn, _, err := websocket.DefaultDialer.Dial(addr+"?token="+token, nil)
	if err != nil {
		log.Fatalf("dial %s: %v", addr, err)
	}
	defer conn.Close()

	for _, topic := range topics {
		msg, _ := json.Marshal(map[string]string{"action": "subscribe", "topic": topic})
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Fatalf("subscribe %s: %v", topic, err)
		}
	}
	fmt.Printf("subscribed to %s, waiting for events (Ctrl-C to quit)\n", strings.Join(topics, ", "))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-interrupt
		conn.Close()
	}()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var event struct {
			Type    string          `json:"type"`
			Topic   string          `json:"topic"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal(raw, &event); err != nil {
			continue
		}
		fmt.Printf("\n[%s] on %s\n%s\n", event.Type, event.Topic, indent(event.Payload))
	}
}

func indent(raw json.RawMessage) string {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetIndent("  ", "  ")
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	_ = enc.Encode(v)
	return "  " + strings.TrimRight(buf.String(), "\n")
}
