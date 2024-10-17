package database

import (
	"database/sql"

	"github.com/assaidy/markdown-note-takin-app/models"
	_ "github.com/mattn/go-sqlite3"
)

type DBService struct {
	*sql.DB
}

var (
	dbName   = "database/notes.db"
	instance *DBService
)

func NewDBService() *DBService {
	if instance != nil {
		return instance
	}
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	instance = &DBService{DB: db}
	return instance
}

func (dbs *DBService) CreateNote(inout *models.Note) error {
	query := `
    INSERT INTO notes (title, content, created_at)
    VALUES (?, ?, ?);
    `
	res, err := dbs.Exec(query, inout.Title, inout.Content, inout.CreatedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	inout.Id = int(id)

	return nil
}

func (dbs *DBService) GetAllNotes() ([]*models.Note, error) {
	query := `
    SELECT id, title, created_at
    FROM notes;
    `
	rows, err := dbs.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]*models.Note, 0)
	for rows.Next() {
		note := models.Note{}
		if err := rows.Scan(&note.Id, &note.Title, &note.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, err
}

func (dbs *DBService) GetNoteById(id int) (*models.Note, error) {
	query := `
    SELECT title, content, created_at
    FROM notes
    WHERE id = ?;
    `
	note := models.Note{Id: id}
	if err := dbs.QueryRow(query, id).Scan(&note.Title, &note.Content, &note.CreatedAt); err != nil {
		return nil, err
	}

    return &note, nil
}
