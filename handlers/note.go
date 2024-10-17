package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/assaidy/markdown-note-takin-app/database"
	"github.com/assaidy/markdown-note-takin-app/models"
	"github.com/assaidy/markdown-note-takin-app/utils"
	"github.com/gofiber/fiber/v2"
)

type NoteHandler struct {
	db *database.DBService
}

func NewNoteHandler(db *database.DBService) *NoteHandler {
	return &NoteHandler{db: db}
}

func (h *NoteHandler) HandleCreateNote(c *fiber.Ctx) error {
	// TODO: filesize limit?
	fh, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "provide a file")
	}

	fd, err := fh.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "couldn't read file")
	}
	defer fd.Close()

	cont, err := io.ReadAll(fd)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "couldn't read file")
	}
	if string(cont) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "file is empty")
	}

	var req struct {
		Title string `json:"title"` // TODO: user validator (required)
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json request")
	}

	note := models.Note{
		Title:     req.Title,
		Content:   string(cont),
		CreatedAt: time.Now().UTC(),
	}

	if err := h.db.CreateNote(&note); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("couldn't store note. error: %v", err))
	}

	return c.Status(fiber.StatusCreated).JSON(note)
}

func (h *NoteHandler) HandleGetAllNotes(c *fiber.Ctx) error {
	notes, err := h.db.GetAllNotes()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("couldn't get notes. error: %v", err))
	}

	return c.Status(fiber.StatusOK).JSON(notes)
}

func (h *NoteHandler) HandleGetNoteById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrInternalServerError
	}

	note, err := h.db.GetNoteById(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("couldn't get note. error: %v", err))
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(fiber.StatusOK).Send(utils.MdToHTML([]byte(note.Content)))
}

func (h *NoteHandler) HandleGrammarCheck(c *fiber.Ctx) error {
	var req struct {
		Text string `json:"text"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json request")
	}

	langToolURL := "https://api.languagetool.org/v2/check"
	data := []byte(fmt.Sprintf("text=%s&language=en-US", req.Text))
	checkReq, err := http.NewRequest("POST", langToolURL, bytes.NewBuffer(data))
	if err != nil {
		return fiber.ErrInternalServerError
	}
	checkReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(checkReq)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var grammarResponse struct {
		Matches []struct {
			Message string `json:"message"`
		} `json:"matches"`
	}
	if err := json.Unmarshal(body, &grammarResponse); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(grammarResponse)
}
