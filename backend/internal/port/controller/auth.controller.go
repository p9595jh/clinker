package controller

import (
	"clinker-backend/common/logger"
	"clinker-backend/internal/domain/model/dto"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/domain/service"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	router         fiber.Router
	authService    *service.AuthService
	processService *service.ProcessService
}

func NewAuthController(
	router fiber.Router,
	authService *service.AuthService,
	processService *service.ProcessService,
) *AuthController {
	return &AuthController{
		router:         router,
		authService:    authService,
		processService: processService,
	}
}

func (c *AuthController) Restricted() {}

func (c *AuthController) Accessible() {
	c.router.Post("/login", c.login)
}

func (c *AuthController) name() string {
	return "AuthController"
}

func (c *AuthController) login(ctx *fiber.Ctx) error {
	var body dto.AuthLoginDto
	if err := c.processService.PreWithParser(ctx.BodyParser, &body, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	ok, user, err := c.authService.Validate(body.Id, body.Password)
	if err != nil {
		logger.Error(c.name()).E(err).W()
		return res.New500Res(ctx)
	} else if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(res.NewErrorClientRes(ctx, "id or password is wrong"))
	}

	if token, err := c.authService.PublishToken(user); err != nil {
		logger.Error(c.name()).E(err).W()
		return res.New500Res(ctx)
	} else {
		return ctx.JSON(&res.AuthLoginRes{Token: token})
	}
}
