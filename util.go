package kafka

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

// TODO These constants could become configurations.
const (

	// DefaultLeaderTimeout is the amount of time to wait before a connection to the Kafka leader has considered to time
	// out.
	DefaultLeaderTimeout = time.Second * 10

	// MessageMaxBytes is the maximum quantity of bytes a Kafka message can be.
	MessageMaxBytes = 10e8

	// MessageBuffer is the maximum number of Kafka messages to buffer from a reading channel.
	MessageBuffer = 5
)

// CreateKafkaMessage creates the Kafka message to publish to Kafka.
func CreateKafkaMessage(topic string) (message *kafka.Message, err error) {

	// Create a fake reading to put in the message to publish.
	reading := createReading()

	// Turn the reading into a JSON byte slice.
	var data []byte
	if data, err = json.Marshal(reading); err != nil {
		return nil, err
	}

	// Create the Kafka message.
	return &kafka.Message{
		Topic: topic,
		Value: data,
	}, nil
}

// createReading creates a random fake temperature reading for the current time.
func createReading() *Reading {

	// Create a random temperature reading based on the temperatures from this website:
	// https://earthobservatory.nasa.gov/global-maps/MOD_LSTD_M
	celsius := rand.Float64() + float64(rand.Intn(45+25)-25)

	// Create the Reading Go structure.
	return &Reading{
		Celsius: celsius,
		Epoch:   time.Now().Unix(),
	}
}

// UnmarshalMessage unmarshalls the Kafka message into the Reading Go structure.
func UnmarshalMessage(message *kafka.Message) (reading *Reading, err error) {
	reading = &Reading{}
	if err = json.Unmarshal(message.Value, reading); err != nil {
		return nil, err
	}
	return reading, nil
}
