package main

import (
	"clinker-backend/common/config"
	"clinker-backend/common/hook"
	"clinker-backend/common/logger"
	_ "clinker-backend/docs"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/domain/service"
	"clinker-backend/internal/infrastructure/database/dbconn"
	"clinker-backend/internal/infrastructure/database/repository"
	"clinker-backend/internal/port/controller"
	"clinker-backend/internal/port/web"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jasonlvhit/gocron"
	"github.com/p9595jh/transform"
)

// @securityDefinitions.apikey Authorization
// @in                         header
// @name                       Authorization
func main() {
	ctx := "Application"

	// show configuration
	logger.Info(ctx).D("config", config.ToMap()).W()

	// database
	logger.Info(ctx).W("DB setting...")
	db, err := dbconn.Connect(&dbconn.DatabaseConfig{
		User:     config.V.GetString("db.user"),
		Password: config.V.GetString("db.password"),
		Host:     config.V.GetString("db.host"),
		Port:     config.V.GetString("db.port"),
		Schema:   config.V.GetString("db.schema"),
	})
	if err != nil {
		panic(err)
	} else {
		defer dbconn.Close(db)
	}
	logger.Info(ctx).W("DB setting done")

	// repository
	var (
		vestigeRepository   = repository.NewVestigeRepository(db)
		appraisalRepository = repository.NewAppraisalRepository(db)
		userRepository      = repository.NewUserRepository(db)
	)

	logger.Info(ctx).W("Repositories loaded")

	// service
	var (
		processService   = service.NewProcessService(validator.New(), transform.New())
		authService      = service.NewAuthService(userRepository)
		vestigeService   = service.NewVestigeService(vestigeRepository, appraisalRepository, userRepository, processService)
		appraisalService = service.NewAppraisalService(appraisalRepository, vestigeRepository, processService)
		userService      = service.NewUserService(userRepository, processService, authService)
	)

	logger.Info(ctx).W("Services loaded")

	// hook
	services := []any{
		processService,
		authService,
		vestigeService,
		appraisalService,
		userService,
	}
	scheduleItems := []hook.ScheduleItem{}
	for _, s := range services {
		if init, ok := s.(hook.Initialize); ok {
			go init.Initializer()
		}
		if schedule, ok := s.(hook.Schedule); ok {
			scheduleItems = append(scheduleItems, schedule.Schedulers()...)
		}
	}
	go func() {
		for _, item := range scheduleItems {
			item.Period.Do(item.Job)
		}
		<-gocron.Start()
	}()
	logger.Info(ctx).W("Service hooks are executed")

	// fiber init
	var (
		w = web.NewWeb(fiber.New(fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(fiber.StatusNotFound).JSON(res.NewErrorClientRes(c, "Page Not Found"))
			},
		}), fmt.Sprintf(":%s", config.V.GetString("port")), authService.PK())
		api = w.App.Group("/api")
	)

	// controller
	var (
		appController       = controller.NewAppController(api)
		authController      = controller.NewAuthController(api.Group("/auth"), authService, processService)
		vestigeController   = controller.NewVestigeController(api.Group("/vestiges"), vestigeService, processService, authService)
		appraisalController = controller.NewAppraisalController(api.Group("/appraisals"), appraisalService, processService, authService)
		userController      = controller.NewUserController(api.Group("/users"), userService, processService)
	)

	controllers := []hook.Controller{
		appController,
		authController,
		vestigeController,
		appraisalController,
		userController,
	}
	w.Attach(controllers)

	logger.Info(ctx).W("Controllers loaded")

	logger.Info(ctx).Wf("Server is listening %s", w.Address)
	if err := w.Listen(); err != nil {
		logger.Error(ctx).E(err).W()
	}
}
