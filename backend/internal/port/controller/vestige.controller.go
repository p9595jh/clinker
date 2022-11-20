package controller

import (
	"clinker-backend/common/logger"
	"clinker-backend/internal/domain/model/dto"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/domain/service"

	"github.com/gofiber/fiber/v2"
)

type VestigeController struct {
	router            fiber.Router
	vestigeService    *service.VestigeService
	preprocessService *service.PreprocessService
}

func NewVestigeController(
	router fiber.Router,
	vestigeService *service.VestigeService,
	preprocessService *service.PreprocessService,
) *VestigeController {
	return &VestigeController{
		router:            router,
		vestigeService:    vestigeService,
		preprocessService: preprocessService,
	}
}

func (c *VestigeController) Accessible() {
	c.router.Get("/orphans", c.findOrphans)
	c.router.Get("/:txHash", c.findOne)
	c.router.Get("/friends/:txHash", c.findFriends)
	c.router.Get("/children/:txHash", c.findChildren)
}

func (c *VestigeController) Restricted() {}

func (*VestigeController) name() string {
	return "VestigeController"
}

// @tags    Vestige
// @summary Inquire vestiges of the main page
// @produce json
// @success 200 {object} res.VestigesRes
// @router  /api/vestiges/orphans [get]
// @param   pagination query dto.VestigePaginationDto true "pagination data"
func (c *VestigeController) findOrphans(ctx *fiber.Ctx) error {
	pagination := new(dto.VestigePanginationDto)
	if err := c.preprocessService.PipeParsing(ctx.QueryParser, pagination, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	vestiges, errRes := c.vestigeService.FindOrphans(pagination.Skip, pagination.Take)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), vestiges)
}

// @tags    Vestige
// @summary Find one vestige with txHash
// @produce json
// @success 200 {object} res.VestigeRes
// @router  /api/vestiges/{txHash} [get]
// @param   txHash path string true "txHash"
func (c *VestigeController) findOne(ctx *fiber.Ctx) error {
	param := new(dto.VestigeTxHashDto)
	if err := c.preprocessService.Pipe(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	vestige, errRes := c.vestigeService.FindOneByTxHash(param.TxHash)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), vestige)
}

// @tags    Vestige
// @summary Find all friends with head txHash
// @produce json
// @success 200 {object} res.VestigeRes
// @router  /api/vestiges/friends/{txHash} [get]
// @param   txHash path string true "txHash"
func (c *VestigeController) findFriends(ctx *fiber.Ctx) error {
	body := new(dto.VestigeTxHashDto)
	if err := c.preprocessService.PipeParsing(ctx.BodyParser, body, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	vestiges, errRes := c.vestigeService.FindFriendsByHead(body.TxHash)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), vestiges)
}

// @tags    Vestige
// @summary Find all children with head txHash
// @produce json
// @success 200 {object} res.VestigesRes
// @router  /api/vestiges/children/{txHash} [get]
// @param   pagination query dto.VestigePaginationDto true "pagination data"
// @param   txHash     path  string                   true "txHash"
func (c *VestigeController) findChildren(ctx *fiber.Ctx) error {
	pagination := new(dto.VestigePanginationDto)
	if err := c.preprocessService.PipeParsing(ctx.QueryParser, pagination, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	body := new(dto.VestigeTxHashDto)
	if err := c.preprocessService.PipeParsing(ctx.BodyParser, body, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	vestiges, errRes := c.vestigeService.FindChildren(body.TxHash, pagination.Skip, pagination.Take)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), vestiges)
}
