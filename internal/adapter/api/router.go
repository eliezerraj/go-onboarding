package api

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"

	"github.com/rs/zerolog/log"

	"github.com/go-onboarding/internal/core/service"
	"github.com/go-onboarding/internal/core/model"
	"github.com/go-onboarding/internal/core/erro"

	"github.com/eliezerraj/go-core/coreJson"
	"github.com/gorilla/mux"

	go_core_observ "github.com/eliezerraj/go-core/observability"
)

var childLogger = log.With().Str("component", "go-onboarding").Str("package", "internal.adapter.api").Logger()

var core_json coreJson.CoreJson
var core_apiError coreJson.APIError
var tracerProvider go_core_observ.TracerProvider

type HttpRouters struct {
	workerService 	*service.WorkerService
}

// Above create routers
func NewHttpRouters(workerService *service.WorkerService) HttpRouters {
	childLogger.Info().Str("func","NewHttpRouters").Send()

	return HttpRouters{
		workerService: workerService,
	}
}

// About return a health
func (h *HttpRouters) Health(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Msg("Health")

	json.NewEncoder(rw).Encode(model.MessageRouter{Message: "true"})
}

// About return a live
func (h *HttpRouters) Live(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Live").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	json.NewEncoder(rw).Encode(model.MessageRouter{Message: "true"})
}

// About show all header received
func (h *HttpRouters) Header(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Header").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()
	
	json.NewEncoder(rw).Encode(req.Header)
}

// About add person
func (h *HttpRouters) AddPerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","AddPerson").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	span := tracerProvider.Span(req.Context(), "adapter.api.AddPerson")
	defer span.End()

	onBoarding := model.Onboarding{}
	err := json.NewDecoder(req.Body).Decode(&onBoarding)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
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

// About get person
func (h *HttpRouters) GetPerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","GetPerson").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

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

// About update person
func (h *HttpRouters) UpdatePerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","UpdatePerson").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	span := tracerProvider.Span(req.Context(), "adapter.api.UpdatePerson")
	defer span.End()

	onBoarding := model.Onboarding{}
	err := json.NewDecoder(req.Body).Decode(&onBoarding)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
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

// About list person
func (h *HttpRouters) ListPerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","ListPerson").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

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

// About list person
func (h *HttpRouters) UploadFile(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","UploadFile").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	// Trace
	span := tracerProvider.Span(req.Context(), "adapter.api.UploadFile")
	defer span.End()

	// Check the size
	err := req.ParseMultipartForm(20 << 20) //20Mb
	if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
		return &core_apiError
	}

	// Open a form
	file, handler, err := req.FormFile("file")
	if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
		return &core_apiError
	}
	defer file.Close()

	onboardingFile := model.OnboardingFile{}
	onboardingFile.FileName = handler.Filename
	onboardingFile.File, err = ioutil.ReadAll(file)
	if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
		return &core_apiError
	}

	childLogger.Info().Str("func","UploadFile").
						Interface("file_data", fmt.Sprintf("%v %v %v",handler.Header ,handler.Filename, handler.Size)).
						Send()

	err = h.workerService.UploadFile(req.Context(), &onboardingFile)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}

	return json.NewEncoder(rw).Encode(model.MessageRouter{Message: "true"})
}