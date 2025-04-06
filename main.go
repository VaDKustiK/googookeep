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

	// Create a new note
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

	// Get all notes
	router.GET("/notes", func(c *gin.Context) {
		notes, err := models.GetAllNotes(db.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting notes"})
			return
		}
		c.JSON(http.StatusOK, notes)
	})

	// Get a note by ID
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

	// Delete a note by ID
	router.DELETE("/notes/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
			return
		}

		err = models.DeleteNoteByID(db.DB, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
	})

	// Update a note by ID
	router.PUT("/notes/:id", func(c *gin.Context) {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Note updated"})
	})

	fmt.Println("Server running on :8080")
	router.Run(":8080")
}
