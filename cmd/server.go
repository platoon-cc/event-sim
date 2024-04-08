package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/platoon-cc/platoon-cli/internal/server"
	"github.com/spf13/cobra"
)

func init() {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "run a test server for prototyping event ingestion",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := server.New()

			stopChan := make(chan os.Signal)
			errChan := make(chan error)

			// Setup the graceful shutdown handler (traps SIGINT and SIGTERM)
			go func() {
				signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

				<-stopChan

				timer, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := s.Stop(timer); err != nil {
					errChan <- err
					return
				}

				errChan <- nil
			}()

			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				return err
			}
			// Start the server
			if err := s.Start(port); err != http.ErrServerClosed {
				return err
			}

			return <-errChan
		},
	}

	serverCmd.Flags().Int("port", 9998, "port to listen on")

	rootCmd.AddCommand(serverCmd)
}
