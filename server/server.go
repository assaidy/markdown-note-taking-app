package server

import (
	"github.com/assaidy/markdown-note-takin-app/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type FiberServer struct {
	*fiber.App
	DB *database.DBService
}

func NewFiberServer() *FiberServer {
	fs := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "markdown-note-takin-app",
			AppName:      "markdown-note-takin-app",
		}),
		DB: database.NewDBService(),
	}
	fs.Use(logger.New())
	return fs
}
