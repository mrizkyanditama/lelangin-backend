package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrizkyanditama/lelangin/api/models"
)

func (server *Server) GetTags(c *gin.Context) {

	tag := models.Tag{}

	tags, err := tag.FindAllTags(server.DB)
	if err != nil {
		errList["No_tag"] = "No tag found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": tags,
	})
}
