package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func eprint(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func main() {
	var (
		url      string
		exchange string
		bind     string
	)

	flag.StringVar(&url, "url", "amqp://guest:guest@localhost:5672/", "AMQP server URL")
	flag.StringVar(&exchange, "exchange", "", "Exchange to bind.")
	flag.StringVar(&bind, "bind", "", "Binding key.")
	flag.Parse()

	conn, err := amqp.Dial(url)
	if err != nil {
		eprint("failed to connect: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		eprint("failed to create channel: %s\n", err)
		os.Exit(1)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		eprint("failed to declare queue: %s\n", err)
		os.Exit(1)
	}

	err = ch.QueueBind(
		queue.Name,
		bind,
		exchange,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		eprint("failed to bind the queue: %s\n", err)
		os.Exit(1)
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",    // consumer
		true,  // auto-ack
		true,  // exclusive
		false, // no-local
		false, // no-wait
		nil,   //arguments
	)
	if err != nil {
		eprint("failed to consume from queue: %s\n", err)
		os.Exit(1)
	}

	for d := range msgs {
		fmt.Println(string(d.Body))
	}
}
