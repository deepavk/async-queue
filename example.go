package main

import (
	"async-queue/backends"
	"async-queue/broker"
	"async-queue/consumer"
	"async-queue/producer"
	"flag"
	"github.com/satori/go.uuid"
	"log"
)

const (
	EmailTask = "send_email"
)

var TaskMap = map[string]interface{}{
	EmailTask: new(Email),
}

type Email struct {
	Id   string `json:"id"`
	Body string `json:"body"`
}

func (e *Email) RunTask(params map[string]interface{}) error {
	if emailId, ok := params["id"].(string); ok {
		e.Id = emailId
	}
	if emailBody, ok := params["body"].(string); ok {
		e.Body = emailBody
	}
	log.Printf("Sending email: %+v", e)
	return nil
}


func main() {
	appType := flag.String("--type", "app","")
	if *appType == "app" {
		startApp()
	} else {
		startConsumer()
	}
}

func startConsumer() {
	backend := getBackend()
	c := consumer.NewConsumer(backend)
	c.RegisterTasks(TaskMap)
	c.RunConsumer()
}

func startApp() {
	backend := getBackend()
	p := producer.NewProducer(backend)
	// Enqueue email tasks
	for i := 0; i < 20; i++ {
		taskId := uuid.NewV1()
		params := map[string]interface{}{"id": "xyz@gmail.com", "body": "Example email"}
		p.EnqueueTask(EmailTask, params, taskId.String())
	}
}

func getBackend() broker.QueueBackendManager {
	return backends.NewSQSService()
}
