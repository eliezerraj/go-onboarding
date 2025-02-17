package service

import(
	"github.com/go-onboarding/internal/adapter/database"
)

type WorkerService struct {
	workerRepository *database.WorkerRepository
}

func NewWorkerService(workerRepository *database.WorkerRepository) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
		workerRepository: workerRepository,
	}
}