package testserver

import (
	"log/slog"
	"net/http/httptest"
	"sync"
	"time"
)

type (
	Manager struct {
		running    []Server
		repository *Provider
		mu         sync.Mutex
	}

	ServerBuilder struct {
		manager *Manager
		config  *Server
	}
)

func NewManager() *Manager {
	state := NewState()

	return &Manager{
		running:    []Server{},
		repository: NewProvider(state),
		mu:         sync.Mutex{},
	}
}

func (m *Manager) generateServerID() int {
	return len(m.running) + 1
}

func (m *Manager) newServer(config *Server) *httptest.Server {
	testsrv := httptest.NewServer(NewMux(config, m.repository))

	config.ID = m.generateServerID()
	config.URL = testsrv.URL

	m.mu.Lock()
	m.running = append(m.running, *config)
	m.mu.Unlock()

	slog.Info("test server running!", "config", config)

	return testsrv
}

func (m *Manager) NewServer() *httptest.Server {
	return m.newServer(&Server{})
}

func (m *Manager) NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		manager: m,
		config:  &Server{},
	}
}

func (b *ServerBuilder) EnableBusy() *ServerBuilder {
	b.config.EnableBusy = true
	return b
}

func (b *ServerBuilder) EnableHeaderDebug() *ServerBuilder {
	b.config.EnableHeaderDebug = true
	return b
}

func (b *ServerBuilder) SleepFor(d time.Duration, jitter float64) *ServerBuilder {
	if jitter == 0 {
		jitter = 2.0
	}

	b.config.Interval = d
	b.config.Jitter = jitter

	return b
}

func (b *ServerBuilder) Build() *httptest.Server {
	return b.manager.newServer(b.config)
}
