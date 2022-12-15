package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/AllPaste/web-bbf/config"
	"github.com/AllPaste/web-bbf/internal/routers"
)

func Run(conf *config.Config) {
	srv := &http.Server{
		Addr: strings.Join(
			[]string{conf.Server.Address, strconv.Itoa(conf.Server.Port)},
			":",
		),
		Handler:        routers.InitRouter(),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("server is runing port: %s\n", strings.Join(
			[]string{conf.Server.Address, strconv.Itoa(conf.Server.Port)},
			":",
		))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGSEGV,
		syscall.SIGABRT,
		syscall.SIGILL,
		syscall.SIGFPE,
		os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
