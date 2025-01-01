package server

import (
	"net"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

func (serverPool *ServerPool) HealthCheckCron() {
	t := time.NewTicker(time.Second * 20)
	defer t.Stop()
	for range t.C {
		log.Debugf("Starting health check...")
		serverPool.healthCheck()
		log.Debugf("Health check completed")
	}
}

func (serverPool *ServerPool) healthCheck() {
	for _, b := range serverPool.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Debugf("%s [%s]\n", b.URL, status)
	}
}

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Errorf("Site unreachable, error: %s", err.Error())
		return false
	}
	_ = conn.Close()
	return true
}

func (serverPool *ServerPool) markBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range serverPool.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}
