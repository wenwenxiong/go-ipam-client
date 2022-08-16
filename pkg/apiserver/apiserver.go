package apiserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"wenwenxiong/go-ipam-client/pkg/client/goipam"
)

type APIServer struct {
	ListenPort int
	Endpoint   string

	// goipam client
	GoipamClient goipam.Client
	Server       *http.Server
}

func (s *APIServer) Run() (err error) {

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	go func() {
		// Block until we receive our signal.
		<-c
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		s.Server.Shutdown(ctx)
	}()

	log.Println("Start listening on %s", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		log.Println(err)
	}

	return err
}
