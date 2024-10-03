package service

import (
	"github.com/agusespa/ecom-be/customer/internal/models"
	"github.com/agusespa/ecom-be/customer/internal/repository"
)

type CustomerService struct {
	CustomerRepo *repository.CustomerRepository
}

func NewCustomerService(customerRepo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{CustomerRepo: customerRepo}
}

func (cs *CustomerService) GetCustomerByUUID(uuid string) (models.Customer, error) {
	entity, err := cs.CustomerRepo.QueryCustomerByUUID(uuid)
	mappedCustomer := models.NewCustomer(
		entity.ID,
		entity.UUID,
		entity.Email,
		entity.FirstName,
		entity.MiddleName,
		entity.LastName,
		entity.CreatedAt,
	)
	return mappedCustomer, err
}

func (cs *CustomerService) CreateCustomer(uuid string) (int32, error) {

}

func (cs *CustomerService) UpdateCustomer(uuid string) (int32, error) {
}
