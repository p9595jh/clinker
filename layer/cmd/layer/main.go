package main

import (
	"fmt"
	"layer/common/config"
	"layer/common/hook"
	"layer/common/logger"
	"layer/common/util"
	"layer/internal/domain/service"
	"layer/internal/infrastructure/database/dbconn"
	"layer/internal/infrastructure/database/repository"
	rpcclient "layer/internal/infrastructure/rpc/client"
	rpcserver "layer/internal/infrastructure/rpc/server"
	"net/http"

	"github.com/jasonlvhit/gocron"
)

func main() {
	ctx := "Application"
	logger.FileLoggerInit()

	// health checker
	go func() {
		res := []byte(`{"success":true}`)
		http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) { w.Write(res) })
		logger.Info(ctx).Wf("HTTP Server Listeneing at %d", config.V.GetInt("port.http"))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.V.GetInt("port.http")), nil); err != nil {
			panic(err)
		}
	}()

	// show configuration
	logger.Info(ctx).D("config", config.ToMap()).W()

	// database
	logger.Info(ctx).W("DB setting...")
	db, err := dbconn.ConnectDB(&dbconn.DatabaseConfig{
		Host:   config.V.GetString("db.host"),
		Port:   config.V.GetInt("db.port"),
		Schema: config.V.GetString("db.schema"),
	})
	if err != nil {
		panic(err)
	}

	// transformer for repositories
	transformer := util.Transformer()

	// repository
	var (
		caRepository    = repository.NewCaRepository(db.Collection("cas"), transformer)
		clinkRepository = repository.NewClinkRepository(db.Collection("clinks"), transformer)
		txnRepository   = repository.NewTxnRepository(db.Collection("txns"), transformer)
	)

	logger.Info(ctx).W("Repositories loaded")

	// rpc client
	var client = rpcclient.NewClinkRpcClient(config.V.GetString("backend.rpc.url"))

	// service
	var transactionService = service.NewTransactionService(
		client,
		clinkRepository,
		caRepository,
		txnRepository,
	)

	logger.Info(ctx).W("Services loaded")

	// hook
	services := []any{
		transactionService,
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

	// rpc server
	var server = rpcserver.NewClinkRpcServer(
		config.V.GetInt("port.rpc"),
		transactionService,
	)

	logger.Info(ctx).Wf("RPC Server Listeneing at %d", config.V.GetInt("port.rpc"))
	for err := range server.Listen() {
		panic(err)
	}
}
