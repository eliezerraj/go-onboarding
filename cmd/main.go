package main

import(
	"time"
	"context"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-onboarding/internal/infra/configuration"
	"github.com/go-onboarding/internal/core/model"
	"github.com/go-onboarding/internal/core/service"
	"github.com/go-onboarding/internal/infra/server"
	"github.com/go-onboarding/internal/handler/api"
	"github.com/eliezerraj/go-core/core"
)

var(
	logLevel = 	zerolog.DebugLevel
	appServer	model.AppServer
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)

	infoPod, server := configuration.GetInfoPod()
	configOTEL := configuration.GetOtelEnv()

	appServer.InfoPod = &infoPod
	appServer.Server = &server
	appServer.ConfigOTEL = &configOTEL
}

func main (){
	log.Debug().Msg("----------------------------------------------------")
	log.Debug().Msg("main")
	log.Debug().Msg("----------------------------------------------------")
	log.Debug().Interface("appServer :",appServer).Msg("")
	log.Debug().Msg("----------------------------------------------------")

	var core core.ToolsCore
	core.Test()

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( appServer.Server.ReadTimeout ) * time.Second)
	defer cancel()

	workerService := service.NewWorkerService()
	httpRouters := api.NewHttpRouters(workerService)
	httpServer := server.NewHttpAppServer(appServer.Server)
	httpServer.StartHttpAppServer(ctx, &httpRouters, &appServer)
}