package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// APINoCache sets Cache-Control headers on API responses to prevent proxy/CDN caching.
// This prevents user-specific responses (e.g. /auth/me) from being cached by reverse
// proxies and served to different users.
func APINoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}
		c.Next()
	}
}
