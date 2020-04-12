package consumer

import (
	"async-queue/broker"
	"encoding/json"
	"log"
	"sync"
)

const workerCount = 10

type Consumer struct {
	Broker       broker.QueueBackendManager
	TaskRegister map[string]interface{}
}

func NewConsumer(broker broker.QueueBackendManager) *Consumer {
	c := Consumer{Broker: broker, TaskRegister: map[string]interface{}{}}
	return &c
}

func (c *Consumer) RegisterTasks(tasks map[string]interface{}) {
	for key, value := range tasks {
		c.TaskRegister[key] = value
	}
}

func (c *Consumer) GetTask(key string) interface{} {
	task, ok := c.TaskRegister[key]
	if !ok {
		return nil
	}
	return task
}

func (c *Consumer) RunTask(message *string) {
	t := new(broker.Task)
	err := json.Unmarshal([]byte(*message), t)
	if err != nil {
		log.Printf("Error in unmarshalling task %s", err)
	}
	task := c.GetTask(t.TaskName)
	if taskInterface, ok := task.(broker.TaskRunner); ok {
		if err := taskInterface.RunTask(t.Params); err != nil {
			log.Printf("error %s in executing task %+v", err, t.Params)
		} else {
			log.Printf("executed task %+v", t.Params)
		}
	} else {
		log.Printf("task is not found in register")
	}
}

func (c *Consumer) RunConsumer() {
	// initialize pool of workers
	for w := 1; w <= workerCount; w++ {
		log.Printf("spawning worker %d", w)
		go c.worker(w)
	}
}

func (c *Consumer) worker(id int) {
	for {
		log.Println("Polling for messages")
		messages := c.Broker.DequeueTask()
		var wg sync.WaitGroup
		for messageId, messageBody := range messages {
			// spawn one go routine to process each message
			wg.Add(1)
			go c.processTask(wg, messageId, messageBody)
		}
		wg.Wait()
	}
}

func (c *Consumer) processTask(wg sync.WaitGroup, msgId *string, body *string) {
	defer wg.Done()
	log.Printf("Processing Task %s", *msgId)
	c.RunTask(body)
	c.Broker.DeleteTask(msgId)
}
