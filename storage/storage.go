package storage

import (
	"github.com/pkg/errors"
	"github.com/sbani/amycr/config"
)

// NewManager is the factory method to create the storage manager from config
func NewManager(c *config.Config) (Manager, error) {
	switch c.Storage {
	case config.BoltDBKey:
		return NewBoltManager(), nil
	}

	return nil, errors.Errorf("Manager %s unknown", c.Storage)
}
