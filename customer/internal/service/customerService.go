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
	GetCustomerByUUID(uuid string) (models.Customer, error)
}

type DefaultCustomerService struct {
	CustomerRepo repository.CustomerRepository
	Logger       logger.Logger
}

func NewDefaultCustomerService(repo *repository.MySqlRepository, logger logger.Logger) *DefaultCustomerService {
	return &DefaultCustomerService{
		CustomerRepo: repo,
		Logger:       logger}
}

func (cs *DefaultCustomerService) GetCustomerByUUID(uuid string) (models.Customer, error) {
	entity, err := cs.CustomerRepo.ReadCustomerByUUID(uuid)
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

	tx, err := cs.CustomerRepo.StartTransaction()
	if err != nil {
		err := errors.New("failed to create transaction: " + err.Error())
		err = httperrors.NewError(err, http.StatusInternalServerError)
		cs.Logger.LogError(err)
		return 0, httperrors.NewError(err, http.StatusInternalServerError)
	}

	id, dbErr := cs.CustomerRepo.CreateCustomerWithTx(tx, body, uuidStr)
	if dbErr != nil {
		err := errors.New("failed to create customer: " + dbErr.Error())
		err = httperrors.NewError(err, http.StatusInternalServerError)
		cs.Logger.LogError(err)
		tx.Rollback()
		return 0, dbErr
	}

	authErr := CreateUser(body, uuidStr)
	if authErr != nil {
		err := errors.New("failed to create customer: " + authErr.Error())
		err = httperrors.NewError(err, http.StatusInternalServerError)
		cs.Logger.LogError(err)
		tx.Rollback()
		return 0, authErr
	}

	if commitErr := tx.Commit(); commitErr != nil {
		authErr := RemoveUser(uuidStr)
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

func CreateUser(body models.CustomerRequest, uuid string) error {
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

	req, err := http.NewRequest("POST", "http://192.168.64.7/a3n/api/register", bytes.NewBuffer(jsonData))
	if err != nil {
		return httperrors.NewError(err, http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
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

func RemoveUser(uuid string) error {
	endpoint := fmt.Sprintf("http://192.168.64.7/a3n/api/user?uuid=%s", uuid)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return httperrors.NewError(err, http.StatusInternalServerError)
	}
	//TODO add the api key
	client := &http.Client{}
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
