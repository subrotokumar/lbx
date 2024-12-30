package server

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/subrotokumar/lbx/internal/config"
)

func Run() {
	log.Infoln("Starting Load Balancer")
	log.Infoln("Reading configuration")

	config, err := config.GetConfigFromPath("./config.yml")
	if err != nil {
		log.Errorf("%s", err.Error())
		os.Exit(0)
	}

	serverPool := NewFromConfig(*config)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: http.HandlerFunc(serverPool.lb),
	}

	go serverPool.HealthCheckCron()

	log.Infof("Load Balancer started at :%d\n", config.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
