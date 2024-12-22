package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/subrotokumar/lbx/internal/config"
)

func Run() {
	fmt.Println("LBX...")
	config, err := config.GetConfigFromPath("./config.yml")
	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(0)
	}
	serverPool := NewFromConfig(*config)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: http.HandlerFunc(serverPool.lb),
	}
	go serverPool.HealthCheckCron()
	log.Printf("Load Balancer started at :%d\n", config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
