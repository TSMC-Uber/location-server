package core

import (
	"location-server/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

type webSocketQuery struct {
	UserID string `form:"user_id" binding:"required"`
	TripID string `form:"trip_id" binding:"required"`
}

var upgrader = ws.Upgrader{
	// TODO: check origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) DriverWebSocketHandler(c *gin.Context) {
	var query webSocketQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: check if user is driver for this trip

	// upgrade get request to websocket protocol
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	_, err = websocket.ServeClient(
		conn,
		s.router.WsRoomsDispatcher,
		query.TripID,
		true,
	)
	if err != nil {
		return
	}
}

func (s *Server) PassengerWebSocketHandler(c *gin.Context) {
	var query webSocketQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: check if user is passenger for this trip

	// upgrade get request to websocket protocol
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	_, err = websocket.ServeClient(
		conn,
		s.router.WsRoomsDispatcher,
		query.TripID,
		false,
	)
	if err != nil {
		return
	}
}
