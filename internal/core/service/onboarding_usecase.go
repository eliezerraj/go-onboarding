package service

import(
	"context"

	"github.com/go-onboarding/internal/adapter/database"
	"github.com/rs/zerolog/log"

	"github.com/go-onboarding/internal/core/model"
	"github.com/go-onboarding/internal/core/erro"

	go_core_pg "github.com/eliezerraj/go-core/database/pg"
	go_core_observ "github.com/eliezerraj/go-core/observability"
	go_core_s3_bucket "github.com/eliezerraj/go-core/aws/bucket_s3"
)

var tracerProvider go_core_observ.TracerProvider
var childLogger = log.With().Str("component","go-payment").Str("package","internal.core.service").Logger()

type WorkerService struct {
	workerRepository 	*database.WorkerRepository
	workerBucketS3 		*go_core_s3_bucket.AwsBucketS3
	awsService			*model.AwsService
}

// About create a new worker service
func NewWorkerService(	workerRepository *database.WorkerRepository,
						workerBucketS3 *go_core_s3_bucket.AwsBucketS3,
						awsService		*model.AwsService) *WorkerService{
	childLogger.Info().Str("func","NewWorkerService").Send()

	return &WorkerService{
		workerRepository: workerRepository,
		workerBucketS3: workerBucketS3,
		awsService: awsService,
	}
}

// About handle/convert http status code
func (s *WorkerService) Stat(ctx context.Context) (go_core_pg.PoolStats){
	childLogger.Info().Str("func","Stat").Interface("trace-resquest-id", ctx.Value("trace-request-id")).Send()

	return s.workerRepository.Stat(ctx)
}

// About create a person
func (s *WorkerService) AddPerson(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Info().Str("func","AddPerson").Interface("trace-resquest-id", ctx.Value("trace-request-id")).Interface("onboarding", onboarding).Send()

	span := tracerProvider.Span(ctx, "service.AddPerson")
	defer span.End()
	
	tx, conn, err := s.workerRepository.DatabasePGServer.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.workerRepository.DatabasePGServer.ReleaseTx(conn)

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
		span.End()
	}()

	res, err := s.workerRepository.AddPerson(ctx, tx, onboarding)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// About get a person
func (s *WorkerService) GetPerson(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Info().Str("func","GetPerson").Interface("trace-resquest-id", ctx.Value("trace-request-id")).Interface("onboarding", onboarding).Send()

	span := tracerProvider.Span(ctx, "service.GetPerson")
	defer span.End()
	
	res, err := s.workerRepository.GetPerson(ctx, onboarding)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// About update a person
func (s *WorkerService) UpdatePerson(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Info().Str("func","UpdatePerson").Interface("trace-resquest-id", ctx.Value("trace-request-id")).Interface("onboarding", onboarding).Send()

	span := tracerProvider.Span(ctx, "service.UpdatePerson")
	defer span.End()
	
	tx, conn, err := s.workerRepository.DatabasePGServer.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.workerRepository.DatabasePGServer.ReleaseTx(conn)

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
		span.End()
	}()

	//Check data exists
	_, err = s.workerRepository.GetPerson(ctx, onboarding)
	if err != nil {
		return nil, err
	}

	// Do update
	res_update, err := s.workerRepository.UpdatePerson(ctx, tx, onboarding)
	if err != nil {
		return nil, err
	}
	if (res_update == 0) {
		return nil, erro.ErrUpdate
	}

	return onboarding, nil
}

// About list a person
func (s *WorkerService) ListPerson(ctx context.Context, onboarding *model.Onboarding) (*[]model.Onboarding, error){
	childLogger.Info().Str("func","ListPerson").Interface("trace-resquest-id", ctx.Value("trace-request-id")).Interface("onboarding", onboarding).Send()

	span := tracerProvider.Span(ctx, "service.ListPerson")
	defer span.End()
	
	res, err := s.workerRepository.ListPerson(ctx, onboarding)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// About upload file
func (s *WorkerService) UploadFile(ctx context.Context, onboardingFile *model.OnboardingFile) (error){
	childLogger.Info().Str("func","UploadFile").Interface("trace-resquest-id", ctx.Value("trace-request-id")).Send()

	span := tracerProvider.Span(ctx, "service.UploadFile")
	defer span.End()
	
	onboardingFile.BucketName = s.awsService.BucketName
	onboardingFile.FilePath = s.awsService.FilePath

	err := s.workerBucketS3.PutObject(	ctx, 
										onboardingFile.BucketName,
										onboardingFile.FilePath, 
										onboardingFile.FileName,
										onboardingFile.File)
	if err != nil {
		return err
	}

	return nil
}