package main

import (
	"context"
	"log"
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
			time.Sleep(time.Second * 1)
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

	// Make a channel to catch the signal from the OS. It'll be notified of Ctrl + C.
	ctrlC := make(chan os.Signal, 1)

	// Tell the program to monitor for an interrupt or SIGTERM and report it on the given channel.
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)

	// Create a channel that will close when the program wants all goroutines to end.
	death := make(chan struct{})

	// Create a channel that will have all received messages sent over it.
	msgChan := readMessageChan(kfpubsub.MessageBuffer, conn, death, logger, kfpubsub.MessageMaxBytes)

	// Print messages from Kafka until told to stop.
	for {
		select {

		// Wait for a condition that means the program should end.
		case <-ctrlC:
			return

		// Wait for a Kafka message to be received.
		case message, ok := <-msgChan:

			// Check if the channel was closed.
			if !ok {
				return
			}

			// Unmarshal the Kafka message into a Go structure.
			var reading *kfpubsub.Reading
			if reading, err = kfpubsub.UnmarshalMessage(message); err != nil {
				logger.Fatalln("Failed to unmarshal Kafka message into a Go struct.\nError: %s", err.Error())
			}

			// Print the received message.
			logger.Printf("New message:\n  Epoch in seconds: %d\n  Celsius reading: %.2f°C", reading.Epoch, reading.Celsius)
		}
	}
}

// readMessageChan creates a channel that will read messages from the Kafka channel.
func readMessageChan(buffer uint, conn *kafka.Conn, death chan struct{}, logger *log.Logger, maxBytes int) (msgChan chan *kafka.Message) {

	// Create the channel to send messages over.
	msgChan = make(chan *kafka.Message, buffer)

	// Launch a separate goroutine that will read messages and send them over the channel.
	go func() {

		// Read messages until there is a failure to do so.
		for {

			// Read the message from Kafka.
			message, err := conn.ReadMessage(maxBytes)
			if err != nil {
				logger.Printf("Failed to read message from Kafka.\nError: %s", err.Error())
				close(msgChan)
				return
			}

			// Send the message over the channel, unless dead.
			select {
			case msgChan <- &message:
			case <-death:
				return
			}
		}
	}()

	return msgChan
}
