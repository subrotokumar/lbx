package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/subrotokumar/lbx/internal/server"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	server.Run()
}
