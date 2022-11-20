package controller

import (
	"github.com/gofiber/fiber/v2"
)

type AppController struct {
	router fiber.Router
}

func NewAppController(router fiber.Router) *AppController {
	return &AppController{
		router: router,
	}
}

func (c *AppController) Accessible() {
	c.router.Get("/health", c.healthCheck)
}

func (c *AppController) Restricted() {}

// @tags   App
// @router /api/health [get]
func (c *AppController) healthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}
