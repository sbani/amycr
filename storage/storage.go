package storage

import (
	"github.com/pkg/errors"
	"github.com/sbani/gcr/config"
)

func NewManager(c *config.Config) (Manager, error) {
	switch c.Storage {
	case "boltdb":
		return newBoltManager(), nil
	}

	return nil, errors.Errorf("Manager %s unknown", c.Storage)
}
