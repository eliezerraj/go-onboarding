package service

import(
	"context"

	"github.com/go-onboarding/internal/core/model"
	"github.com/go-onboarding/internal/core/erro"
	go_core_observ "github.com/eliezerraj/go-core/observability"
)

var tracerProvider go_core_observ.TracerProvider

func (s WorkerService) AddPerson(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Debug().Msg("AddPerson")
	childLogger.Debug().Interface("onboarding: ",onboarding).Msg("")

	span := tracerProvider.Span(ctx, "service.AddPerson")
	defer span.End()
	
	tx, conn, err := s.workerRepository.DatabasePGServer.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
		s.workerRepository.DatabasePGServer.ReleaseTx(conn)
		span.End()
	}()

	res, err := s.workerRepository.AddPerson(ctx, tx, onboarding)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s WorkerService) GetPerson(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Debug().Msg("GetPerson")
	childLogger.Debug().Interface("onboarding: ",onboarding).Msg("")

	span := tracerProvider.Span(ctx, "service.GetPerson")
	defer span.End()
	
	res, err := s.workerRepository.GetPerson(ctx, onboarding)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s WorkerService) UpdatePerson(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Debug().Msg("UpdatePerson")
	childLogger.Debug().Interface("onboarding: ",onboarding).Msg("")

	span := tracerProvider.Span(ctx, "service.UpdatePerson")
	defer span.End()
	
	tx, conn, err := s.workerRepository.DatabasePGServer.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
		s.workerRepository.DatabasePGServer.ReleaseTx(conn)
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

func (s WorkerService) ListPerson(ctx context.Context, onboarding *model.Onboarding) (*[]model.Onboarding, error){
	childLogger.Debug().Msg("ListPerson")
	childLogger.Debug().Interface("onboarding: ",onboarding).Msg("")

	span := tracerProvider.Span(ctx, "service.ListPerson")
	defer span.End()
	
	res, err := s.workerRepository.ListPerson(ctx, onboarding)
	if err != nil {
		return nil, err
	}
	return res, nil
}