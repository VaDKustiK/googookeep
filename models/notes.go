package models

import (
	"database/sql"
	"log"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

func CreateNote(db *sql.DB, title, content string) (int, error) {
	var id int
	query := `INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, title, content).Scan(&id)
	if err != nil {
		log.Println("Error inserting note:", err)
		return 0, err
	}
	return id, nil
}

func GetAllNotes(db *sql.DB) ([]Note, error) {
	rows, err := db.Query("SELECT id, title, content, date FROM notes")
	if err != nil {
		log.Println("Error getting notes:", err)
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Date); err != nil {
			log.Println("Error scanning notes:", err)
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func GetNoteByID(db *sql.DB, id int) (Note, error) {
	var note Note
	query := "SELECT id, title, content, date FROM notes WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&note.ID, &note.Title, &note.Content, &note.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return Note{}, nil
		}
		log.Println("Error getting note:", err)
		return Note{}, err
	}
	return note, nil
}
