package controller

import (
	"clinker-backend/common/logger"
	"clinker-backend/internal/domain/model/dto"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/domain/service"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	router            fiber.Router
	authService       *service.AuthService
	preprocessService *service.PreprocessService
}

func NewAuthController(
	router fiber.Router,
	authService *service.AuthService,
	preprocessService *service.PreprocessService,
) *AuthController {
	return &AuthController{
		router:            router,
		authService:       authService,
		preprocessService: preprocessService,
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
	if err := c.preprocessService.PipeParsing(ctx.BodyParser, &body, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	if ok, err := c.authService.Validate(body.Id, body.Password); err != nil {
		logger.Error(c.name()).E(err).W()
		return res.New500Res(ctx)
	} else if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(res.NewErrorClientRes(ctx, "id or password is wrong"))
	}

	if token, err := c.authService.PublishToken(body.Id); err != nil {
		logger.Error(c.name()).E(err).W()
		return res.New500Res(ctx)
	} else {
		return ctx.JSON(&res.AuthLoginRes{Token: token})
	}
}
