package api

import (
	"fmt"
	"time"
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"io/ioutil"
	"strings"

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
	ctxTimeout		time.Duration
}

// Above create routers
func NewHttpRouters(workerService *service.WorkerService,
					ctxTimeout	time.Duration) HttpRouters {
	childLogger.Info().Str("func","NewHttpRouters").Send()

	return HttpRouters{
		workerService: workerService,
		ctxTimeout: ctxTimeout,
	}
}

// About return a health
func (h *HttpRouters) Health(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Health").Send()

	json.NewEncoder(rw).Encode(model.MessageRouter{Message: "true"})
}

// About return a live
func (h *HttpRouters) Live(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Live").Send()

	json.NewEncoder(rw).Encode(model.MessageRouter{Message: "true"})
}

// About show all header received
func (h *HttpRouters) Header(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Header").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()
	
	json.NewEncoder(rw).Encode(req.Header)
}

// About show all context values
func (h *HttpRouters) Context(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Context").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()
	
	contextValues := reflect.ValueOf(req.Context()).Elem()
	json.NewEncoder(rw).Encode(fmt.Sprintf("%v",contextValues))
}

// About show pgx stats
func (h *HttpRouters) Stat(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Stat").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()
	
	res := h.workerService.Stat(req.Context())

	json.NewEncoder(rw).Encode(res)
}

// About add person
func (h *HttpRouters) AddPerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","AddPerson").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	ctx, cancel := context.WithTimeout(req.Context(), h.ctxTimeout * time.Second)
    defer cancel()

	span := tracerProvider.Span(ctx, "adapter.api.AddPerson")
	defer span.End()

	trace_id := fmt.Sprintf("%v", ctx.Value("trace-request-id"))

	onBoarding := model.Onboarding{}
	err := json.NewDecoder(req.Body).Decode(&onBoarding)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusBadRequest)
		return &core_apiError
    }
	defer req.Body.Close()

	res, err := h.workerService.AddPerson(ctx, &onBoarding)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About get person
func (h *HttpRouters) GetPerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","GetPerson").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	ctx, cancel := context.WithTimeout(req.Context(), h.ctxTimeout * time.Second)
    defer cancel()

	span := tracerProvider.Span(ctx, "adapter.api.GetPerson")
	defer span.End()

	trace_id := fmt.Sprintf("%v", ctx.Value("trace-request-id"))

	vars := mux.Vars(req)
	varID := vars["id"]

	onBoarding := model.Onboarding{}
	person := model.Person{}
	person.PersonID = varID
	onBoarding.Person = &person

	res, err := h.workerService.GetPerson(ctx, &onBoarding)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About update person
func (h *HttpRouters) UpdatePerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","UpdatePerson").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	ctx, cancel := context.WithTimeout(req.Context(), h.ctxTimeout * time.Second)
    defer cancel()

	span := tracerProvider.Span(ctx, "adapter.api.UpdatePerson")
	defer span.End()

	trace_id := fmt.Sprintf("%v", ctx.Value("trace-request-id"))

	onBoarding := model.Onboarding{}
	err := json.NewDecoder(req.Body).Decode(&onBoarding)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusBadRequest)
		return &core_apiError
    }
	defer req.Body.Close()

	res, err := h.workerService.UpdatePerson(ctx, &onBoarding)
	if err != nil {

		if strings.Contains(err.Error(), "context deadline exceeded") {
    		err = erro.ErrTimeout
		} 

		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusNotFound)
		case erro.ErrTimeout:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusGatewayTimeout)
		default:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About list person
func (h *HttpRouters) ListPerson(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","ListPerson").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()
	
	ctx, cancel := context.WithTimeout(req.Context(), h.ctxTimeout * time.Second)
    defer cancel()

	span := tracerProvider.Span(ctx, "adapter.api.ListPerson")
	defer span.End()

	trace_id := fmt.Sprintf("%v", ctx.Value("trace-request-id"))

	vars := mux.Vars(req)
	varID := vars["id"]

	onBoarding := model.Onboarding{}
	person := model.Person{}
	person.PersonID = varID
	onBoarding.Person = &person

	res, err := h.workerService.ListPerson(ctx, &onBoarding)
	if err != nil {

		if strings.Contains(err.Error(), "context deadline exceeded") {
    		err = erro.ErrTimeout
		} 

		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusNotFound)
		case erro.ErrTimeout:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusGatewayTimeout)
		default:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About list person
func (h *HttpRouters) UploadFile(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","UploadFile").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	ctx, cancel := context.WithTimeout(req.Context(), h.ctxTimeout * time.Second)
    defer cancel()

	// Trace
	span := tracerProvider.Span(ctx, "adapter.api.UploadFile")
	defer span.End()

	trace_id := fmt.Sprintf("%v",ctx.Value("trace-request-id"))

	// Check the size
	err := req.ParseMultipartForm(20 << 20) //20Mb
	if err != nil {
		core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusBadRequest)
		return &core_apiError
	}

	// Open a form
	file, handler, err := req.FormFile("file")
	if err != nil {
		core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusBadRequest)
		return &core_apiError
	}
	defer file.Close()

	onboardingFile := model.OnboardingFile{}
	onboardingFile.FileName = handler.Filename
	onboardingFile.File, err = ioutil.ReadAll(file)
	if err != nil {
		core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusBadRequest)
		return &core_apiError
	}

	childLogger.Info().Str("func","UploadFile").
						Interface("file_data", fmt.Sprintf("%v %v %v",handler.Header ,handler.Filename, handler.Size)).
						Send()

	err = h.workerService.UploadFile(ctx, &onboardingFile)
	if err != nil {

		if strings.Contains(err.Error(), "context deadline exceeded") {
    		err = erro.ErrTimeout
		} 

		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusNotFound)
		case erro.ErrTimeout:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusGatewayTimeout)
		default:
			core_apiError = core_apiError.NewAPIError(err, trace_id, http.StatusInternalServerError)
		}
		return &core_apiError
	}

	return json.NewEncoder(rw).Encode(model.MessageRouter{Message: "true"})
}