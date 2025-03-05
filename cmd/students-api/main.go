package main

import (
	"context"
	// "fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jalad-shrimali/students-api/internal/config"
)

func main(){
	//load config using mustload
	cfg := config.MustLoad()

	//database setup

	//setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Welcome to students api!"))
	})

	//start server
	// before starting the server, we will create a channel to listen for the interrupt signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //listen for interrupt signal and terminate signal
	server := http.Server{
		Addr: cfg.HTTPServer.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String("address", cfg.HTTPServer.Addr))
	// fmt.Printf("server started at %s", cfg.HTTPServer.Addr)
	
	go func ()  {
		err := server.ListenAndServe() //start server
	if err != nil{
		log.Fatal("failed to start server")
	}
	} ()

	<-done
	slog.Info("server stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil{
		log.Fatal("failed to stop server ", slog.String("error", err.Error())) 
	}
	slog.Info("server stopped successfully")
}