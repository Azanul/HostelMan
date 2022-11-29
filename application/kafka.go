package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		fmt.Println("xxx")
		return nil, err
	}
	return conn, nil
}

func PushCommentToQueue(topic string, message []byte) error {
	kafkaServer := Getenv("KAFKA_BROKER", "172.27.0.6:9092")

	brokersUrl := []string{kafkaServer}

	producer, err := ConnectProducer(brokersUrl)
	if err != nil {
		return err
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		println("abc")
		return err
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
