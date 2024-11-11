package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"

	"github.com/agusespa/ecom-be/customer/internal/httperrors"
	"github.com/agusespa/ecom-be/customer/internal/models"
	"github.com/agusespa/ecom-be/customer/internal/repository"
	logger "github.com/agusespa/flogg"
	"github.com/google/uuid"
)

type CustomerService interface {
	PostCustomer(body models.CustomerRequest) (int64, error)
	GetCustomerByID(id int64) (models.Customer, error)
}

type DefaultCustomerService struct {
	CustomerRepo repository.CustomerRepository
	AuthApyKey   string
	AuthDomain   string
	Logger       logger.Logger
}

func NewDefaultCustomerService(repo *repository.MySqlRepository, authApyKey, authDomain string, logger logger.Logger) *DefaultCustomerService {
	return &DefaultCustomerService{
		CustomerRepo: repo,
		AuthApyKey:   authApyKey,
		AuthDomain:   authDomain,
		Logger:       logger}
}

func (cs *DefaultCustomerService) GetCustomerByID(id int64) (models.Customer, error) {
	entity, err := cs.CustomerRepo.ReadCustomerByID(id)
	mappedCustomer := models.NewCustomer(
		entity.CustomerID,
		entity.CustomerUUID,
		entity.Email,
		entity.FirstName,
		entity.MiddleName,
		entity.LastName,
		entity.CreatedAt,
	)
	return mappedCustomer, err
}

func (cs *DefaultCustomerService) PostCustomer(body models.CustomerRequest) (int64, error) {
	if !IsValidEmail(body.Email) {
		err := errors.New("not a valid email address")
		err = httperrors.NewError(err, http.StatusBadRequest)
		cs.Logger.LogError(err)
		return 0, err
	}

	uuidStr := uuid.New().String()
	body.CustomerUUID = uuidStr

	tx, err := cs.CustomerRepo.StartTransaction()
	if err != nil {
		err := errors.New("failed to create transaction: " + err.Error())
		err = httperrors.NewError(err, http.StatusInternalServerError)
		cs.Logger.LogError(err)
		return 0, httperrors.NewError(err, http.StatusInternalServerError)
	}

	id, dbErr := cs.CustomerRepo.CreateCustomerWithTx(tx, body)
	if dbErr != nil {
		err := errors.New("failed to create customer: " + dbErr.Error())
		err = httperrors.NewError(err, http.StatusInternalServerError)
		cs.Logger.LogError(err)
		tx.Rollback()
		return 0, dbErr
	}

	authErr := cs.createAuthUser(body, uuidStr)
	if authErr != nil {
		err := errors.New("failed to create auth user: " + authErr.Error())
		err = httperrors.NewError(err, http.StatusInternalServerError)
		cs.Logger.LogError(err)
		tx.Rollback()
		return 0, authErr
	}

	if commitErr := tx.Commit(); commitErr != nil {
		authErr := cs.removeAuthUser(uuidStr)
		if authErr != nil {
		}
		err := errors.New("failed to commit customer: " + commitErr.Error())
		err = httperrors.NewError(err, http.StatusInternalServerError)
		cs.Logger.LogError(err)
		return 0, httperrors.NewError(commitErr, http.StatusInternalServerError)
	}

	return id, nil
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (cs *DefaultCustomerService) createAuthUser(body models.CustomerRequest, uuid string) error {
	customer := models.CustomerRequest{
		CustomerUUID: uuid,
		Email:        body.Email,
		Password:     body.Password,
		FirstName:    body.FirstName,
		MiddleName:   body.MiddleName,
		LastName:     body.LastName,
	}

	jsonData, err := json.Marshal(customer)
	if err != nil {
		return httperrors.NewError(err, http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", cs.AuthDomain+"/api/user", bytes.NewBuffer(jsonData))
	if err != nil {
		return httperrors.NewError(err, http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	cs.Logger.LogInfo(fmt.Sprintf("making auth request: %s %s", req.Method, req.URL))

	resp, err := client.Do(req)
	if err != nil {
		return httperrors.NewError(err, http.StatusServiceUnavailable)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return httperrors.NewError(errors.New("auth API error: unable to read error response"), resp.StatusCode)
		}
		errorMessage := string(bodyBytes)
		return httperrors.NewError(errors.New("auth API error: "+errorMessage), resp.StatusCode)
	}

	return nil
}

func (cs *DefaultCustomerService) removeAuthUser(uuid string) error {
	endpoint := fmt.Sprintf("%s/api/user?uuid=%s", cs.AuthDomain, uuid)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return httperrors.NewError(err, http.StatusInternalServerError)
	}

	req.Header.Set("Authentication", cs.AuthApyKey)
	client := &http.Client{}

	cs.Logger.LogInfo(fmt.Sprintf("making auth request: %s", req.URL))

	resp, err := client.Do(req)
	if err != nil {
		return httperrors.NewError(err, http.StatusServiceUnavailable)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return httperrors.NewError(errors.New("auth API error: unable to read error response"), resp.StatusCode)
		}
		errorMessage := string(bodyBytes)
		return httperrors.NewError(errors.New("auth API error: "+errorMessage), resp.StatusCode)
	}

	return nil
}
