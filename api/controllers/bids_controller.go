package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mrizkyanditama/lelangin/api/models"
	"gopkg.in/olahol/melody.v1"
)

type Resp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UserID  uint32 `json:"user_id"`
	Amount  uint32 `json:"amount"`
}

func (s *Server) initializeBidsController() {
	s.Mel.HandleMessage(s.ProcessBids)
}

func (s *Server) HandleBids(c *gin.Context) {
	s.Mel.HandleRequest(c.Writer, c.Request)
}

func (s *Server) BroadcastBidsResult(m *melody.Session, status string, msg string, amount uint32) {
	resp := &Resp{}
	resp.Status = status
	resp.Message = msg
	respByte, _ := json.Marshal(resp)
	fmt.Println(respByte)
	s.Mel.BroadcastFilter([]byte(respByte), func(q *melody.Session) bool {
		return q.Request.URL.Path == m.Request.URL.Path
	})
}

func (s *Server) ProcessBids(m *melody.Session, msg []byte) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	b := models.Bid{}
	err := json.Unmarshal([]byte(msg), &b)
	if err != nil {
		fmt.Println(err)
		s.BroadcastBidsResult(m, "failed", "Cannot process bid", 0)
		return
	}

	errorMessages := b.Validate(s.DB)
	if len(errorMessages) > 0 {
		s.BroadcastBidsResult(m, "failed", "Failed validating bid", 0)
		return
	}

	_, err = b.AddBid(s.DB)
	if err != nil {
		fmt.Println(err)
		s.BroadcastBidsResult(m, "failed", "Failed saving bid", 0)
		return
	}

	s.BroadcastBidsResult(m, "success", "Success bid", b.Value)

}
