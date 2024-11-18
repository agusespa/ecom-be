package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/agusespa/ecom-be/customer/cmd/middleware"
	"github.com/agusespa/ecom-be/customer/internal/database"
	"github.com/agusespa/ecom-be/customer/internal/handlers"
	"github.com/agusespa/ecom-be/customer/internal/helpers"
	"github.com/agusespa/ecom-be/customer/internal/models"
	"github.com/agusespa/ecom-be/customer/internal/repository"
	"github.com/agusespa/ecom-be/customer/internal/service"
	logger "github.com/agusespa/flogg"
)

func main() {
	var devFlag bool
	flag.BoolVar(&devFlag, "dev", false, "enable development mode")
	flag.Parse()

	logg := logger.NewLogger(devFlag, ".ecom_customer")

	dbUser, dbAddr, dbPassword, err := helpers.GetDatabaseVars()
	if err != nil {
		logg.LogFatal(fmt.Errorf("failed to read database env variables: %s", err.Error()))
	}

	authApiKey, authDomain, err := helpers.GetAppVars()
	if err != nil {
		logg.LogFatal(fmt.Errorf("failed to read app env variables: %s", err.Error()))
	}

	allowedIPs := os.Getenv("ECOM_ALLOWED_GATEWAY_IPS")
	if allowedIPs == "" {
		logg.LogFatal(fmt.Errorf("failed to read ALLOWED_GATEWAY_IPS env variable"))
	}
	allowedIPList := strings.Split(allowedIPs, ",")

	databaseConfig := models.Database{User: dbUser, Address: dbAddr, Password: dbPassword}
	db, err := database.ConnectDB(databaseConfig)
	if err != nil {
		logg.LogFatal(fmt.Errorf("failed to establish database connection: %s", err.Error()))
	}

	customerRepository := repository.NewMySqlRepository(db)
	customerService := service.NewDefaultCustomerService(customerRepository, authApiKey, authDomain, logg)
	customerHandler := handlers.NewDefaultCustomerHandler(customerService, logg)

	mux := http.NewServeMux()
	mux.HandleFunc("/customerapi/customer/", customerHandler.HandleCustomer)

	handler := middleware.ChainMiddleware(
		mux,
		func(h http.Handler) http.Handler { return middleware.CorsMiddleware(h, "todo") },
		func(h http.Handler) http.Handler { return middleware.GatewayMiddleware(h, allowedIPList, logg) },
	)

	port := os.Getenv("ECOM_CUSTOMER_PORT")
	if port == "" {
		port = "3002"
	}
	logg.LogInfo(fmt.Sprintf("Listening on port %v", port))
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		logg.LogFatal(fmt.Errorf("failed to start HTTP server: %s", err.Error()))
	}
}
