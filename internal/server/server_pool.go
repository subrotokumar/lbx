package server

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/subrotokumar/lbx/internal/config"
)

type ServerPool struct {
	backends []*Backend
	current  uint64
}

func New() *ServerPool {
	return &ServerPool{
		backends: make([]*Backend, 0),
		current:  0,
	}
}
func NewFromConfig(cfg config.Config) *ServerPool {
	serverPool := New()
	for _, serverConfig := range cfg.Servers {
		serverUrl, err := url.Parse(serverConfig.URL)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		/*
			proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
				log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
				retries := GetRetryFromContext(request)
				if retries < 3 {
					select {
					case <-time.After(10 * time.Millisecond):
						ctx := context.WithValue(request.Context(), Retry, retries+1)
						proxy.ServeHTTP(writer, request.WithContext(ctx))
					}
					return
				}

				// after 3 retries, mark this backend as down
				serverPool.MarkBackendStatus(serverUrl, false)

				// if the same request routing for few attempts with different backends, increase the count
				attempts := GetAttemptsFromContext(request)
				log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
				ctx := context.WithValue(request.Context(), Attempts, attempts+1)
				lb(writer, request.WithContext(ctx))
			}
		*/
		serverPool.AddBackend(&Backend{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Debugf("Configured server: %s\n", serverUrl)
	}
	return serverPool
}

func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()
	l := len(s.backends) + next

	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx)) // mark the current one
			}
			return s.backends[idx]
		}
	}

	return nil
}

func (s *ServerPool) lb(w http.ResponseWriter, r *http.Request) {
	peer := s.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	_ = conn.Close()
	return true
}

func (s *ServerPool) healthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

func (s *ServerPool) HealthCheckCron() {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			s.healthCheck()
			log.Println("Health check completed")
		}
	}
}
