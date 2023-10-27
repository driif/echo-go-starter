package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/driif/echo-go-starter/internal/api/router"
	"github.com/driif/echo-go-starter/internal/server"
	"github.com/driif/echo-go-starter/internal/server/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the server",
	Long: `Starts the stateless RESTful JSON server

Requires configuration through ENV and
and a fully migrated PostgreSQL database.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

// init adds the server command to the root command.
func init() {
	rootCmd.AddCommand(serverCmd)
}

// runServer starts the server.
func runServer() {
	config := config.DefaultServiceConfigFromEnv()

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(config.Logger.Level)
	if config.Logger.PrettyPrintConsole {
		log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "15:04:05"
		}))
	}

	fmt.Println("Starting server...")
	s := server.New(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := s.InitDB(ctx); err != nil {
		cancel()
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	cancel()

	if err := s.Initialize(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize server")
		os.Exit(1)
	}

	router.InitGroups(s)
	router.AttachRoutes(s)

	go func() {
		if err := s.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info().Msg("Server closed")
			} else {
				log.Fatal().Err(err).Msg("Failed to start server")
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("Failed to gracefully shut down server")
	}
}
