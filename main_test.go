package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/save", func(c *gin.Context) {
		code := c.PostForm("code")
		err := saveCode(code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save code"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Code saved successfully"})
	})

	return r
}

func TestIndexRoute(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<title>Code Editor</title>")
}

func TestSaveRoute(t *testing.T) {
	router := setupRouter()

	formData := "code=package+main%0Aimport+%28%0A%09%22fmt%22%0A%29%0Afunc+main%28%29+%7B%0A%09fmt.Println%28%22Hello%2C+world!%22%29%0A%7D"
	req, _ := http.NewRequest("POST", "/save", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Code saved successfully")
}
