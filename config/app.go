package config

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"latihan-fiber/helper"
	"latihan-fiber/middleware"
	"latihan-fiber/route"
)

func NewApp(logger *slog.Logger, deps route.Dependencies) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      GetEnv("APP_NAME", "Praktikum Backend Lanjut"),
		ErrorHandler: newErrorHandler(logger),
		// Membatasi ukuran body mencegah satu request besar menghabiskan memori server
		BodyLimit:    1 * 1024 * 1024, // 1 MB
	})

	// Panggil middleware.Register dengan allowedOrigins dari .env
	middleware.Register(app, logger, GetEnv("ALLOWED_ORIGINS", ""))

	route.Register(app, deps)

	return app
}

func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		logger.Error("terjadi panic atau error tidak tertangani",
			slog.String("error", err.Error()),
			slog.String("path", c.Path()),
		)
		return helper.Fail(c, code, "terjadi kesalahan pada server")
	}
}