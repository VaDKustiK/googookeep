package models

import (
	"database/sql"
	"googookeep/utils"
	"log"
)

type Note struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Date      string `json:"date"`
	Position  int    `json:"position"`
	ShareCode string `json:"share_code"`
}

func CreateNote(db *sql.DB, title, content string) (int, error) {
	var id int
	var maxPosition int
	err := db.QueryRow(`SELECT COALESCE(MAX(position), 0) FROM notes`).Scan(&maxPosition)
	if err != nil {
		log.Println("Error getting max position:", err)
		return 0, err
	}

	// Generate the share code once here
	shareCode := utils.GenerateShareCode(8)

	query := `
		INSERT INTO notes (title, content, position, share_code)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = db.QueryRow(query, title, content, maxPosition+1, shareCode).Scan(&id)
	if err != nil {
		log.Println("Error inserting note:", err)
		return 0, err
	}

	return id, nil
}

func GetAllNotes(db *sql.DB) ([]Note, error) {
	rows, err := db.Query(`
		SELECT id, title, content, position, share_code
		FROM notes
		ORDER BY position ASC
	`)
	if err != nil {
		log.Println("Error querying notes:", err)
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Position, &n.ShareCode); err != nil {
			log.Println("Error scanning notes:", err)
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func GetNoteByID(db *sql.DB, id int) (Note, error) {
	var note Note
	query := `
		SELECT id, title, content, date, position, share_code
		FROM notes
		WHERE id = $1
	`
	err := db.QueryRow(query, id).
		Scan(&note.ID, &note.Title, &note.Content, &note.Date, &note.Position, &note.ShareCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return Note{}, nil
		}
		log.Println("Error getting note:", err)
		return Note{}, err
	}
	return note, nil
}

func DeleteNoteByID(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		log.Println("Error deleting note:", err)
	}
	return err
}

func UpdateNoteByID(db *sql.DB, id int, title, content string) error {
	_, err := db.Exec(
		`UPDATE notes SET title = $1, content = $2 WHERE id = $3`,
		title, content, id,
	)
	if err != nil {
		log.Println("Error updating note:", err)
	}
	return err
}

func UpdateNotePosition(db *sql.DB, id, position int) error {
	_, err := db.Exec(
		`UPDATE notes SET position = $1 WHERE id = $2`,
		position, id,
	)
	if err != nil {
		log.Println("Error updating note position:", err)
	}
	return err
}

func GetNoteByShareCode(db *sql.DB, code string) (Note, error) {
	var n Note
	query := `
		SELECT id, title, content, date, position, share_code
		FROM notes
		WHERE share_code = $1
	`
	err := db.QueryRow(query, code).
		Scan(&n.ID, &n.Title, &n.Content, &n.Date, &n.Position, &n.ShareCode)
	if err != nil {
		return Note{}, err
	}
	return n, nil
}
