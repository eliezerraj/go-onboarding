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
	"github.com/go-onboarding/internal/adapter/api"
	"github.com/go-onboarding/internal/adapter/database"

	go_core_pg "github.com/eliezerraj/go-core/database/pg"
	go_core_aws_config "github.com/eliezerraj/go-core/aws/aws_config"
	go_core_s3_bucket "github.com/eliezerraj/go-core/aws/bucket_s3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

var(
	logLevel = 	zerolog.InfoLevel // zerolog.InfoLevel zerolog.DebugLevel
	appServer	model.AppServer
	databaseConfig 		go_core_pg.DatabaseConfig
	databasePGServer 	go_core_pg.DatabasePGServer
	goCoreAwsConfig 	go_core_aws_config.AwsConfig
	goCoreAwsBucketS3	go_core_s3_bucket.AwsBucketS3

	childLogger = log.With().Str("component","go-onboarding").Str("package", "main").Logger()
)

// Above init
func init(){
	childLogger.Info().Str("func","init").Send()
	zerolog.SetGlobalLevel(logLevel)

	infoPod, server := configuration.GetInfoPod()
	configOTEL 		:= configuration.GetOtelEnv()
	databaseConfig 	:= configuration.GetDatabaseEnv()
	certsTls 		:= configuration.GetCertEnv() 

	awsService 		:= configuration.GetAwsServiceEnv() 

	appServer.InfoPod = &infoPod
	appServer.Server = &server
	appServer.ConfigOTEL = &configOTEL
	appServer.AwsService = &awsService
	appServer.Cert = &certsTls
	appServer.DatabaseConfig = &databaseConfig
}

// Above main
func main (){
	childLogger.Info().Str("func","main").Interface("appServer",appServer).Send()

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( appServer.Server.ReadTimeout ) * time.Second)
	defer cancel()

	// Open Database
	count := 1
	var err error
	for {
		databasePGServer, err = databasePGServer.NewDatabasePGServer(ctx, *appServer.DatabaseConfig)
		if err != nil {
			if count < 3 {
				log.Error().Err(err).Msg("error open database... trying again !!")
			} else {
				log.Error().Err(err).Msg("fatal error open Database aborting")
				panic(err)
			}
			time.Sleep(3 * time.Second) //backoff
			count = count + 1
			continue
		}
		break
	}

	// Prepare aws services
	awsConfig, err := goCoreAwsConfig.NewAWSConfig(ctx, appServer.AwsService.AwsRegion)
	if err != nil {
		panic("error create new aws session " + err.Error())
	}

	// Otel over aws services
	otelaws.AppendMiddlewares(&awsConfig.APIOptions)

	// Create a S3 worker
	s3BucketWorker := goCoreAwsBucketS3.NewAwsS3Bucket(awsConfig)

	// wire	
	database := database.NewWorkerRepository(&databasePGServer)
	workerService := service.NewWorkerService(database, s3BucketWorker, appServer.AwsService)
	httpRouters := api.NewHttpRouters(workerService)

	// start server
	httpServer := server.NewHttpAppServer(appServer.Server)
	httpServer.StartHttpAppServer(ctx, &httpRouters, &appServer)
}