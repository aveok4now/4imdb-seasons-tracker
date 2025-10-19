package main

import (
	"4imdb-seasons-tracker/internal/config"
	"4imdb-seasons-tracker/internal/handler"
	"4imdb-seasons-tracker/internal/repository"
	"4imdb-seasons-tracker/internal/scheduler"
	"4imdb-seasons-tracker/internal/scraper"
	"4imdb-seasons-tracker/internal/server"
	"4imdb-seasons-tracker/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := log.New(os.Stdout, "[IMDB-TRACKER] ", log.LstdFlags|log.Lshortfile)

	repo := repository.NewJSONRepository(cfg.Storage.FilePath, logger)
	if err := repo.Load(); err != nil {
		logger.Printf("Warning: failed to load data: %v", err)
	}

	scraperSvc := scraper.NewIMDBScraper(cfg.Scraper, logger)
	trackerSvc := service.NewTrackerService(repo, scraperSvc, logger)

	sched := scheduler.NewScheduler(cfg.Scheduler, trackerSvc, logger)
	sched.Start()
	defer sched.Stop()

	h := handler.NewHandler(trackerSvc, logger)
	srv := server.NewServer(cfg.Server, h, logger)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown error: %v", err)
	}

	logger.Println("Server stopped")
}
