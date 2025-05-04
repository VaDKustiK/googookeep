package main

import (
	"database/sql"
	"fmt"
	"googookeep/db"
	"googookeep/models"
	"googookeep/utils"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/note.html",
	))
}

func main() {
	db.ConnectDB()
	defer db.DB.Close()

	router := gin.Default()
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		notes, err := models.GetAllNotes(db.DB)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error loading notes")
			return
		}

		c.Header("Content-Type", "text/html")
		err = templates.ExecuteTemplate(c.Writer, "index.html", gin.H{
			"Title": "Your notes",
			"Notes": notes,
		})
		if err != nil {
			fmt.Println("TEMPLATE ERROR:", err)
			c.String(http.StatusInternalServerError, "Template rendering error")
		}
	})

	router.GET("/note/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid note ID")
			return
		}

		note, err := models.GetNoteByID(db.DB, id)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error loading note")
			return
		}

		c.Header("Content-Type", "text/html")
		err = templates.ExecuteTemplate(c.Writer, "note.html", gin.H{
			"Title":  note.Title,
			"Note":   note,
			"IsEdit": true,
		})
		if err != nil {
			fmt.Println("TEMPLATE ERROR:", err)
			c.String(http.StatusInternalServerError, "Template rendering error")
		}
	})

	router.POST("/notes/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
			return
		}

		title := c.PostForm("title")
		content := c.PostForm("content")

		if title == "" || content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title and content cannot be empty"})
			return
		}

		err = models.UpdateNoteByID(db.DB, id, title, content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating note"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/api/notes", func(c *gin.Context) {
		var newNote models.Note
		if err := c.ShouldBindJSON(&newNote); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		noteID, err := models.CreateNote(db.DB, newNote.Title, newNote.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating note"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": noteID, "message": "Note created"})
	})

	router.GET("/api/notes", func(c *gin.Context) {
		notes, err := models.GetAllNotes(db.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting notes"})
			return
		}
		c.JSON(http.StatusOK, notes)
	})

	router.PUT("/api/notes/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
			return
		}

		var updatedNote models.Note
		if err := c.ShouldBindJSON(&updatedNote); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		err = models.UpdateNoteByID(db.DB, id, updatedNote.Title, updatedNote.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating note"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Note updated"})
	})

	router.DELETE("/api/notes/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
			return
		}

		err = models.DeleteNoteByID(db.DB, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting note"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
	})

	router.POST("/api/notes/order", func(c *gin.Context) {
		var newOrder []struct {
			ID    int `json:"id"`
			Order int `json:"order"`
		}

		if err := c.ShouldBindJSON(&newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		for _, item := range newOrder {
			_, err := db.DB.Exec("UPDATE notes SET position = $1 WHERE id = $2", item.Order, item.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order updated"})
	})

	router.POST("/api/notes/import", func(c *gin.Context) {
		var payload struct {
			Code string `json:"code"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		shared, err := models.GetNoteByShareCode(db.DB, payload.Code)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup failed"})
			return
		}

		newID, err := models.CreateNote(db.DB, shared.Title, shared.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not import"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": newID})
	})

	router.POST("/api/notes/share/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		note, err := models.GetNoteByID(db.DB, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup failed"})
			return
		}

		if note.ShareCode == "" {
			code := utils.GenerateShareCode(8)
			if _, err := db.DB.Exec("UPDATE notes SET share_code=$1 WHERE id=$2", code, id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save code"})
				return
			}
			note.ShareCode = code
		}

		c.JSON(http.StatusOK, gin.H{"code": note.ShareCode})
	})

	fmt.Println("Server running on :8080")
	router.Run(":8080")
}
