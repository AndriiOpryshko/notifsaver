package main

import (
	"github.com/AndriiOpryshko/notifsaver/notifications"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	InitConfig()
	notsConsumer := notifications.InitConsumer(GetConsumerConfig(), GetS3Conf())

	time.Sleep(10 * time.Second)
	go notsConsumer.Run()

	for {
		time.Sleep(2 * time.Second)
	}
}
