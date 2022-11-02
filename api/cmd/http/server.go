package main

import (
	"net/http"
	"time"

	"github.com/IndominusByte/catalyst-be/api/internal/config"
)

func startServer(mux *http.ServeMux, cfg *config.Config) error {
	readTimeout, err := time.ParseDuration(cfg.Server.HTTP.ReadTimeout)
	if err != nil {
		return err
	}
	writeTimeout, err := time.ParseDuration(cfg.Server.HTTP.WriteTimeout)
	if err != nil {
		return err
	}

	srv := http.Server{
		Addr:         cfg.Server.HTTP.Address,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      mux,
	}

	return srv.ListenAndServe()
}
