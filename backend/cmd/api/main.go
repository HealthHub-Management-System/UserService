package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wader/gormstore/v2" // Add this import

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

// @host      127.0.0.1:8080
// @BasePath  /api/v1
func main() {
	c := config.New()
	l := logger.New(c.Server.Debug)
	v := validatorUtil.New()

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

	store := gormstore.New(db, []byte(c.Server.Secret))
	store.SessionOpts.SameSite = http.SameSiteNoneMode
	store.SessionOpts.Secure = true

	r := router.New(l, db, v, store, c.Server.APIKey)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go store.PeriodicCleanup(1*time.Hour, closed)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		l.Info().Msgf("Shutting down server %v", s.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), c.Server.TimeoutIdle)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			l.Error().Err(err).Msg("Server shutdown failure")
		}

		sqlDB, err := db.DB()
		if err == nil {
			if err = sqlDB.Close(); err != nil {
				l.Error().Err(err).Msg("DB connection closing failure")
			}
		}

		close(closed)
	}()

	l.Info().Msgf("Starting server %v", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	l.Info().Msgf("Server shutdown successfully")
}
