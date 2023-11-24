package main

import (
	"context"
	"fmt"
	"job-portal-api/config"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/caching"

	"job-portal-api/internal/database"
	"job-portal-api/internal/handlers"
	"job-portal-api/internal/repository"
	"job-portal-api/internal/services"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	err := startApp()
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("Hello this is our app")
}

// =========================================================================
// Initializing  Authentication Support
func startApp() error {
	cfg := config.GetConfig()
	log.Info().Msg("Config done")

	log.Info().Msg("main : Started : Initializing authentication support")
	// privatePEM, err := os.ReadFile(`private.pem`)
	// if err != nil {
	// 	return fmt.Errorf("reading auth private key %w", err)
	// }
	privatePEM:=[]byte(cfg.AuthConfig.PrivateKey)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("parsing auth private key %w", err)
	}

	// publicPEM, err := os.ReadFile(`pubkey.pem`)
	// if err != nil {
	// 	return fmt.Errorf("reading auth public key %w", err)
	// }
    publicPEM:=[]byte(cfg.AuthConfig.PublicKey)
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return fmt.Errorf("parsing auth public key %w", err)
	}

	a, err := auth.NewAuth(privateKey, publicKey)
	if err != nil {
		return fmt.Errorf("constructing auth %w", err)
	}

	// =========================================================================
	// Starting the Database
	log.Info().Msg("main : Started : Initializing database support")
	db, err := database.OpenConnection()
	if err != nil {
		return fmt.Errorf("connecting to db %w", err)
	}
	pg, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = pg.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("database is not connected: %w ", err)
	}
	// redis database connection
	rdb := database.RedisConnection()
	redisLayer, err := caching.NewRedis(rdb)
	if err != nil {
		return fmt.Errorf("redis db is not connected: %w ", err)
	}
	// =========================================================================
	//Initialize Conn layer support

	r, err := repository.NewRepository(db)
	if err != nil {
		return err
	}

	ms, err := services.NewService(r, a, redisLayer)
	if err != nil {
		return err
	}

	// =========================================================================
	// Initialize http service
	api := http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppConfig.Port),
		ReadTimeout:  time.Duration(cfg.AppConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.AppConfig.WriteTimeout)*time.Second,
		IdleTimeout:  time.Duration(cfg.AppConfig.IdleTimeout)*time.Second,
		Handler:      handlers.API(a, ms),
	}

	//Server termination
	serverErrors := make(chan error, 1)
	go func() {
		log.Info().Str("port", api.Addr).Msg("main: API listening")
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error %w", err)
	case sig := <-shutdown:
		log.Info().Msgf("main: Start shutdown %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			err = api.Close()
			return fmt.Errorf("could not stop server gracefully %w", err)
		}

	}
	return nil

}
