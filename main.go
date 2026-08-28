package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"latihan-fiber/app/repository"
	"latihan-fiber/config"
	"latihan-fiber/database"
)

// requireJSON memastikan request memiliki tipe konten JSON (Berasal dari pertemuan sebelumnya)
func requireJSON(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet && c.Get("Content-Type") != "application/json" {
		return fail(c, fiber.StatusUnsupportedMediaType, "Hanya menerima JSON")
	}
	return c.Next()
}

func main() {
	// 1. Konfigurasi
	config.LoadEnv() // Memuat variabel environment

	// 2. Koneksi basis data
	pool, err := database.NewPool(context.Background()) // Membuat pool koneksi PostgreSQL[cite: 1]
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	// 3. Perakitan: pool -> repository -> handler
	userRepository := repository.NewUserRepository(pool) // Inject pool ke repository[cite: 1]
	userHandler := NewUserHandler(userRepository)        // Inject repository ke handler[cite: 1]

	// 4. Aplikasi
	app := fiber.New()
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		// Kesehatan layanan kini ikut bergantung pada basis data[cite: 1]
		if err := pool.Ping(ctx); err != nil {
			return fail(c, fiber.StatusServiceUnavailable, "database tidak dapat dihubungi")
		}
		return ok(c, "server dan database berjalan", nil)
	})

	u := api.Group("/users", requireJSON)
	u.Get("/", userHandler.List)
	u.Get("/:id", userHandler.Get)
	u.Post("/", userHandler.Create)
	u.Put("/:id", userHandler.Replace)
	u.Patch("/:id", userHandler.Patch)
	u.Delete("/:id", userHandler.Delete)

	port := config.GetEnv("APP_PORT", "3000")
	log.Fatal(app.Listen(":" + port))
}