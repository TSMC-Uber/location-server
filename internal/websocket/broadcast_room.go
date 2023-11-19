package websocket

import (
	"context"
	"location-server/internal/redis"
	"log"
)

const (
	// Intend to block the join/leave operations.
	roomParticChannelBufferSize = 0
)

type RoomChannel struct {
	Joining chan *Client
	Leaving chan *Client
	// This channel would be closed when the room is closed
	Closed chan struct{}
}

func (c *RoomChannel) IsClosed() bool {
	select {
	case <-c.Closed:
		return true
	default:
		return false
	}
}

func (c *RoomChannel) Close() {
	close(c.Closed)
}

// BroadcastRoom maintains the set of active clients and broadcasts messages from the redis pub-sub-channel to them.
type BroadcastRoom struct {
	id      string
	clients map[*Client]bool
	channel *RoomChannel
}

func newBroadcastRoom(id string) *BroadcastRoom {
	return &BroadcastRoom{
		id:      id,
		clients: make(map[*Client]bool),
		channel: &RoomChannel{
			Joining: make(chan *Client, roomParticChannelBufferSize),
			Leaving: make(chan *Client, roomParticChannelBufferSize),
			Closed:  make(chan struct{}),
		},
	}
}

func (r *BroadcastRoom) GetClientCount() int {
	return len(r.clients)
}

// Handles the corresponding interactions from channels.
// This room will automatically close when the last client leaves.
func (r *BroadcastRoom) run(ctx context.Context) {
	defer func() {
		r.close()
	}()

	pubSub := redis.NewClient().Subscribe(ctx, r.id)
	defer func() {
		err := pubSub.Unsubscribe(ctx, r.id)
		if err != nil {
			log.Println("failed to unsubscribe from redis pub-sub-channel", err)
		}
	}()

	channel := r.channel
	redisChannel := pubSub.Channel()
	for {
		if channel.IsClosed() {
			return
		}

		select {
		case client := <-channel.Joining:
			r.clients[client] = true
		case client := <-channel.Leaving:
			r.remove(client)
		case <-ctx.Done():
			return
		case message, ok := <-redisChannel:
			if !ok {
				return
			}

			for client := range r.clients {
				if err := client.SendMessage(message.Payload); err != nil {
					log.Println("failed to send message to client", err)
					r.remove(client)
				}
			}
		}
	}
}

func (r *BroadcastRoom) PublishMessage(ctx context.Context, message string) error {
	return redis.NewClient().Publish(ctx, r.id, message).Err()
}

// remove remove the given client when it is joined. And close this room when the last client been removed.
func (r *BroadcastRoom) remove(client *Client) {
	if _, ok := r.clients[client]; !ok {
		return
	}

	delete(r.clients, client)

	if len(r.clients) == 0 {
		r.close()
	}
}

// close do corresponding cleanups (when it is necessary)
func (r *BroadcastRoom) close() {
	if r.channel.IsClosed() {
		return
	}

	r.channel.Close()
	clear(r.clients)
}
