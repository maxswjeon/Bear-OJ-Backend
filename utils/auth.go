package utils

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthUser(c *gin.Context) bool {
	session := sessions.Default(c)

	if session.Get("user") == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"result": false,
			"error":  "Unauthorized"})
		return false
	}

	return true
}

func AuthAdmin(c *gin.Context) bool {
	session := sessions.Default(c)

	if session.Get("admin") == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"result": false,
			"error":  "Unauthorized"})
		return false
	}

	return true
}

func AuthAll(c *gin.Context) bool {
	session := sessions.Default(c)

	if session.Get("user") == nil && session.Get("admin") == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"result": false,
			"error":  "Unauthorized"})
		return false
	}

	return true
}
