package service

import(
	"context"

	"github.com/go-onboarding/internal/core/model"
	"github.com/go-onboarding/internal/adapter/database"
	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("core", "service").Logger()

type WorkerInterface interface {
	AddPerson(context.Context, *model.Onboarding) (*model.Onboarding, error)
	GetPerson(context.Context, *model.Onboarding) (*model.Onboarding, error)
	UpdatePerson(context.Context, *model.Onboarding) (*model.Onboarding, error)
	ListPerson(context.Context, *model.Onboarding) (*[]model.Onboarding, error)
}

type WorkerService struct {
	workerRepository *database.WorkerRepository
}

func NewWorkerService(workerRepository *database.WorkerRepository) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
		workerRepository: workerRepository,
	}
}