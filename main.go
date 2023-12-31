package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
)

func main() {

	ctx := context.Background()
	dbString := envString("DB_URL", "postgresql://postgres:password@localhost:5432/kontrakan?sslmode=disable")
	portString := envString("PORT", "3000")
	pool, err := pgxpool.New(ctx, dbString)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database")
	}
	defer pool.Close()

	m := newMigration(pool)

	cacheConf := cache.New(5*time.Minute, 10*time.Minute)

	app := fiber.New()

	//log all request
	app.Use(func(c *fiber.Ctx) error {
		log.Trace().Msgf("Method: %s \t Path: %s", c.Method(), c.Path())
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//setup migration setting
	app.Post("/migrate/status", func(c *fiber.Ctx) error {
		var user MigratorStatus
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(Response{
				Status:  400,
				Message: "Invalid request",
			})
		}
		u := user.validateUser()
		if len(u) > 0 {
			log.Info().Msg("User not valid while try to turn on migration")
			return c.JSON(Response{
				Status:     400,
				Message:    "Migration failed",
				FormErrors: u,
			})
		}

		log.Info().Interface("user", user).Msg("User valid while try to turn on migration")
		if user.Status {
			//set cache migrate to 1. it means migration is on
			cacheConf.Set("migrate", 1, cache.DefaultExpiration)
			return c.JSON(MigrateResponse{
				Response: Response{
					Status:  200,
					Message: "Migration on success",
				},
				Data: user.Status,
			})
		}

		//set cache migrate to 0. it means migration is off
		cacheConf.Set("migrate", 0, cache.DefaultExpiration)
		return c.JSON(MigrateResponse{
			Response: Response{
				Status:  200,
				Message: "Migration off success",
			},
			Data: user.Status,
		})
	})

	api := app.Group("/api/v1")
	api.Post("/migrate/", func(c *fiber.Ctx) error {

		getCache, ok := cacheConf.Get("migrate")
		if !ok {
			log.Info().Msg("cache not found")
			return c.JSON(MigrateResponse{
				Response: Response{
					Status:  200,
					Message: "migration failed",
				},
				Data: ok,
			})
		}

		if getCache == 0 {
			return c.JSON(MigrateResponse{
				Response: Response{
					Status:  200,
					Message: "migration failed",
				},
				Data: false,
			})
		}

		err = m.migrate(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Migration failed")
			return c.JSON(MigrateResponse{
				Response: Response{
					Status:  200,
					Message: "Migration failed",
				},
				Data: false,
			})
		}

		return c.JSON(MigrateResponse{
			Response: Response{
				Status:  200,
				Message: "Migration success",
			},
			Data: true,
		})
	})

	// report endpoint
	api.Get("/reports", func(c *fiber.Ctx) error {
		var reports []Report
		rows, err := pool.Query(ctx, "SELECT * FROM reports")
		if err != nil {
			log.Error().Err(err).Msg("Unable to get reports")
			return c.JSON(Response{
				Status:  500,
				Message: "Unable to get reports",
			})
		}
		defer rows.Close()

		for rows.Next() {
			var report Report
			err = rows.Scan(
				&report.Id,
				&report.Title,
				&report.Type,
				&report.Description,
				&report.Image,
				&report.UserReport,
				&report.Price,
				&report.UserId,
				&report.CreatedAt,
				&report.UpdatedAt,
			)
			if err != nil {
				log.Error().Err(err).Msg("Unable to scan report")
				return c.JSON(Response{
					Status:  500,
					Message: "Unable to scan report",
				})
			}
			reports = append(reports, report)
		}

		return c.JSON(ListResReport{
			Response: Response{
				Status:  200,
				Message: "Reports retrieved successfully",
			},
			Data: reports,
		})
	})
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

		newReport := NewReportConnection(pool)

		res, err := newReport.Save(ctx, report)
		if err != nil {
			log.Error().Err(err).Msg("Unable to save report")
			return c.JSON(Response{
				Status:  500,
				Message: "Unable to save report",
			})
		}

		return c.JSON(ResponseReport{
			Response: Response{
				Status:  200,
				Message: "Report created successfully",
			},
			Data: res,
		})
	})

	log.Fatal().Err(app.Listen(fmt.Sprintf(":%s", portString))).Msg("Server is running")
}
