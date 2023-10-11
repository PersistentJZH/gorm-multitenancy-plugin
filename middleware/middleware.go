package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ScopeIDSetter(dbs ...*gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopeID := c.Request.Header.Get("Scope-ID")
		if strings.EqualFold(scopeID, "") {
			log.Println("Scope-ID empty!")
			c.Abort()
		}
		ctx := context.WithValue(c.Request.Context(), "scope_id", scopeID)
		for _, db := range dbs {
			db = db.WithContext(ctx)
		}
		c.Next()
	}
}
