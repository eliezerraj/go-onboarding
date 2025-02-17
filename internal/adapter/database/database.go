package database

import (
	"context"
	//"time"
	//"errors"
	
	go_core_observ "github.com/eliezerraj/go-core/observability"
	"github.com/go-onboarding/internal/core/model"
	go_core_pg "github.com/eliezerraj/go-core/database/pg"

	"github.com/rs/zerolog/log"
)

var tracerProvider go_core_observ.TracerProvider
var childLogger = log.With().Str("adapter", "database").Logger()

type WorkerRepository struct {
	databasePGServer *go_core_pg.DatabasePGServer
}

func NewWorkerRepository(databasePGServer *go_core_pg.DatabasePGServer) *WorkerRepository{
	childLogger.Debug().Msg("NewWorkerRepository")

	return &WorkerRepository{
		databasePGServer: databasePGServer,
	}
}

func (w WorkerRepository) Add(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Debug().Msg("Add")

	span := tracerProvider.Span(ctx, "database.Add")
	defer span.End()

	return onboarding, nil
}