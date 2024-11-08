package helpers

import (
	"errors"
	"os"
)

func GetApiKeyVars() (string, error) {
	encryptionKey := os.Getenv("ECOM_CUSTOMER_ENCRYPTION_KEY")
	if encryptionKey == "" {
		return "", errors.New("failed to get ENCRYPTION_KEY variable")
	}
	return encryptionKey, nil
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
