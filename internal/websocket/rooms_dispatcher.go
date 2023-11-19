package websocket

import "context"

type RoomsDispatcher struct {
	broadcastRoomMap map[string]*BroadcastRoom
}

func NewRoomsDispatcher() *RoomsDispatcher {
	return &RoomsDispatcher{
		broadcastRoomMap: make(map[string]*BroadcastRoom),
	}
}

func (d *RoomsDispatcher) AcquireBroadcastRoomChannel(id string) *RoomChannel {
	if _, ok := d.broadcastRoomMap[id]; !ok {
		d.broadcastRoomMap[id] = newBroadcastRoom(id)
		ctx := context.Background()
		go d.broadcastRoomMap[id].run(ctx)
	}
	return d.broadcastRoomMap[id].channel
}
