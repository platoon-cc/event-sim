package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/platoon-cc/platoon-cli/internal/model"
	"github.com/platoon-cc/platoon-cli/internal/processor"
)

type Server struct {
	srv       *http.Server
	processor *processor.Processor
}

func New() *Server {
	server := &Server{}

	return server
}

func errorWrapper(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			log.Printf("Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error"))
		}
	}
}

func (s *Server) ingestHandler(w http.ResponseWriter, r *http.Request) error {
	reader := io.LimitReader(r.Body, maxIngestPayloadSize)
	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if b[0] == '[' {
		p := []model.Event{}
		err = json.Unmarshal(b, &p)
		if err != nil {
			if err != io.EOF {
				return err
			}
		}
		for _, e := range p {
			if err := s.processor.IngestEvent(e); err != nil {
				return err
			}
		}
	} else {
		e := model.Event{}
		err = json.Unmarshal(b, &e)
		if err != nil {
			if err != io.EOF {
				return err
			}
		}
		if err := s.processor.IngestEvent(e); err != nil {
			return err
		}
	}

	return nil
}

var maxIngestPayloadSize int64 = 1024 * 512

func (s *Server) Start(port int) error {
	addr := fmt.Sprintf("0.0.0.0:%d", port)

	processor, err := processor.New("local")
	if err != nil {
		return err
	}
	s.processor = processor

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/ingest", errorWrapper(s.ingestHandler))

	s.srv = &http.Server{
		Addr:              addr,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		Handler:           mux,
	}

	log.Printf("Starting server at %s\n", addr)
	s.srv.ListenAndServe()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Printf("Stopping server\n")
	s.processor.Close()
	s.srv.Shutdown(ctx)
	return nil
}
