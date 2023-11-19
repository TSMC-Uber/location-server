package core

import (
	"location-server/internal/websocket"

	"github.com/gin-gonic/gin"
)

type Router struct {
	WsRoomsDispatcher *websocket.RoomsDispatcher
}

func NewRouter() *Router {
	return &Router{
		WsRoomsDispatcher: websocket.NewRoomsDispatcher(),
	}
}

func (s *Server) SetupRouter() *gin.Engine {
	router := gin.Default()

	// public api
	router.GET("/ws/driver", s.DriverWebSocketHandler)
	router.GET("/ws/passenger", s.PassengerWebSocketHandler)

	// private api
	api := router.Group("/api")
	{
		api.POST("/task", s.CreateTask)
	}
	return router
}
