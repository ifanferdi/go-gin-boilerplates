package middlewares

import (
	"go-gin-boilerplate/internal/constants"

	"github.com/gin-gonic/gin"
)

// NotFoundHandler menangani route yang tidak ditemukan (404) (ekuivalen dengan route-not-found-handler.ts).
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(constants.StatusNotFound, gin.H{
			"message": constants.MessageRouteNotFound,
		})
	}
}
