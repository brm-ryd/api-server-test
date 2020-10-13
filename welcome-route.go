package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Welcome page for GET handler
func WelcomeHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	f, err := resourcesBox.Open("welcome-page.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	_, err = io.Copy(c.Writer, f)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}
