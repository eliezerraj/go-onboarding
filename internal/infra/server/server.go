package server

import (
	"time"
	"encoding/json"
	"net/http"
	"strconv"
	"os"
	"os/signal"
	"syscall"
	"context"

	"github.com/go-onboarding/internal/core/model"
	go_core_observ "github.com/eliezerraj/go-core/observability"  
	"github.com/go-onboarding/internal/adapter/api"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/eliezerraj/go-core/middleware"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var childLogger = log.With().Str("component","go-payment").Str("package","internal.infra.server").Logger()

var core_middleware middleware.ToolsMiddleware
var tracerProvider go_core_observ.TracerProvider
var infoTrace go_core_observ.InfoTrace

type HttpServer struct {
	httpServer	*model.Server
}

func NewHttpAppServer(httpServer *model.Server) HttpServer {
	childLogger.Info().Str("func","NewHttpAppServer").Send()
	return HttpServer{httpServer: httpServer }
}

// About start http server
func (h HttpServer) StartHttpAppServer(	ctx context.Context, 
										httpRouters *api.HttpRouters,
										appServer *model.AppServer) {
	childLogger.Info().Str("func","StartHttpAppServer").Send()
			
	// ---------------------- OTEL ---------------
	infoTrace.PodName = appServer.InfoPod.PodName
	infoTrace.PodVersion = appServer.InfoPod.ApiVersion
	infoTrace.ServiceType = "k8-workload"
	infoTrace.Env = appServer.InfoPod.Env
	infoTrace.AccountID = appServer.InfoPod.AccountID

	tp := tracerProvider.NewTracerProvider(	ctx, 
											appServer.ConfigOTEL, 
											&infoTrace)

	if tp != nil {
		otel.SetTextMapPropagator(xray.Propagator{})
		otel.SetTracerProvider(tp)
	}

	defer func() { 
		if tp != nil {
			err := tp.Shutdown(ctx)
			if err != nil{
				childLogger.Error().Err(err).Send()
			}
		}
		childLogger.Info().Msg("stop done !!!")
	}()
	
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(core_middleware.MiddleWareHandlerHeader)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/")
		json.NewEncoder(rw).Encode(appServer)
	})

	health := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    health.HandleFunc("/health", httpRouters.Health)

	live := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    live.HandleFunc("/live", httpRouters.Live)

	header := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    header.HandleFunc("/header", httpRouters.Header)

	wk_ctx := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    wk_ctx.HandleFunc("/context", httpRouters.Context)

	stat := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    stat.HandleFunc("/stat", httpRouters.Stat)
	
	myRouter.HandleFunc("/info", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Info().Str("HandleFunc","/info").Send()

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(appServer)
	})
	
	addPerson := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	addPerson.HandleFunc("/person/add", core_middleware.MiddleWareErrorHandler(httpRouters.AddPerson))		
	addPerson.Use(otelmux.Middleware("go-onboarding"))

	getPerson := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getPerson.HandleFunc("/person/{id}", core_middleware.MiddleWareErrorHandler(httpRouters.GetPerson))		
	getPerson.Use(otelmux.Middleware("go-onboarding"))

	updatePerson := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	updatePerson.HandleFunc("/person/update", core_middleware.MiddleWareErrorHandler(httpRouters.UpdatePerson))		
	updatePerson.Use(otelmux.Middleware("go-onboarding"))

	listPerson := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	listPerson.HandleFunc("/person/list/{id}", core_middleware.MiddleWareErrorHandler(httpRouters.ListPerson))		
	listPerson.Use(otelmux.Middleware("go-onboarding"))

	uploadFile := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	uploadFile.HandleFunc("/uploadFile", core_middleware.MiddleWareErrorHandler(httpRouters.UploadFile))		
	uploadFile.Use(otelmux.Middleware("go-onboarding"))

	srv := http.Server{
		Addr:         ":" +  strconv.Itoa(h.httpServer.Port),      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(h.httpServer.ReadTimeout) * time.Second,   
		WriteTimeout: time.Duration(h.httpServer.WriteTimeout) * time.Second,  
		IdleTimeout:  time.Duration(h.httpServer.IdleTimeout) * time.Second, 
	}

	childLogger.Info().Str("Service Port", strconv.Itoa(h.httpServer.Port)).Send()

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			childLogger.Error().Err(err).Msg("canceling http mux server !!!")
		}
	}()

	// Get SIGNALS
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig := <-ch

		switch sig {
		case syscall.SIGHUP:
			childLogger.Info().Msg("Received SIGHUP: reloading configuration...")
		case syscall.SIGINT, syscall.SIGTERM:
			childLogger.Info().Msg("Received SIGINT/SIGTERM termination signal. Exiting")
			return
		default:
			childLogger.Info().Interface("Received signal:", sig).Send()
		}
	}

	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		childLogger.Error().Err(err).Msg("warning dirty shutdown !!!")
		return
	}
}