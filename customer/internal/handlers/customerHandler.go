package handlers

import (
	"net/http"
	"strings"

	"github.com/agusespa/ecom-be/customer/internal/errors"
	"github.com/agusespa/ecom-be/customer/internal/models"
	"github.com/agusespa/ecom-be/customer/internal/payload"
	"github.com/agusespa/ecom-be/customer/internal/service"
	"github.com/golang-jwt/jwt"
)

type CustomerHandler struct {
	CustomerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{CustomerService: customerService}
}

func (h *CustomerHandler) HandleCustomerByUUID(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		err := errors.NewError(nil, http.StatusUnauthorized)
		payload.WriteError(w, r, err)
		return
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		err := errors.NewError(nil, http.StatusUnauthorized)
		payload.WriteError(w, r, err)
		return
	}

	bearerToken := authParts[1]

	claims := &models.CustomClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(bearerToken, claims)
	if err != nil {
		err := errors.NewError(nil, http.StatusUnauthorized)
		payload.WriteError(w, r, err)
		return
	}

	if r.Method == http.MethodGet {
		customer, err := h.CustomerService.GetCustomerByUUID(claims.User.UserUUID)
		if err != nil {
			payload.WriteError(w, r, err)
			return
		}
		payload.Write(w, r, customer)
	}

	if r.Method == http.MethodPost {
		customer, err := h.CustomerService.CreateCustomer(claims.User.UserUUID)
		if err != nil {
			payload.WriteError(w, r, err)
			return
		}
		payload.Write(w, r, customer)
	}

	if r.Method == http.MethodPut {
		customer, err := h.CustomerService.UpdateCustomer(claims.User.UserUUID)
		if err != nil {
			payload.WriteError(w, r, err)
			return
		}
		payload.Write(w, r, customer)
	}

	err = errors.NewError(nil, http.StatusMethodNotAllowed)
	payload.WriteError(w, r, err)
}
