package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"loan-service/controllers"
	"loan-service/middleware"
)

// SetupRouter initializes all routes with middleware and controller dependencies
func SetupRouter(controller *controllers.LoanController) *mux.Router {
	router := mux.NewRouter()

	// Global middleware
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.LoggingMiddleware)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
		})
	}).Methods("GET")

	// Loan routes (protected)clear
	loanRouter := router.PathPrefix("/loans").Subrouter()
	loanRouter.Use(middleware.AuthMiddleware)
	loanRouter.HandleFunc("", controller.CreateLoan).Methods("POST")
	loanRouter.HandleFunc("/{id}/approve", controller.ApproveLoan).Methods("PUT")
	loanRouter.HandleFunc("/{id}/invest", controller.InvestLoan).Methods("POST")
	loanRouter.HandleFunc("/{id}/disburse", controller.DisburseLoan).Methods("PUT")

	return router
}