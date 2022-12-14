package web

import (
	"crypto/rsa"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v4"

	"clinker-backend/common/hook"
	"clinker-backend/common/util"
	_ "clinker-backend/docs"
	"clinker-backend/internal/domain/model/res"
)

type Web struct {
	App        *fiber.App
	Address    string
	privateKey *rsa.PrivateKey
}

func NewWeb(app *fiber.App, address string, privateKey *rsa.PrivateKey) *Web {
	w := &Web{app, address, privateKey}
	p := &processor{}

	w.App.Use(cors.New())
	w.App.Use(limiter.New(limiter.Config{
		Expiration: 60 * time.Second,
		Max:        1000,
	}))
	w.App.Use(func(c *fiber.Ctx) error {
		return fiberlog.New(fiberlog.Config{
			Format:     reqLogFormat,
			TimeFormat: util.DateFormat,
			Output: &writer{
				keys: []string{
					"status_code", "http_method", "request_uri", "request_params",
					"request_body", "response_body", "response_time", "request_ip",
				},
				processors: map[string]func([]byte) any{
					"response_body":  p.responseBody,
					"request_body":   p.requestBody,
					"request_params": p.requestParams,
					"response_time":  p.responseTime,
					"status_code":    p.statusCode,
					"request_uri":    p.requestUri,
				},
				id: func() string {
					token, ok := c.Locals("user").(*jwt.Token)
					if !ok {
						return ""
					}
					claim, ok := token.Claims.(jwt.MapClaims)
					if !ok {
						return ""
					}
					id, ok := claim["id"]
					if !ok {
						return ""
					}
					if s, ok := id.(string); !ok {
						return ""
					} else {
						return s
					}
				}(),
			},
		})(c)
	})

	w.App.Get("/swagger/*", swagger.HandlerDefault)

	return w
}

func (w *Web) Attach(controllers []hook.Controller) {
	w.App.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	for _, c := range controllers {
		c.Accessible()
	}

	w.App.Use(jwtware.New(jwtware.Config{
		SigningMethod: "RS256",
		SigningKey:    w.privateKey,
		SuccessHandler: func(c *fiber.Ctx) error {
			claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
			for _, field := range []string{
				"id",
				"authority",
				"confirmed",
				"availableAt",
			} {
				c.Locals(field, claims[field])
			}
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(res.NewErrorClientRes(c, "unauthorized"))
		},
	}))

	for _, c := range controllers {
		c.Restricted()
	}
}

func (w *Web) Listen() <-chan error {
	ch := make(chan error)
	go func() { ch <- w.App.Listen(w.Address) }()
	return ch
}
