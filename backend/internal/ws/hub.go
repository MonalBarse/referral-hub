// Package ws is the real-time layer: an in-process topic hub plus socket plumbing.
package ws

import (
	"encoding/json"
	"log"
	"regexp"
	"sync"
	"time"
)

const (
	EventJobCreated      = "job.created"
	EventReferralCreated = "referral.created"
	EventCommentCreated  = "comment.created"
)

const TopicJobs = "jobs"

var topicPattern = regexp.MustCompile(`^(jobs|job:[0-9a-fA-F-]{36}|referral:[0-9a-fA-F-]{36})$`)

func ValidTopic(name string) bool { return topicPattern.MatchString(name) }

type Event struct {
	Type      string    `json:"type"`
	Topic     string    `json:"topic"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

type Hub struct {
	mu     sync.RWMutex
	topics map[string]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{topics: make(map[string]map[*Client]struct{})}
}

func (h *Hub) Publish(topic, eventType string, payload any) {
	frame, err := json.Marshal(Event{
		Type:      eventType,
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("ws: marshal %s: %v", eventType, err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.topics[topic] {
		select {
		case c.send <- frame:
		default:
			// Never block the fan-out on one slow client.
			log.Printf("ws: dropping %s for a slow client", eventType)
		}
	}
}

func (h *Hub) SubscriberCount(topic string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.topics[topic])
}

func (h *Hub) subscribe(c *Client, topic string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.topics[topic] == nil {
		h.topics[topic] = make(map[*Client]struct{})
	}
	h.topics[topic][c] = struct{}{}
}

func (h *Hub) unsubscribe(c *Client, topic string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.removeLocked(c, topic)
}

func (h *Hub) remove(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for topic := range h.topics {
		h.removeLocked(c, topic)
	}
}

func (h *Hub) removeLocked(c *Client, topic string) {
	subs, ok := h.topics[topic]
	if !ok {
		return
	}
	delete(subs, c)
	if len(subs) == 0 {
		delete(h.topics, topic)
	}
}
