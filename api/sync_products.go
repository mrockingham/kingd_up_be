package api

import (
	"kingdup/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SyncProductsHandler(c *gin.Context) {
	if err := services.SyncProductsFromPrintful(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Products synced!"})
}
