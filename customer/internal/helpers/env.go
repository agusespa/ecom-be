package helpers

import (
	"errors"
	"os"
)

func GetAppVars() (string, string, error) {
	authApiKey := os.Getenv("ECOM_AUTH_API_KEY")
	if authApiKey == "" {
		return "", "", errors.New("failed to get AUTH_API_KEY variable")
	}
	authDomain := os.Getenv("ECOM_AUTH_DOMAIN")
	if authApiKey == "" {
		return "", "", errors.New("failed to get AUTH_DOMAIN variable")
	}
	return authApiKey, authDomain, nil
}

func GetDatabaseVars() (string, string, string, error) {
	dbUser := os.Getenv("ECOM_CUSTOMER_DB_USER")
	if dbUser == "" {
		return "", "", "", errors.New("failed to get DB_USER variable")
	}
	dbAddr := os.Getenv("ECOM_CUSTOMER_DB_ADDR")
	if dbAddr == "" {
		return "", "", "", errors.New("failed to get DB_ADDR variable")
	}
	dbPassword := os.Getenv("ECOM_CUSTOMER_DB_PASSWORD")
	if dbPassword == "" {
		return "", "", "", errors.New("failed to get DB_PASSWORD variable")
	}
	return dbUser, dbAddr, dbPassword, nil
}
