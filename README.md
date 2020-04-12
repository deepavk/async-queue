Implementation of an asynchronous task queueing capability


**Broker**
The broker allows connection to a backend (sqs,redis), enqueue/dequeue operations.
It defines the interfaces required by the producer and consumer to interact with queue backends

**Tasks:**
Tasks have to be defined in a task map and registered at the consumer.

**Backends:**
SQS is implemented here as the backend. Other backends such as redis can be added by extending the QueueBackendManager 
interface

**Consumer:**
The consumer spawns a pool of workers that dequeue messages in parallel. 



**Example usage:**
An application is added in example.go: As seen here, the producer and consumer can be imported in the application 
* producer: is used to send tasks to a queue backend to be processed asynchronously.
* task: The email task is registered with the consumer while running the consumer. 
* consume tasks: Start a consumer to poll the backend infinitely and execute the task. 

To start the sample application: 
```$xslt
    go run example.go --type=app
```

To start the consumer:
```$xslt
    go run example.go --type=worker 
```

