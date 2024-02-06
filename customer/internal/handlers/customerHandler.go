package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/agusespa/ecom-be/customer/internal/errors"
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
	fmt.Printf("CustomerHandler.HandleCustomerByUUID %v\n", r)

	if r.Method != http.MethodGet {
		err := errors.NewError(nil, http.StatusMethodNotAllowed)
		payload.WriteError(w, r, err)
		return
	}

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

	type CustomClaims struct {
		jwt.StandardClaims
		UserUUID string `json:"UserUUID"`
	}

	claims := &CustomClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(bearerToken, claims)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}

	customer, err := h.CustomerService.GetCustomerByUUID(claims.UserUUID)
	if err != nil {
		payload.WriteError(w, r, err)
		return
	}

	fmt.Printf("Customer: %v\n", customer)

	payload.Write(w, r, customer)
}
