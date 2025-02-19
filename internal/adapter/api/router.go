package api

import (
	"encoding/json"
	"net/http"
	"github.com/rs/zerolog/log"

	"github.com/go-onboarding/internal/core/service"
	"github.com/go-onboarding/internal/core/model"
	"github.com/go-onboarding/internal/core/erro"
	go_core_observ "github.com/eliezerraj/go-core/observability"
	"github.com/eliezerraj/go-core/coreJson"
	"github.com/gorilla/mux"
)

var childLogger = log.With().Str("adapter", "api.router").Logger()

var core_json coreJson.CoreJson
var core_apiError coreJson.APIError
var tracerProvider go_core_observ.TracerProvider

type HttpRouters struct {
	workerService 	*service.WorkerService
}

func NewHttpRouters(workerService *service.WorkerService) HttpRouters {
	return HttpRouters{
		workerService: workerService,
	}
}

func (h *HttpRouters) Health(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Health")

	health := true
	json.NewEncoder(rw).Encode(health)
}

func (h *HttpRouters) Live(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Live")

	live := true
	json.NewEncoder(rw).Encode(live)
}

func (h *HttpRouters) Header(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Header")
	
	json.NewEncoder(rw).Encode(req.Header)
}

func (h *HttpRouters) AddPerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("AddPerson")

	span := tracerProvider.Span(req.Context(), "adapter.api.AddPerson")
	defer span.End()

	onBoarding := model.Onboarding{}
	err := json.NewDecoder(req.Body).Decode(&onBoarding)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(erro.ErrUnmarshal, http.StatusBadRequest)
		return &core_apiError
    }
	defer req.Body.Close()

	res, err := h.workerService.AddPerson(req.Context(), &onBoarding)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

func (h *HttpRouters) GetPerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("GetPerson")

	span := tracerProvider.Span(req.Context(), "adapter.api.GetPerson")
	defer span.End()

	vars := mux.Vars(req)
	varID := vars["id"]

	onBoarding := model.Onboarding{}
	person := model.Person{}
	person.PersonID = varID
	onBoarding.Person = &person

	res, err := h.workerService.GetPerson(req.Context(), &onBoarding)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

func (h *HttpRouters) UpdatePerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("UpdatePerson")

	span := tracerProvider.Span(req.Context(), "adapter.api.UpdatePerson")
	defer span.End()

	onBoarding := model.Onboarding{}
	err := json.NewDecoder(req.Body).Decode(&onBoarding)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(erro.ErrUnmarshal, http.StatusBadRequest)
		return &core_apiError
    }
	defer req.Body.Close()

	res, err := h.workerService.UpdatePerson(req.Context(), &onBoarding)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

func (h *HttpRouters) ListPerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("ListPerson")

	span := tracerProvider.Span(req.Context(), "adapter.api.ListPerson")
	defer span.End()

	vars := mux.Vars(req)
	varID := vars["id"]

	onBoarding := model.Onboarding{}
	person := model.Person{}
	person.PersonID = varID
	onBoarding.Person = &person

	res, err := h.workerService.ListPerson(req.Context(), &onBoarding)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}