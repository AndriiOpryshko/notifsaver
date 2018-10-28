package notifications

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os/signal"
	"strings"
	"time"
	"fmt"
	"os"
)

type ConsumerConfig interface {
	GetAddr() string
	GetTopic() string
	GetClientId() string
}

func InitConsumer(confConsumer ConsumerConfig, s3Conf S3Conf) *NotificationConsumer {
	return &NotificationConsumer{
		conf: confConsumer,
		S3Sevice: InitS3(s3Conf),
		Nots: make(chan *Notifications),
	}
}

type NotificationConsumer struct {
	conf ConsumerConfig
	S3Sevice *S3Sevice
	Nots chan *Notifications
}

func (nc NotificationConsumer) Run() {
	config := sarama.NewConfig()
	config.ClientID = nc.conf.GetClientId() //"go-kafka-consumer"
	config.Consumer.Return.Errors = true
	config.Consumer.Retry.Backoff = 2 * time.Second

	brokers := []string{nc.conf.GetAddr()} //"kafka:9092"}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Error while connection to kafka")

		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)
		if r := recover(); r != nil {
			log.WithFields(log.Fields{
				"r": r,
			}).Info("Recovered")
		}
	}()

	consumer, errors := consume([]string{nc.conf.GetTopic()}, master)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// Get signnal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				msgCount++
				var nots Notifications
				err := proto.Unmarshal(msg.Value, &nots)

				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Error("Error parsing received proto3 messages")
				}

				path := fmt.Sprintf("%d_%d", time.Now().Unix(), msgCount)

				nc.S3Sevice.AddObjectToS3(path, msg.Value)

				log.WithFields(log.Fields{
					"key":  msg.Key,
					"nots": nots,
				}).Info("Received messages")
			case consumerError := <-errors:
				msgCount++
				log.WithFields(log.Fields{
					"err_topic":     string(consumerError.Topic),
					"err_partition": string(consumerError.Partition),
					"err_consumer":  consumerError.Err,
				}).Error("Received messages error")
				doneCh <- struct{}{}
			case <-signals:
				log.WithFields(log.Fields{
					"signals": signals,
				}).Warning("Interrupt notification consumer")

				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.WithFields(log.Fields{
		"count": msgCount,
	}).Debug("Processed messages")

}

func consume(topics []string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {
	consumers := make(chan *sarama.ConsumerMessage)
	errors := make(chan *sarama.ConsumerError)
	for _, topic := range topics {
		if strings.Contains(topic, "__consumer_offsets") {
			continue
		}
		partitions, _ := master.Partitions(topic)
		// this only consumes partition no 1, you would probably want to consume all partitions
		consumer, err := master.ConsumePartition(topic, partitions[0], sarama.OffsetOldest)
		if nil != err {
			log.WithFields(log.Fields{
				"topic":     topic,
				"partition": partitions,
				"err":       err,
			}).Error("Error while receiving notifications")
			panic(err)
		}
		log.WithFields(log.Fields{
			"topic": topic,
		}).Info("Start consuming topic")
		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					errors <- consumerError
					log.WithFields(log.Fields{
						"err": consumerError.Err,
					}).Error("Consuming error")

				case msg := <-consumer.Messages():
					consumers <- msg

					var nots Notifications
					err := proto.Unmarshal(msg.Value, &nots)

					if err != nil {
						log.WithFields(log.Fields{
							"err": err,
						}).Error("Error parsing received protoe messages")
					}

					log.WithFields(log.Fields{
						"topic": topic,
						"value": nots,
					}).Info("Got notification")
				}
			}
		}(topic, consumer)
	}

	return consumers, errors
}
