package scheduler

import (
	"4imdb-seasons-tracker/internal/config"
	"4imdb-seasons-tracker/internal/service"
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron    *cron.Cron
	service *service.TrackerService
	logger  *log.Logger
	config  config.SchedulerConfig
}

func NewScheduler(config config.SchedulerConfig, service *service.TrackerService, logger *log.Logger) *Scheduler {
	c := cron.New()

	return &Scheduler{
		cron:    c,
		service: service,
		logger:  logger,
		config:  config,
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc(s.config.CronExpression, func() {
		s.logger.Println("Starting scheduled check...")
		if err := s.service.CheckAll(); err != nil {
			s.logger.Printf("Scheduled check error: %v", err)
		}
		s.logger.Println("Scheduled check completed")
	})

	s.cron.Start()
	s.logger.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.logger.Println("Scheduler stopped")
}
