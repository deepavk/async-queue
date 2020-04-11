package backends

import (
	"async-queue/broker"
	"fmt"
	"strconv"
)

var messages []string

type PlaceholderBackend struct {}

func NewPlaceholderBackend() broker.QueueBackendManager {
	return new(PlaceholderBackend)
}

func (pb *PlaceholderBackend) EnqueueTask(message string) error {
	messages = append(messages, message)
	fmt.Printf("Enqueued task: %s", message)
	return nil
}

func (pb *PlaceholderBackend) DequeueTask() map[*string]*string {
	fmt.Println("messages", messages)
	m := make(map[*string]*string)
	ctr := 0
	for _, msg := range messages {
		key := "task_id" + strconv.Itoa(ctr)
		ctr += 1
		m[&key] = &msg
	}
	return m
}

func (pb *PlaceholderBackend) DeleteTask(msgId *string) {
	fmt.Printf("Task deleted %s", *msgId)
}