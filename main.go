package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func main() {

	


	// migration endpoint
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/migrate/:name/:password", func(c *fiber.Ctx) error {
		name := c.Params("name")
		password := c.Params("password")

		user := User{
			Name:     name,
			Password: password,
		}
		u := user.validateUser()
		if len(u) > 0 {
			return c.JSON(Response{
				Status:     400,
				Message:    "Migration failed",
				FormErrors: u,
			})
		}

		return c.JSON(fiber.Map{
			"status":  200,
			"message": "Migration success",
		})
	})

	api := app.Group("/api/v1")
	api.Post("/reports", func(c *fiber.Ctx) error {
		var report Report
		if err := c.BodyParser(&report); err != nil {
			return c.Status(400).JSON(Response{
				Status:  400,
				Message: "Invalid request",
			})
		}

		r := report.validateReport()
		if len(r) > 0 {
			return c.JSON(Response{
				Status:     400,
				Message:    "Report creation failed",
				FormErrors: r,
			})
		}

		return c.JSON(ResponseReport{
			Response: Response{
				Status:  200,
				Message: "Report created successfully",
			},
			Data: report,
		})
	})

	log.Fatal().Err(app.Listen(":3000")).Msg("Server is running")
}

const (
	aji           = "aji"
	migrationOnPw = "kontrakan"
)

type ResponseReport struct {
	Response
	Data Report `json:"data"`
}
