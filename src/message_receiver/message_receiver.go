package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"your_project/file_list"

	"github.com/streadway/amqp"
)

type MessageReceiver interface {
	Listen(tx chan string, wg *sync.WaitGroup)
}

// RabbitMQMessageReceiver listens to messages from a RabbitMQ queue
type RabbitMQMessageReceiver struct {
	channel   *amqp.Channel
	queueName string
}

// NewRabbitMQMessageReceiver initializes the RabbitMQ connection and declares the queue
func NewRabbitMQMessageReceiver() (*RabbitMQMessageReceiver, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		"",    // auto-generated queue name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(
		q.Name,          // queue name
		"",              // routing key
		"eram-messages", // exchange
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQMessageReceiver{channel: ch, queueName: q.Name}, nil
}

// Listen listens for messages from the RabbitMQ queue and sends them over the provided channel
func (r *RabbitMQMessageReceiver) Listen(tx chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	msgs, err := r.channel.Consume(
		r.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %s", err)
	}

	// Use a goroutine to handle incoming messages
	go func() {
		for d := range msgs {
			tx <- string(d.Body)
		}
	}()
}

// FileListMessageReceiver reads files and sends their content over a channel
type FileListMessageReceiver struct {
	fileList file_list.FileList
}

// NewFileListMessageReceiver creates a new instance of FileListMessageReceiver
func NewFileListMessageReceiver(fileList file_list.FileList) *FileListMessageReceiver {
	return &FileListMessageReceiver{fileList: fileList}
}

// Listen reads the content of files from the list and sends it over the channel
func (f *FileListMessageReceiver) Listen(tx chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		for {
			filename := f.fileList.NextFile()
			fmt.Printf("Parsing file %s ...\n", filename)

			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Printf("Failed to read file %s: %s", filename, err)
				continue
			}

			tx <- string(data)
			time.Sleep(time.Second / 60) // Sleep to simulate 60fps-like frequency
		}
	}()
}
