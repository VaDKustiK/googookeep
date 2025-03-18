package main

import (
	"fmt"
	"googookeep/db"
	"googookeep/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()
	defer db.DB.Close()

	router := gin.Default()

	router.POST("/notes", func(c *gin.Context) {
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

	router.GET("/notes", func(c *gin.Context) {
		notes, err := models.GetAllNotes(db.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting notes"})
			return
		}
		c.JSON(http.StatusOK, notes)
	})

	router.GET("/notes/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
			return
		}

		note, err := models.GetNoteByID(db.DB, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting note"})
			return
		}
		if note.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		c.JSON(http.StatusOK, note)
	})

	fmt.Println("Server running on :8080")
	router.Run(":8080")
}
