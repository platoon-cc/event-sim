package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/pkg/browser"
	"github.com/platoon-cc/platoon-cli/settings"
	"github.com/platoon-cc/platoon-cli/utils"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use: "auth",
}

var loginCmd = &cobra.Command{
	Use: "login",
	RunE: func(cmd *cobra.Command, args []string) error {
		return login()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(loginCmd)
}

func login() error {
	// check whether we are already logged in
	fmt.Println("Auth Login")

	callback, err := newServer()
	if err != nil {
		return err
	}

	route := "http://pl.localhost:9999/app/login"

	url := fmt.Sprintf("%s?redirect=true&hash=%s&port=%d", route, callback.Hash, callback.Port)
	if err := browser.OpenURL(url); err != nil {
		return err
	}

	fmt.Println("Opening your browser at:")
	fmt.Println(url)
	fmt.Println("Waiting for authentication...")
	jwt := callback.Result()

	settings.SetAuthToken(jwt)
	return nil
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
