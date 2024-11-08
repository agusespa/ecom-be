package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/agusespa/ecom-be/customer/internal/httperrors"
	"github.com/agusespa/ecom-be/customer/internal/models"
	"github.com/agusespa/ecom-be/customer/internal/payload"
	"github.com/agusespa/ecom-be/customer/internal/service"
	logger "github.com/agusespa/flogg"
)

type CustomerHandler interface {
	HandleCustomer(w http.ResponseWriter, r *http.Request)
	HandleCustomerRegister(w http.ResponseWriter, r *http.Request)
}

type DefaultCustomerHandler struct {
	CustomerService service.CustomerService
	Logger          logger.Logger
}

func NewDefaultCustomerHandler(customerService service.CustomerService, logger logger.Logger) *DefaultCustomerHandler {
	return &DefaultCustomerHandler{CustomerService: customerService, Logger: logger}
}

func (h *DefaultCustomerHandler) HandleCustomerRegister(w http.ResponseWriter, r *http.Request) {
	h.Logger.LogInfo(fmt.Sprintf("%s %v", r.Method, r.URL))

	if r.Method != http.MethodPost {
		h.Logger.LogError(fmt.Errorf("%s method not allowed for %v", r.Method, r.URL))
		err := httperrors.NewError(nil, http.StatusMethodNotAllowed)
		payload.WriteError(w, r, err)
		return
	}

	var userReq models.CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		err = httperrors.NewError(err, http.StatusBadRequest)
		h.Logger.LogError(err)
		payload.WriteError(w, r, err)
		return
	}

	if userReq.FirstName == "" || userReq.LastName == "" {
		err := errors.New("name not provided")
		err = httperrors.NewError(err, http.StatusBadRequest)
		h.Logger.LogError(err)
		payload.WriteError(w, r, err)
		return
	}

	if userReq.Email == "" || userReq.Password == "" {
		err := errors.New("missing credentials")
		err = httperrors.NewError(err, http.StatusUnauthorized)
		h.Logger.LogError(err)
		payload.WriteError(w, r, err)
		return
	}

	id, err := h.CustomerService.PostCustomer(userReq)
	if err != nil {
		payload.WriteError(w, r, err)
		return
	}

	res := models.RegistrationResponse{
		CustomerID: id,
	}

	payload.Write(w, r, res, nil)
}

func (h *DefaultCustomerHandler) HandleCustomer(w http.ResponseWriter, r *http.Request) {
	h.Logger.LogInfo(fmt.Sprintf("%s %v", r.Method, r.URL))

	if r.Method == http.MethodGet {
		uuid := r.Header.Get("X-User-UUID")
		if uuid == "" {
			err := httperrors.NewError(nil, http.StatusUnauthorized)
			payload.WriteError(w, r, err)
			return
		}

		customer, err := h.CustomerService.GetCustomerByUUID(uuid)
		if err != nil {
			payload.WriteError(w, r, err)
			return
		}
		payload.Write(w, r, customer, nil)
	}
}
