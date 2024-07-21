package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Snippet struct {
	ID      int
	Content string
}

func main() {
	// Connect to PostgreSQL database
	var err error
	connStr := "postgres://username:password@hostname:port/dbname?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Set up Gin router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Serve static files
	router.Static("/static", "./static")

	// Routes
	router.GET("/", showEditor)
	router.POST("/save", saveSnippet)
	router.GET("/snippets", listSnippets)

	// Run the server
	router.Run(":8080")
}

func showEditor(c *gin.Context) {
	c.HTML(http.StatusOK, "editor.html", nil)
}

func saveSnippet(c *gin.Context) {
	content := c.PostForm("content")
	_, err := db.Exec("INSERT INTO snippets (content) VALUES ($1)", content)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving snippet")
		return
	}
	c.String(http.StatusOK, "Snippet saved")
}

func listSnippets(c *gin.Context) {
	rows, err := db.Query("SELECT id, content FROM snippets")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching snippets")
		return
	}
	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
		var snippet Snippet
		if err := rows.Scan(&snippet.ID, &snippet.Content); err != nil {
			c.String(http.StatusInternalServerError, "Error scanning snippet")
			return
		}
		snippets = append(snippets, snippet)
	}

	c.HTML(http.StatusOK, "snippets.html", gin.H{
		"snippets": snippets,
	})
}
