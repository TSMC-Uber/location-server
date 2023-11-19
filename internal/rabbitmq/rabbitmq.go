package rabbitmq

import (
	"location-server/internal/config"
	"log"

	"github.com/streadway/amqp"
)

type Client interface {
	Close() error

	NewChannel() Channel
}

type ClientImpl struct {
	*amqp.Connection
}

func NewClient() Client {
	conn, err := amqp.Dial(config.MustGetEnv("RABBITMQ_HOST"))
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}

	return &ClientImpl{
		Connection: conn,
	}
}

func (c *ClientImpl) NewChannel() Channel {
	ch, err := c.Connection.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	return &ChannelImpl{
		Channel: ch,
	}
}

type Channel interface {
	declareQueue(qName string) error
	declareDelayExchange(exName string) error
	bindQueue(string, string, string) error
	consume(qName string) (<-chan amqp.Delivery, error)
	SendMsg(string, string, []byte, int) error
	NewMsgReceiver(string, string, string) <-chan amqp.Delivery

	Close() error
}

type ChannelImpl struct {
	*amqp.Channel
}

func (c *ChannelImpl) declareQueue(qName string) error {
	_, err := c.QueueDeclare(
		qName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (c *ChannelImpl) declareDelayExchange(exName string) error {
	args := amqp.Table{
		"x-delayed-type": "direct",
	}
	return c.ExchangeDeclare(
		exName,              // name
		"x-delayed-message", // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		args,                // arguments
	)
}

func (c *ChannelImpl) bindQueue(qName, exName, routingKey string) error {
	return c.QueueBind(
		qName,      // queue name
		routingKey, // routing key
		exName,     // exchange
		false,      // no-wait
		nil,
	)
}

func (c *ChannelImpl) consume(qName string) (<-chan amqp.Delivery, error) {
	return c.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func (c *ChannelImpl) SendMsg(exName, routingKey string, msg []byte, delayPeriod int) error {
	return c.Publish(
		exName,     // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
			Headers: amqp.Table{
				"x-delay": delayPeriod,
			},
		})
}

func (c *ChannelImpl) NewMsgReceiver(qName, exName, routingKey string) <-chan amqp.Delivery {
	if err := c.declareQueue(qName); err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	if err := c.declareDelayExchange(exName); err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	if err := c.bindQueue(qName, exName, routingKey); err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	msgs, err := c.consume(qName)
	if err != nil {
		log.Fatalf("failed to consume: %v", err)
	}

	return msgs
}
