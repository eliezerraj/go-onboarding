package service

import(
	"context"

	"github.com/rs/zerolog/log"

	"github.com/go-onboarding/internal/core/model"
)

var childLogger = log.With().Str("core", "service").Logger()

func (s WorkerService) Add(ctx context.Context, onboarding *model.Onboarding) (*model.Onboarding, error){
	childLogger.Debug().Msg("Add")
	childLogger.Debug().Interface("onboarding: ",onboarding).Msg("")

	res := onboarding
	return res, nil
}