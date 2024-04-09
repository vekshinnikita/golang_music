package tools

import "github.com/gin-gonic/gin"

func GetDomain(c *gin.Context) string {
	domain := ""
	tls := c.Request.TLS
	host := c.Request.Host

	if tls == nil {
		domain += "http://"
	} else {
		domain += "https://"
	}

	domain += host
	return domain
}
