package main

import (
	"net/http"

	"github.com/lonelyshoeh/api/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func listenPrometheus(cfg config.APIConfig) {
	router := chi.NewRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    cfg.PrometheusBindAddress,
	}

	router.Handle("/metrics", promhttp.Handler())

	go srv.ListenAndServe()
}
