package core

import (
	"github.com/gin-gonic/gin"
)

type createTaskJSON struct {
	Email     string `json:"email"`
	DelayTime int    `json:"delay_time"`
}

func (s *Server) CreateTask(c *gin.Context) {
	var json createTaskJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	s.rabbitmqChannel.SendMsg(s.RabbitmqExchangeName, s.RabbitmqRoutingKey, []byte(json.Email), json.DelayTime)

	c.JSON(200, gin.H{"message": "task created"})
}
