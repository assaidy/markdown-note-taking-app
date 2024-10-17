package server

import (
	"github.com/assaidy/markdown-note-takin-app/handlers"
)

func (s *FiberServer) RegisterRoutes() {
	NoteH := handlers.NewNoteHandler(s.DB)

	s.Post("/notes", NoteH.HandleCreateNote)
	s.Get("/notes", NoteH.HandleGetAllNotes)
    s.Get("/notes/:id<int>", NoteH.HandleGetNoteById)
    s.Post("/grammar/", NoteH.HandleGrammarCheck)
}
