package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"backend/api/router"
	"backend/config"
	"backend/utils/logger"
	validatorUtil "backend/utils/validator"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

// @title           Users API
// @version         1.0
// @description		RESTful API enabling CRUD operations (Create, Read, Update, Delete) for user management in web application.

// @license.name  MIT License
// @license.url   https://github.com/HealthHub-Management-System/UserService/blob/master/LICENSE

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	c := config.New()
	l := logger.New(c.Server.Debug)
	v := validatorUtil.New()
	store := sessions.NewCookieStore([]byte("chuj"))

	var logLevel gormlogger.LogLevel
	if c.Database.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	dbString := fmt.Sprintf(fmtDBString, c.Database.Host, c.Database.Username, c.Database.Password, c.Database.Name, c.Database.Port)
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
	if err != nil {
		log.Fatal("DB connection start failure")
		return
	}

	r := router.New(l, db, v, store)

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
