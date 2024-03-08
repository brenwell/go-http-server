package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const Addr = ":80"

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))
	prometheus := fiberprometheus.New("buzz")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/ping", pingHandler)
	app.Get("/block", jsonHandler)
	app.Get("/api/self", jsonHandler)
	app.Get("/api/cluster/a", requestHandler(os.Getenv("NGINX_A")))
	app.Get("/api/cluster/b", requestHandler(os.Getenv("NGINX_B")))
	app.Get("/api/external/a", requestHandler("https://httpbin.org/json"))
	app.Get("/api/external/b", requestHandler("https://jsonplaceholder.typicode.com/todos/1"))


	log.Printf("Server started on port%s...", Addr)
	if err := app.Listen(Addr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func pingHandler(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func jsonHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
		"source": "self",
		"status":  "ok",
	})
}

func requestHandler(url string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if url == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Not a valid url")
		}
		resp, err := http.Get(url)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Send(body)
	}
}