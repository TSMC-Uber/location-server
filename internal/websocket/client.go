package websocket

import (
	"errors"
	"log"
	"time"

	ws "github.com/gorilla/websocket"
)

const (
	// The upper bound on the total number of senders:
	// 1. the response for the websocket request
	// 2. the broadcast from the room
	msgChannelBufferSize = 2

	// Time allowed to write a message to the peer.
	writeMaxWait = 20 * time.Second

	// Send pings to peer with this period.
	pingPeriod  = 20 * time.Second
	pingTimeout = 20 * time.Second

	// Time allowed to read the next message from the peer.
	readMaxWait = pingPeriod + pingTimeout

	roomResponseTimeout = 2 * time.Minute
)

type Client struct {
	*ws.Conn
	dispatcher    *RoomsDispatcher
	pingTicker    *time.Ticker
	messageToSend chan string
	roomChannel   *RoomChannel
}

func (c *Client) SendMessage(msg string) error {
	timer := time.NewTimer(writeMaxWait)
	defer timer.Stop()

	select {
	case c.messageToSend <- msg:
		return nil
	case <-timer.C:
		return errors.New("failed to write into message-to-send channel on time")
	}
}

// sendLoop is a loop to send messages to the websocket connection
func (c *Client) sendLoop() {
	defer c.Close()

	for {
		select {
		case <-c.pingTicker.C:
			c.updateWriteDeadline()
			err := c.WriteMessage(ws.PingMessage, nil)
			if err != nil {
				log.Printf("failed to write ping message to the client %v", err)
				return
			}
		case msg := <-c.messageToSend:
			c.updateWriteDeadline()
			err := c.WriteMessage(ws.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("failed to write message to the client %v", err)
				return
			}
		}
	}
}

func (c *Client) receiveLoop() {
	defer c.Close()

	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			return // Usually caused by a normal client disconnection
		}

		c.updateReadDeadline()

		if msgType == ws.TextMessage { // only handle text message
			err := HandleRequest(c, msg)
			if err != nil {
				log.Printf("failed to handle request %v", err)
			}
			continue
		}
	}
}

func (c *Client) updateReadDeadline() {
	_ = c.SetReadDeadline(time.Now().Add(readMaxWait)) // ignore error
}

func (c *Client) updateWriteDeadline() {
	_ = c.SetWriteDeadline(time.Now().Add(writeMaxWait)) // ignore error
}

func newClient(baseConn *ws.Conn, dispatcher *RoomsDispatcher) *Client {
	return &Client{
		Conn:          baseConn,
		dispatcher:    dispatcher,
		pingTicker:    time.NewTicker(pingPeriod),
		messageToSend: make(chan string, msgChannelBufferSize),
	}
}

func (c *Client) listenRoomClose(closeChan chan struct{}) {
	select {
	case <-closeChan:
		c.Close()
		return
	}
}

func (c *Client) leaveRoom() error {
	channel := c.roomChannel
	if channel == nil {
		return nil
	}

	// If the room channel is already closed, then we don't need to send the
	// `Leaving` channel.
	if !channel.IsClosed() {
		timer := time.NewTimer(roomResponseTimeout)
		defer timer.Stop()

		select {
		case channel.Leaving <- c:
		case <-channel.Closed:
		case <-timer.C:
			return errors.New("leave from broadcast room timeout")
		}
	}

	c.roomChannel = nil
	return nil
}

func (c *Client) doJoin(
	tripID string,
) error {
	if c.roomChannel != nil {
		err := c.leaveRoom()
		if err != nil {
			return err
		}
	}

	timer := time.NewTimer(roomResponseTimeout)
	defer timer.Stop()

	channel := c.dispatcher.AcquireBroadcastRoomChannel(tripID)

	select {
	case channel.Joining <- c:
		c.roomChannel = channel
		go c.listenRoomClose(channel.Closed)
		return nil
	case <-timer.C:
		return errors.New("join to broadcast room timeout")
	}

}

// ServeClient creates a new websocket connection and handles the corresponding interactions.
func ServeClient(
	baseConn *ws.Conn,
	dispatcher *RoomsDispatcher,
	tripID string,
	isDriver bool,
) (*Client, error) {
	client := newClient(baseConn, dispatcher)

	client.SetCloseHandler(func(int, string) error {
		return nil
	})

	client.updateReadDeadline()
	client.SetPongHandler(func(string) error {
		client.updateReadDeadline()
		return nil
	})

	err := client.doJoin(tripID)
	if err != nil {
		return nil, err
	}

	go client.sendLoop()

	if isDriver {
		go client.receiveLoop()
	}

	return client, nil
}
