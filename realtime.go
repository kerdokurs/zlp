package zlp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
)

type EventType string

const (
	EventTypeMessage      EventType = "messages"
	EventTypeSubscription EventType = "subscriptions"
	EventTypeRealmUser    EventType = "realm_user"
)

type Event struct {
	ApplyMarkdown bool `json:"apply_markdown"`
}

type EventQueue struct {
	Id          string `json:"queue_id"`
	LastEventId int    `json:"last_event_id"`
}

func (b *Bot) RegisterEvent(types ...EventType) *EventQueue {
	values := url.Values{}
	typesBytes, err := json.Marshal(types)
	if err != nil {
		log.Printf("Could not marshal types: %v\n", err)
		return nil
	}
	values.Set("event_types", string(typesBytes))

	bytes, err := b.getResponseData("POST", "register", &values)
	if err != nil {
		log.Printf("Could not construct request: %v\n", err)
		return nil
	}

	var queue EventQueue
	if err := json.Unmarshal(bytes, &queue); err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return &queue
}

func (b *Bot) DeleteQueue(queue *EventQueue) {
	b.DeleteQueueById(queue.Id)
}

func (b *Bot) DeleteQueueById(queueId string) {
	values := url.Values{}
	values.Set("queue_id", queueId)

	bytes, err := b.getResponseData("DELETE", "events", &values)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

func (b *Bot) FetchQueue(queue *EventQueue) {
	values := url.Values{}
	values.Set("queue_id", queue.Id)
	values.Set("last_event_id", strconv.Itoa(queue.LastEventId))

	bytes, err := b.getResponseData("GET", "events", &values)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}
