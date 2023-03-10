package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Edbeer/api-gateway/pkg/auth"
	"github.com/Edbeer/api-gateway/pkg/payment"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "github.com/Edbeer/api-gateway/docs"
	"github.com/swaggo/http-swagger"
)

// @title           API-Gateway
// @version         1.0
// @description     Simple payment system

// @securitydefinitions.apikey
// @in header
// @name Authorization
func main() {
	router := mux.NewRouter()

	// auth
	auth.RegisterAuthRoutes(router)
	// payment
	payment.RegisterPaymentRouter(router)
	
	// SWAGGER
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// in development mode, there will be "*"
	ch := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))
	server := &http.Server{
		Addr:         ":3000",
		Handler:      ch(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	
	log.Println("start")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<- quit

	ctx, shutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown()

	server.Shutdown(ctx)
}
