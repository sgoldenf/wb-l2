package main

import (
	"context"
	"dev11/internal/controller/event"
	httphandler "dev11/internal/handler/http"
	"dev11/internal/repository/memory"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	repo := memory.New()
	ctrl := event.New(repo)
	h := httphandler.New(ctrl)
	m := http.NewServeMux()
	m.Handle("/create_event", h.Post(http.HandlerFunc(h.PostCreateEvent)))
	m.Handle("/update_event", h.Post(http.HandlerFunc(h.PostUpdateEvent)))
	m.Handle("/delete_event", h.Post(http.HandlerFunc(h.PostDeleteEvent)))
	m.Handle("/events_for_day", h.Get(http.HandlerFunc(h.GetEventsForDay)))
	m.Handle("/events_for_week", h.Get(http.HandlerFunc(h.GetEventsForWeek)))
	m.Handle("/events_for_month", h.Get(http.HandlerFunc(h.GetEventsForMonth)))
	s := http.Server{Handler: h.Log(m), Addr: ":8080"}
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM)
	<-sigTerm
	s.Shutdown(context.Background())
}
