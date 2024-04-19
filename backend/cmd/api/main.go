package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/api/router"
	"backend/config"
)

// @title           Users API
// @version         1.0
// @description		RESTful API enabling CRUD operations (Create, Read, Update, Delete) for user management in web application.

// @license.name  MIT License
// @license.url   https://github.com/HealthHub-Management-System/UserService/blob/master/LICENSE

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	c := config.New()
	r := router.New()

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	log.Println("Starting server " + s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
}
