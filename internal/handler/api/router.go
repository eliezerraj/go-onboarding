package api

import (
	"encoding/json"
	"net/http"
	"github.com/rs/zerolog/log"

	"github.com/go-onboarding/internal/core/service"
	"github.com/go-onboarding/internal/core/model"
	"github.com/go-onboarding/internal/core/erro"
	"github.com/go-onboarding/internal/infra/observ"
	"github.com/eliezerraj/go-core/coreJson"
	//"github.com/gorilla/mux"
)

var childLogger = log.With().Str("handler", "api.router").Logger()
var core_json coreJson.CoreJson
var core_apiError coreJson.APIError

/*type APIError struct {
	StatusCode	int  `json:"statusCode"`
	Msg			string `json:"msg"`
}

func (e APIError) Error() string {
	return e.Msg
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:		err.Error(),
	}
}*/

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

func (h *HttpRouters) Add(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("Add")

	span := observ.Span(req.Context(), "handler.Add")
	defer span.End()

	onBoarding := model.Onboarding{}
	err := json.NewDecoder(req.Body).Decode(&onBoarding)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(erro.ErrUnmarshal, http.StatusBadRequest)
		return &core_apiError
    }
	defer req.Body.Close()

	res, err := h.workerService.Add(req.Context(), &onBoarding)
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