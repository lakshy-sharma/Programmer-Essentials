package mysql

import (
	"database/sql"
	"preserve/pkg/models"
)

type NoteModel struct {
	DB *sql.DB
}

// Insert a new note in the database.
func (m *NoteModel) Insert(title, content, expires string) (int, error) {
	// The SQL statement which will be executed for creating a new note.
	sqlStatement := `INSERT INTO notes (title, content, created, expires)
	VALUES(?,?,UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(sqlStatement, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Retrieve a note from the database.
func (m *NoteModel) Get(id int) (*models.Note, error) {
	// The SQL statement to be executed.
	sqlStatement := `SELECT id, title, content, created, expires FROM notes WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// A simple pointer to the model for reading in the values by scanning the rows.
	capturedNote := models.Note{}

	// Query and scan the rows of the database for a match.
	err := m.DB.QueryRow(sqlStatement, id).Scan(&capturedNote.ID, &capturedNote.Title, &capturedNote.Content, &capturedNote.Created, &capturedNote.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return &capturedNote, nil
}

// Get 10 latest notes.
func (m *NoteModel) Latest() ([]*models.Note, error) {
	// The SQL statement to search data.
	sqlStatement := `SELECT id,title,content,created,expires FROm notes WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	notes := []*models.Note{}

	for rows.Next() {
		note := &models.Note{}
		err = rows.Scan(&note.ID, &note.Content, &note.Created, &note.Expires)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	// Check if all the iterations were completely successful or not.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
