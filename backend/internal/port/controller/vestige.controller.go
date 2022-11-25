package controller

import (
	"clinker-backend/common/logger"
	"clinker-backend/internal/domain/model/dto"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/domain/service"

	"github.com/gofiber/fiber/v2"
)

type VestigeController struct {
	router         fiber.Router
	vestigeService *service.VestigeService
	processService *service.ProcessService
	authService    *service.AuthService
}

func NewVestigeController(
	router fiber.Router,
	vestigeService *service.VestigeService,
	processService *service.ProcessService,
	authService *service.AuthService,
) *VestigeController {
	return &VestigeController{
		router:         router,
		vestigeService: vestigeService,
		processService: processService,
		authService:    authService,
	}
}

func (c *VestigeController) Accessible() {
	c.router.Get("/ancestors", c.findAncestors)
	c.router.Get("/:txHash", c.findOne)
	c.router.Get("/friends/:txHash", c.findFriends)
	c.router.Get("/children/:txHash", c.findChildren)
}

func (c *VestigeController) Restricted() {
	c.router.Post("/", c.save)
}

func (*VestigeController) name() string {
	return "VestigeController"
}

// @tags    Vestige
// @summary Inquire vestiges of the main page
// @produce json
// @success 200 {object} res.ProfuseRes[res.VestigeRes]
// @router  /api/vestiges/ancestors [get]
// @param   pagination query dto.QueryPaginationDto true "pagination data"
func (c *VestigeController) findAncestors(ctx *fiber.Ctx) error {
	pagination := new(dto.QueryPaginationDto)
	if err := c.processService.PreWithParser(ctx.QueryParser, pagination, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	vestiges, errRes := c.vestigeService.FindAncestors(*pagination.Page, pagination.Take)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), vestiges)
}

// @tags    Vestige
// @summary Find one vestige with txHash
// @produce json
// @success 200 {object} res.VestigeRes
// @router  /api/vestiges/{txHash} [get]
// @param   txHash path string true "txHash"
func (c *VestigeController) findOne(ctx *fiber.Ctx) error {
	param := new(dto.ParamTxHashDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	vestige, errRes := c.vestigeService.FindOneByTxHash(param.TxHash)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), vestige)
}

// @tags    Vestige
// @summary Find all friends with head txHash
// @produce json
// @success 200 {array} res.VestigeRes
// @router  /api/vestiges/friends/{txHash} [get]
// @param   txHash path string true "txHash"
func (c *VestigeController) findFriends(ctx *fiber.Ctx) error {
	param := new(dto.ParamTxHashDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	vestiges, errRes := c.vestigeService.FindFriendsByHead(param.TxHash)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), vestiges)
}

// @tags    Vestige
// @summary Find all children with head txHash
// @produce json
// @success 200 {object} res.ProfuseRes[res.VestigeRes]
// @router  /api/vestiges/children/{txHash} [get]
// @param   pagination query dto.QueryPaginationDto true "pagination data"
// @param   txHash     path  string                 true "txHash"
func (c *VestigeController) findChildren(ctx *fiber.Ctx) error {
	pagination := new(dto.QueryPaginationDto)
	if err := c.processService.PreWithParser(ctx.QueryParser, pagination, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	param := new(dto.ParamTxHashDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	vestiges, errRes := c.vestigeService.FindChildren(param.TxHash, *pagination.Page, pagination.Take)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), vestiges)
}

// @tags     Vestige
// @summary  Save new vestige
// @produce  json
// @success  200 {object} res.SaveTxHashRes
// @router   /api/vestiges [post]
// @param    vestige body dto.VestigeDto true "vestige"
// @security Authorization
func (c *VestigeController) save(ctx *fiber.Ctx) error {
	if err := c.authService.Available(ctx); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(err)
	}

	body := new(dto.VestigeDto)
	if err := c.processService.PreWithParser(ctx.BodyParser, body, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	vestige, errRes := c.vestigeService.Save(ctx.Locals("id").(string), body)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), vestige)
}
