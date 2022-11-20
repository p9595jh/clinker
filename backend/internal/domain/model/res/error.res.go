package res

import (
	"clinker-backend/common/util"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type ErrorContext interface {
	Path(...string) string
	Method(...string) string
	Response() *fasthttp.Response
}

type ErrorClientRes struct {
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Message   string `json:"message"`
	Status    int    `json:"status"`
}

func NewErrorClientRes(ctx ErrorContext, message string, a ...any) *ErrorClientRes {
	return &ErrorClientRes{
		Timestamp: time.Now().Format(util.DateFormat),
		Method:    ctx.Method(),
		Path:      ctx.Path(),
		Message:   fmt.Sprintf(message, a...),
		Status:    ctx.Response().StatusCode(),
	}
}

type ErrorServerRes struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
}

func NewErrorServerRes(ctx ErrorContext) *ErrorServerRes {
	return &ErrorServerRes{
		Timestamp: time.Now().Format(util.DateFormat),
		Status:    ctx.Response().StatusCode(),
	}
}

func New500Res(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(NewErrorServerRes(ctx))
}

type ErrorRes struct {
	Status int
	Error  error
}

func NewErrorRes(status int, err error) *ErrorRes {
	return &ErrorRes{
		Status: status,
		Error:  err,
	}
}

func NewErrorfRes(status int, message string, a ...any) *ErrorRes {
	return &ErrorRes{
		Status: status,
		Error:  fmt.Errorf(message, a...),
	}
}

func NewInternalErrorRes(err error) *ErrorRes {
	return &ErrorRes{
		Status: fiber.StatusInternalServerError,
		Error:  err,
	}
}

func (e *ErrorRes) String() string {
	return e.Error.Error()
}
