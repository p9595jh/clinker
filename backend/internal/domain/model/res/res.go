package res

import (
	"clinker-backend/common/logger"

	"github.com/gofiber/fiber/v2"
)

// type NormalCaseFunc func(int) error
type ErrorCaseFunc func(int, error) error

type Response struct {
	errRes     *ErrorRes
	normal     func() error
	err4, err5 ErrorCaseFunc
}

func New(errRes *ErrorRes) *Response {
	return &Response{
		errRes: errRes,
	}
}

func (r *Response) Normal(f func() error) *Response {
	r.normal = f
	return r
}

func (r *Response) Error4(f ErrorCaseFunc) *Response {
	r.err4 = f
	return r
}

func (r *Response) Error5(f ErrorCaseFunc) *Response {
	r.err5 = f
	return r
}

func (r *Response) Return() error {
	if r.normal != nil {
		return r.normal()
	} else {
		switch r.errRes.Status / 100 {
		case 4:
			return r.err4(r.errRes.Status, r.errRes.Error)
		case 5:
			return r.err5(r.errRes.Status, r.errRes.Error)
		}
	}
	return nil
}

func (r *Response) JustReturn(ctx *fiber.Ctx, errLogItem logger.LogItem, data any) error {
	if r.errRes != nil {
		switch r.errRes.Status / 100 {
		case 4:
			return ctx.Status(r.errRes.Status).JSON(NewErrorClientRes(ctx, r.errRes.String()))
		case 5:
			errLogItem.E(r.errRes.Error).W()
			return New500Res(ctx)
		}
	}
	return ctx.JSON(data)
}
