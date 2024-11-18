package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/agusespa/ecom-be/customer/internal/helpers"
	"github.com/agusespa/ecom-be/customer/internal/httperrors"
	"github.com/agusespa/ecom-be/customer/internal/models"
	"github.com/agusespa/ecom-be/customer/internal/payload"
	"github.com/agusespa/ecom-be/customer/internal/service"
	logger "github.com/agusespa/flogg"
)

type CustomerHandler interface {
	HandleCustomer(w http.ResponseWriter, r *http.Request)
}

type DefaultCustomerHandler struct {
	CustomerService service.CustomerService
	Logger          logger.Logger
}

func NewDefaultCustomerHandler(customerService service.CustomerService, logger logger.Logger) *DefaultCustomerHandler {
	return &DefaultCustomerHandler{CustomerService: customerService, Logger: logger}
}

func (h *DefaultCustomerHandler) HandleCustomer(w http.ResponseWriter, r *http.Request) {
	h.Logger.LogInfo(fmt.Sprintf("%s %v", r.Method, r.URL))

	path := strings.TrimPrefix(r.URL.Path, "/customerapi/customer")
	parts := strings.Split(path, "/")
	customerID, err := helpers.StringToInt64(parts[0])
	if err != nil {
		err = httperrors.NewError(err, http.StatusInternalServerError)
		h.Logger.LogError(err)
		payload.WriteError(w, r, err)
		return
	}

	switch {
	case len(parts) < 2:
		// customer/{id}
		h.handleCustomerData(w, r, customerID)
	default:
		http.NotFound(w, r)
	}
}

func (h *DefaultCustomerHandler) handleCustomerData(w http.ResponseWriter, r *http.Request, customerID int64) {
	if r.Method == http.MethodPost {
		h.handlePostCustomer(w, r)
		return
	} else {
		if customerID == 0 {
			err := errors.New("User ID required")
			err = httperrors.NewError(err, http.StatusBadRequest)
			h.Logger.LogError(err)
			payload.WriteError(w, r, err)
			return
		}

		reqUUID := r.Header.Get("X-Auth-Uuid")
		err := h.verifyRequestPermission(customerID, reqUUID)
		if err != nil {
			h.Logger.LogError(err)
			payload.WriteError(w, r, err)
			return
		}

		if r.Method == http.MethodGet {
			h.handleGetCustomer(w, r, customerID)
			return
		} else if r.Method == http.MethodDelete {
			// h.HandleDeleteCustomer(w, r)
			return
		}
	}

	h.Logger.LogError(fmt.Errorf("%s method not allowed for %v", r.Method, r.URL))
	err := httperrors.NewError(nil, http.StatusMethodNotAllowed)
	payload.WriteError(w, r, err)
}

func (h *DefaultCustomerHandler) verifyRequestPermission(customerID int64, userUUID string) error {
	uuid, err := h.CustomerService.GetCustomerUUID(customerID)
	if err != nil {
		return err
	}

	if uuid != userUUID {
		err := errors.New("Customer doesn't have permission to access")
		err = httperrors.NewError(err, http.StatusUnauthorized)
		return err
	}

	return nil
}

func (h *DefaultCustomerHandler) handleGetCustomer(w http.ResponseWriter, r *http.Request, customerID int64) {
	customer, err := h.CustomerService.GetCustomerByID(customerID)
	if err != nil {
		h.Logger.LogError(err)
		payload.WriteError(w, r, err)
		return
	}

	payload.Write(w, r, customer, nil)
}

func (h *DefaultCustomerHandler) handlePostCustomer(w http.ResponseWriter, r *http.Request) {
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

	if !helpers.IsValidEmail(userReq.Email) {
		err := errors.New("not a valid email address")
		err = httperrors.NewError(err, http.StatusBadRequest)
		h.Logger.LogError(err)
		payload.WriteError(w, r, err)
		return
	}

	if !helpers.IsValidPassword(userReq.Password) {
		err := errors.New("password doesn't meet minimum criteria")
		err = httperrors.NewError(err, http.StatusBadRequest)
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
