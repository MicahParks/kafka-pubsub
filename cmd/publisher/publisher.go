package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"

	kfpubsub "github.com/MicahParks/kafka-pubsub"
)

func main() {

	// Gather the required information for environment variables.
	hostPort := os.Getenv("KAFKA_HOST_PORT")
	topic := os.Getenv("KAFKA_TOPIC")

	// Create a logger.
	logger := log.New(os.Stdout, "", 0)

	// Create a context that will time out the connection to the leader if it takes too long.
	ctx, cancel := context.WithTimeout(context.Background(), kfpubsub.DefaultLeaderTimeout)

	// Create a Kafka connection to publish messages with. Close it when the program is done. Try this for 5 seconds.
	var conn *kafka.Conn
	var err error
	for i := 0; i < 5; i++ {
		conn, err = kafka.DialLeader(ctx, "tcp", hostPort, topic, 0)
		if err != nil {
			logger.Printf("Failed to connect to Kafka leader.\nError: %s", err.Error())
			time.Sleep(time.Second)
			continue
		}
		break
	}
	cancel()

	// If no connection has been formed, fail the program.
	if conn == nil {
		logger.Fatalf("Failed to dial Kafka leader.\nError: %s", err.Error())
	}

	// The connection has been made. Make sure it properly closes.
	logger.Println("Connected to Kafka leader.")
	defer conn.Close() // Ignore any error.

	// Set the random seem from the current time.
	rand.Seed(time.Now().UnixNano())

	// Make a channel to catch the signal from the OS. It'll be notified of Ctrl + C.
	ctrlC := make(chan os.Signal, 1)

	// Tell the program to monitor for an interrupt or SIGTERM and report it on the given channel.
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)

	// Create a ticker that will tick every second. This will be how often messages are sent.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Publish messages to Kafka until told to stop.
	for {
		select {

		// Wait for a condition that means the program should end.
		case <-ctrlC:
			return

		// Wait for the ticker to tick, then send a message.
		case <-ticker.C:

			// Create the Kafka message.
			var message *kafka.Message
			message, err = kfpubsub.CreateKafkaMessage(topic)
			if err != nil {
				logger.Fatalf("Failed to create Kafka message.\nError: %s", err.Error())
			}

			// Publish the message to Kafka.
			if _, err = conn.WriteMessages(*message); err != nil {
				logger.Fatalf("Failed to publish Kafka message.\nError: %s", err.Error())
			}
		}
	}
}
