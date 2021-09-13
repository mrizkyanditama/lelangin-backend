package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrizkyanditama/lelangin/api/models"
)

func (server *Server) GetCategories(c *gin.Context) {

	category := models.Category{}

	categories, err := category.FindAllCategories(server.DB)
	if err != nil {
		errList["No_category"] = "No category found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": categories,
	})
}
