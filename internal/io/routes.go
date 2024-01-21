package io

import (
	"em_test/internal/io/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, controller handlers.Controller) {
	app.Post("/create", controller.CreatePerson)
	app.Delete("/del/:id", controller.DeletePerson)
	app.Get("/people", controller.GetPeople)
	app.Patch("/update/:id", controller.UpdatePerson)
}
