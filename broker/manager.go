//Package broker includes interfaces required by the producer and consumer to interact with queue backends
//The backend implements specific logic for connecting to a backend (sqs,redis), enqueue/dequeue operations
//The producer enqueues task, consumer dequeues it and deletes task

package broker

type QueueBackendManager interface {
	EnqueueTask(message string) error
	DequeueTask() map[*string]*string
	DeleteTask(messageId *string)
}

// Interface implemented by the consumer to execute the dequeued task
type TaskRunner interface {
	RunTask(params map[string]interface{}) error
}

type Task struct {
	TaskId   string                 `json:"task_id"`
	TaskName string                 `json:"task_name"`
	Params   map[string]interface{} `json:"params"`
}
