package main

import (
	"location-server/internal/core"
	"log"
)

func main() {
	server := core.NewServer()
	defer server.Close()

	router := server.SetupRouter()

	// start worker to listen to rabbitmq queue and send email
	go server.SendEmailWorker(
		server.RabbitmqQueueName,
		server.RabbitmqExchangeName,
		server.RabbitmqRoutingKey,
	)

	// start server
	if err := router.Run(); err != nil {
		log.Panicf("failed to run server: %v", err)
	}
}
