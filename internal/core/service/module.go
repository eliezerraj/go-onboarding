package service

type WorkerService struct {
}

func NewWorkerService() *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
	}
}