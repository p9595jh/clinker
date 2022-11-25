package controller

import (
	"clinker-backend/common/enum/Authority"
	"clinker-backend/common/logger"
	"clinker-backend/internal/domain/model/dto"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/domain/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	router         fiber.Router
	userService    *service.UserService
	processService *service.ProcessService
}

func NewUserController(
	router fiber.Router,
	userService *service.UserService,
	processService *service.ProcessService,
) *UserController {
	return &UserController{
		router:         router,
		userService:    userService,
		processService: processService,
	}
}

func (c *UserController) Accessible() {
	c.router.Post("/", c.register)
}

func (c *UserController) Restricted() {
	c.router.Get("/", c.findUsers)
	c.router.Get("/:userId", c.findOneUser)
	c.router.Get("/:address", c.findByAddress)
	c.router.Put("/stops/:userId", c.stop)
}

func (c *UserController) name() string {
	return "UserController"
}

// @tags     User
// @summary  Inquire users
// @produce  json
// @success  200 {object} res.ProfuseRes[res.UserRes]
// @router   /api/users [get]
// @param    pagination query dto.QueryPaginationDto true "pagination data"
// @security Authorization
func (c *UserController) findUsers(ctx *fiber.Ctx) error {
	pagination := new(dto.QueryPaginationDto)
	if err := c.processService.PreWithParser(ctx.QueryParser, pagination, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	users, errRes := c.userService.FindUsers(ctx.Locals("authority").(Authority.Type), *pagination.Page, pagination.Take)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), users)
}

// @tags     User
// @summary  Find one user with userId
// @produce  json
// @success  200 {object} res.UserRes
// @router   /api/users/{userId} [get]
// @param    userId path string true "userId"
// @security Authorization
func (c *UserController) findOneUser(ctx *fiber.Ctx) error {
	param := new(dto.ParamUserIdDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	user, errRes := c.userService.FindOneUser(param.UserId)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), user)
}

// @tags     User
// @summary  Find one user with address
// @produce  json
// @success  200 {object} res.UserRes
// @router   /api/users/{address} [get]
// @param    address path string true "address"
// @security Authorization
func (c *UserController) findByAddress(ctx *fiber.Ctx) error {
	param := new(dto.ParamAddressDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	user, errRes := c.userService.FindByAddress(param.Address)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), user)
}

// @tags     User
// @summary  Find one user with userId
// @produce  json
// @success  200 {object} res.UserIdRes
// @router   /api/users/{userId} [put]
// @param    userId path string           true "userId"
// @param    stop   body dto.UserStopDtom true "stop data"
// @security Authorization
func (c *UserController) stop(ctx *fiber.Ctx) error {
	param := new(dto.ParamUserIdDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	bodym := new(dto.UserStopDtom)
	body := new(dto.UserStopDto)
	if err := c.processService.PreWithParser(ctx.BodyParser, bodym, body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	stopRes, errRes := c.userService.Stop(ctx.Locals("authority").(Authority.Type), param.UserId, body.Reason, body.Date)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), stopRes)
}

// @tags    User
// @summary Save new User
// @produce json
// @success 200 {object} res.UserIdRes
// @router  /api/users [post]
// @param   user body dto.UserDto true "user"
func (c *UserController) register(ctx *fiber.Ctx) error {
	body := new(dto.UserDto)
	if err := c.processService.PreWithParser(ctx.BodyParser, body, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	user, errRes := c.userService.Register(body)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), user)
}
