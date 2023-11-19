package websocket

import (
	"context"
	"encoding/json"
)

type clientRequest struct {
	UserID   string `json:"id" binding:"required"`
	TripID   string `json:"trip_id" binding:"required"`
	Location string `json:"location"`
}

func HandleRequest(c *Client, plainReq []byte) error {
	var req clientRequest
	err := json.Unmarshal(plainReq, &req)
	if err != nil {
		return err
	}

	ctx := context.Background()

	c.dispatcher.broadcastRoomMap[req.TripID].PublishMessage(ctx, req.Location)

	return nil
}
