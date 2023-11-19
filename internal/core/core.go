package core

import (
	"location-server/internal/config"
	"location-server/internal/rabbitmq"
	"location-server/internal/redis"
)

type Server struct {
	router               *Router
	redisClient          redis.Client
	rabbitmqClient       rabbitmq.Client
	rabbitmqChannel      rabbitmq.Channel
	RabbitmqQueueName    string
	RabbitmqExchangeName string
	RabbitmqRoutingKey   string
}

func NewServer() *Server {
	rabbitmqClient := rabbitmq.NewClient()
	rabbitmqChannel := rabbitmqClient.NewChannel()

	return &Server{
		router:               NewRouter(),
		redisClient:          redis.NewClient(),
		rabbitmqClient:       rabbitmqClient,
		rabbitmqChannel:      rabbitmqChannel,
		RabbitmqQueueName:    config.MustGetEnv("RABBITMQ_QUEUE_NAME"),
		RabbitmqExchangeName: config.MustGetEnv("RABBITMQ_EXCHANGE_NAME"),
		RabbitmqRoutingKey:   config.MustGetEnv("RABBITMQ_ROUTING_KEY"),
	}
}

func (s *Server) Close() {
	s.redisClient.Close()
	s.rabbitmqChannel.Close()
	s.rabbitmqClient.Close()
}
