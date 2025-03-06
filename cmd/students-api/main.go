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
	"github.com/jalad-shrimali/students-api/internal/http/handlers/student"
	"github.com/jalad-shrimali/students-api/internal/storage/sqlite"
)

func main(){
	//load config using mustload
	cfg := config.MustLoad()

	//database setup
	storage, err := sqlite.New(cfg) //create a new sqlite database // you can also switch databases from here
	if err!=nil {
		log.Fatal(err)
	}
	slog.Info("database initialized", slog.String("env", cfg.Env))
	//setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage) ) //create a new student
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage) ) //get all students

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
	err = server.Shutdown(ctx)
	if err != nil{
		log.Fatal("failed to stop server ", slog.String("error", err.Error())) 
	}
	slog.Info("server stopped successfully")
}