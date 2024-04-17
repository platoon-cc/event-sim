package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/pkg/browser"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/platoon-cc/platoon-cli/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "For connecting with the Platoon backend",
	}

	authCmd.PersistentFlags().BoolP("debug", "d", false, "connect to a local debug server")
	authCmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Connect to the backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			callback, err := newServer()
			if err != nil {
				return err
			}

			server := "https://platoon.cc"
			isDebug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				return err
			}
			if isDebug {
				server = "http://pl.localhost:9999"
			}

			route := server + "/app/login"

			url := fmt.Sprintf("%s?redirect=true&hash=%s&port=%d", route, callback.Hash, callback.Port)
			if err := browser.OpenURL(url); err != nil {
				return err
			}

			fmt.Println("Opening your browser at:")
			fmt.Println(url)
			fmt.Println("Waiting for authentication...")
			jwt := callback.Result()

			settings.SetAuth("server", server)
			settings.SetAuth("token", jwt)
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "logout",
		Short: "Disconnect from the backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			settings.ClearAuth("server")
			settings.ClearAuth("token")
			return nil
		},
	})

	rootCmd.AddCommand(authCmd)
}

type authServer struct {
	ch   chan string
	s    *http.Server
	Hash string
	Port int
}

func newServer() (authServer, error) {
	server := authServer{
		ch:   make(chan string, 1),
		Hash: utils.RandString(32),
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return server, fmt.Errorf("could not allocate port for http server: %w", err)
	}

	server.Port = listener.Addr().(*net.TCPAddr).Port
	server.s = &http.Server{Handler: server}

	go func() {
		server.s.Serve(listener)
	}()

	return server, nil
}

func (s authServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}

	header := w.Header()
	header.Set("Access-Control-Allow-Origin", origin)
	header.Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		header.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD")
		header.Set("Access-Control-Allow-Headers", "*")
		header.Set("Access-Control-Max-Age", "86400")
		return
	}

	q := r.URL.Query()
	if q.Get("hash") != s.Hash {
		w.WriteHeader(400)
		return
	}

	s.ch <- q.Get("token")

	w.WriteHeader(200)
	w.Write([]byte(`
<div class="text-center pb-4 text-lg">
	  You should have been successfully signed into your CLI. You can now close this window
</div>`))
}

func (s authServer) Result() string {
	result := <-s.ch
	_ = s.s.Shutdown(context.Background())
	return result
}
