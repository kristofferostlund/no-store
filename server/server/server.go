package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/matryer/way"
	"github.com/sirupsen/logrus"
)

const ServerClosedErrorMessage = "http: Server closed"

type server struct {
	router     *way.Router
	httpServer *http.Server
}

func Serve(address string, port int) {
	addr := fmt.Sprintf("%s:%d", address, port)
	srv := server{
		router:     way.NewRouter(),
		httpServer: &http.Server{Addr: addr},
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGTERM)

	go serve(srv)

	<-stop

	closeServer(srv)
}

func serve(srv server) {
	logrus.Infof("Starting server on %s", srv.httpServer.Addr)

	srv.initRoutes()

	srv.httpServer.Handler = handlers.CompressHandler(
		handlers.LoggingHandler(os.Stdout, srv.router),
	)

	if err := srv.httpServer.ListenAndServe(); err != nil {
		if err.Error() != ServerClosedErrorMessage {
			logrus.Fatal(err)
		}
	}
}

func closeServer(srv server) {
	logrus.Info("Shutting down server...")

	defer os.Exit(0)
	defer logrus.Info("Server successfully shut down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.httpServer.Shutdown(ctx); err != nil {
		logrus.Fatalf("Failed to shut down server, I am literally dying: %v", err)
	}
}
