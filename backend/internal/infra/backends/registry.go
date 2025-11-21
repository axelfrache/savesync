package backends

import (
	"fmt"

	"github.com/axelfrache/savesync/internal/domain"
	"github.com/axelfrache/savesync/internal/infra/backends/local"
	"github.com/axelfrache/savesync/internal/infra/backends/s3"
	"github.com/axelfrache/savesync/internal/infra/backends/sftp"
)

type Registry struct {
	backends map[string]func() domain.Backend
}

func NewRegistry() *Registry {
	r := &Registry{
		backends: make(map[string]func() domain.Backend),
	}

	r.Register("local", func() domain.Backend { return &local.Backend{} })
	r.Register("s3", func() domain.Backend { return &s3.Backend{} })
	r.Register("sftp", func() domain.Backend { return &sftp.Backend{} })

	return r
}

func (r *Registry) Register(backendType string, factory func() domain.Backend) {
	r.backends[backendType] = factory
}

func (r *Registry) Create(backendType string, config map[string]string) (domain.Backend, error) {
	factory, exists := r.backends[backendType]
	if !exists {
		return nil, fmt.Errorf("unknown backend type: %s", backendType)
	}

	backend := factory()
	if err := backend.Init(config); err != nil {
		return nil, fmt.Errorf("failed to initialize backend: %w", err)
	}

	return backend, nil
}

func (r *Registry) IsSupported(backendType string) bool {
	_, exists := r.backends[backendType]
	return exists
}

func (r *Registry) SupportedTypes() []string {
	types := make([]string, 0, len(r.backends))
	for t := range r.backends {
		types = append(types, t)
	}
	return types
}
