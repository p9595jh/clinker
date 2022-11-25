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
	authService      *service.AuthService
}

func NewAppraisalController(
	router fiber.Router,
	appraisalService *service.AppraisalService,
	processService *service.ProcessService,
	authService *service.AuthService,
) *AppraisalController {
	return &AppraisalController{
		router:           router,
		appraisalService: appraisalService,
		processService:   processService,
		authService:      authService,
	}
}

func (c *AppraisalController) Accessible() {
	c.router.Get("/:txHash", c.findByVestigeHead)
}

func (c *AppraisalController) Restricted() {
	c.router.Get("/users/:userId", c.findByUserId)
	c.router.Post("/", c.save)
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
	param := new(dto.ParamTxHashDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	appraisal, errRes := c.appraisalService.FindByVestigeHead(param.TxHash)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), appraisal)
}

// @tags     Appraisal
// @summary  Find by given user id
// @produce  json
// @success  200 {object} res.ProfuseRes[res.AppraisalSpecificRes]
// @router   /api/appraisals/users/{userId} [get]
// @param    pagination query dto.QueryPaginationDto true "pagination data"
// @param    userId     path  string                 true "user id"
// @security Authorization
func (c *AppraisalController) findByUserId(ctx *fiber.Ctx) error {
	pagination := new(dto.QueryPaginationDto)
	if err := c.processService.PreWithParser(ctx.QueryParser, pagination, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	param := new(dto.ParamUserIdDto)
	if err := c.processService.Pre(ctx.AllParams(), param, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	appraisals, errRes := c.appraisalService.FindByUserId(*pagination.Page, pagination.Take, param.UserId)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), appraisals)
}

// @tags     Appraisal
// @summary  Save new appraisal
// @produce  json
// @success  200 {object} res.SaveTxHashRes
// @router   /api/appraisals [post]
// @param    appraisal body dto.AppraisalDto true "appraisal"
// @security Authorization
func (c *AppraisalController) save(ctx *fiber.Ctx) error {
	if err := c.authService.Available(ctx); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(err)
	}

	body := new(dto.AppraisalDto)
	if err := c.processService.PreWithParser(ctx.BodyParser, body, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(res.NewErrorClientRes(ctx, err.Error()))
	}

	appraisal, errRes := c.appraisalService.Save(ctx.Locals("id").(string), body)
	return res.New(errRes).JustReturn(ctx, logger.Error(c.name()), appraisal)
}
