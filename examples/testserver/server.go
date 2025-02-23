package testserver

import (
	"log/slog"
	"net/http/httptest"
	"sync"
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
	m.mu.Lock()

	testsrv := httptest.NewServer(NewMux(config, m.repository))
	config.ID = m.generateServerID()
	config.URL = testsrv.URL
	m.running = append(m.running, *config)

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

func (b *ServerBuilder) Build() *httptest.Server {
	return b.manager.newServer(b.config)
}
