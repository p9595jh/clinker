package controller

import (
	"clinker-backend/common/logger"
	"clinker-backend/internal/domain/model/dto"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/domain/service"

	"github.com/gofiber/fiber/v2"
)

type AppraisalController struct {
	router           fiber.Router
	appraisalService *service.AppraisalService
	processService   *service.ProcessService
}

func NewAppraisalController(
	router fiber.Router,
	appraisalService *service.AppraisalService,
	processService *service.ProcessService,
) *AppraisalController {
	return &AppraisalController{
		router:           router,
		appraisalService: appraisalService,
		processService:   processService,
	}
}

func (c *AppraisalController) Accessible() {
	c.router.Get("/:txHash", c.findByVestigeHead)
}

func (c *AppraisalController) Restricted() {
	c.router.Get("/users/:userId", c.findByUserId)
}

func (*AppraisalController) name() string {
	return "AppraisalController"
}

// @tags    Appraisal
// @summary Find calculated appraisal with head's txHash
// @produce json
// @success 200 {object} res.AppraisalRes
// @router  /api/appraisals/{txHash} [get]
// @param   txHash path string true "txHash"
func (c *AppraisalController) findByVestigeHead(ctx *fiber.Ctx) error {
	param := new(dto.AppraisalTxHashDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	appraisal, errRes := c.appraisalService.FindByVestigeHead(param.TxHash)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), appraisal)
}

// @tags    Appraisal
// @summary Find by given user id
// @produce json
// @success 200 {object} res.AppraisalSpecificsRes
// @router  /api/appraisals/users/{userId} [get]
// @param   pagination query dto.VestigePaginationDto true "pagination data"
// @param   userId     path  string                   true "user id"
func (c *AppraisalController) findByUserId(ctx *fiber.Ctx) error {
	pagination := new(dto.AppraisalPanginationDto)
	if err := c.processService.PreWithParser(ctx.QueryParser, pagination, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	param := new(dto.AppraisalUserIdDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	appraisals, errRes := c.appraisalService.FindByUserId(pagination.Skip, pagination.Take, param.UserId)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), appraisals)
}
