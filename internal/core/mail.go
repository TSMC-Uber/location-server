package core

import (
	"log"

	"github.com/streadway/amqp"
)

func (s *Server) SendEmailWorker(qName, exName, routingKey string) {
	msgs := s.rabbitmqChannel.NewMsgReceiver(qName, exName, routingKey)
	for msg := range msgs {
		s.sendEmail(msg)
	}
}

func (s *Server) sendEmail(msg amqp.Delivery) {
	log.Printf("send email to %s", string(msg.Body))
	// TODO: send email
}
