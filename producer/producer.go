package producer

import (
	"async-queue/broker"
	"encoding/json"
	"log"
)

type Producer struct {
	broker broker.QueueBackendManager
}

func NewProducer(b broker.QueueBackendManager) Producer {
	return Producer{broker: b}
}

func (p Producer) EnqueueTask(taskName string, params map[string]interface{}, corrId string) {
	tm := broker.Task{TaskName: taskName, Params: params, TaskId: corrId}
	message, err := json.Marshal(tm)
	if err != nil {
		log.Printf("Error in unmarshalling task %s", err)
	}
	err = p.broker.EnqueueTask(string(message))
	if err != nil {
		log.Printf("Error in message enqueue %s", err)
	}
}
