package server

import (
	"context"
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
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Debugf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := GetRetryFromContext(request)
			if retries < 3 {
				<-time.After(10 * time.Millisecond)
				ctx := context.WithValue(request.Context(), Retry, retries+1)
				proxy.ServeHTTP(writer, request.WithContext(ctx))
				return
			}

			serverPool.markBackendStatus(serverUrl, false)

			attempts := GetAttemptsFromContext(request)
			log.Debugf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
			serverPool.lb(writer, request.WithContext(ctx))
		}
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

func (serverPool *ServerPool) GetNextPeer() *Backend {
	next := serverPool.NextIndex()
	l := len(serverPool.backends) + next

	for i := next; i < l; i++ {
		idx := i % len(serverPool.backends)
		if serverPool.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&serverPool.current, uint64(idx)) // mark the current one
			}
			return serverPool.backends[idx]
		}
	}

	return nil
}

func (serverPool *ServerPool) lb(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Debugf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := serverPool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
