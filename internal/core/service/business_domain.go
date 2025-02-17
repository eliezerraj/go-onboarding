package service

import(
	"context"

	"github.com/rs/zerolog/log"
	"github.com/go-onboarding/internal/core/model"
	go_core_observ "github.com/eliezerraj/go-core/observability"
)

var tracerProvider go_core_observ.TracerProvider

var childLogger = log.With().Str("core", "service").Logger()

func (s *WorkerService) Add(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Debug().Msg("Add")
	childLogger.Debug().Interface("onboarding: ",onboarding).Msg("")

	span := tracerProvider.Span(ctx, "service.Add")
	defer span.End()
	
	s.workerRepository.Add(ctx, onboarding)

	res := onboarding
	return res, nil
}