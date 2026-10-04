package ws

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(h *Hub) *Client {
	return &Client{hub: h, send: make(chan []byte, sendBuffer)}
}

func TestPublishReachesOnlySubscribers(t *testing.T) {
	hub := NewHub()

	subscriber := newTestClient(hub)
	bystander := newTestClient(hub)

	hub.subscribe(subscriber, TopicJobs)
	hub.subscribe(bystander, "referral:11111111-1111-1111-1111-111111111111")

	hub.Publish(TopicJobs, EventJobCreated, map[string]string{"title": "Backend Engineer"})

	require.Len(t, subscriber.send, 1, "subscriber should have received the event")
	assert.Empty(t, bystander.send, "a client on another topic must not receive it")

	var event Event
	require.NoError(t, json.Unmarshal(<-subscriber.send, &event))
	assert.Equal(t, EventJobCreated, event.Type)
	assert.Equal(t, TopicJobs, event.Topic)
	assert.False(t, event.Timestamp.IsZero())

	payload, ok := event.Payload.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Backend Engineer", payload["title"])
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	hub := NewHub()
	c := newTestClient(hub)

	hub.subscribe(c, TopicJobs)
	hub.unsubscribe(c, TopicJobs)
	hub.Publish(TopicJobs, EventJobCreated, nil)

	assert.Empty(t, c.send)
	assert.Equal(t, 0, hub.SubscriberCount(TopicJobs))
}

func TestRemoveDropsClientFromEveryTopic(t *testing.T) {
	hub := NewHub()
	c := newTestClient(hub)
	referral := "referral:11111111-1111-1111-1111-111111111111"

	hub.subscribe(c, TopicJobs)
	hub.subscribe(c, referral)
	hub.remove(c)

	assert.Equal(t, 0, hub.SubscriberCount(TopicJobs))
	assert.Equal(t, 0, hub.SubscriberCount(referral))
}

func TestPublishDropsFramesForASlowClient(t *testing.T) {
	hub := NewHub()
	slow := newTestClient(hub)
	hub.subscribe(slow, TopicJobs)

	for i := 0; i < sendBuffer+10; i++ {
		hub.Publish(TopicJobs, EventJobCreated, map[string]int{"n": i})
	}

	assert.Len(t, slow.send, sendBuffer, "buffer should be full, with the excess dropped")
}

func TestValidTopic(t *testing.T) {
	valid := []string{
		"jobs",
		"job:11111111-1111-1111-1111-111111111111",
		"referral:11111111-1111-1111-1111-111111111111",
	}
	for _, topic := range valid {
		assert.True(t, ValidTopic(topic), "%q should be allowed", topic)
	}

	invalid := []string{"", "everything", "job:not-a-uuid", "jobs; DROP TABLE", "referral:*"}
	for _, topic := range invalid {
		assert.False(t, ValidTopic(topic), "%q should be rejected", topic)
	}
}
