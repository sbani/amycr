package storage

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sbani/gcr/pkg"
)

type MemoryManager struct {
	contentTypes map[string]string
	records      map[string]string

	sync.RWMutex
}

func newMemoryManager() *MemoryManager {
	return &MemoryManager{}
}

func (m *MemoryManager) CreateRecord(content string) error {
	m.Lock()
	defer m.Unlock()

	m.records[m.newId()] = content
	return nil
}

func (m *MemoryManager) GetRecord(id string) (string, error) {
	m.Lock()
	defer m.Unlock()

	r, ok := m.records[id]
	if !ok {
		return "", errors.Wrap(pkg.ErrNotFound, "MemoryManager")
	}
	return r, nil
}

func (m *MemoryManager) DeleteRecord(id string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.records, id)
	return nil
}

func (m *MemoryManager) CreateContentType(id, content string) error {
	m.Lock()
	defer m.Unlock()

	m.contentTypes[id] = content
	return nil
}

func (m *MemoryManager) GetContentType(id string) (string, error) {
	m.Lock()
	defer m.Unlock()

	r, ok := m.contentTypes[id]
	if !ok {
		return "", errors.Wrap(pkg.ErrNotFound, "MemoryManager")
	}
	return r, nil
}

func (m *MemoryManager) DeleteContentType(id string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.contentTypes, id)
	return nil
}

// newId creates a new unqique
func (m *MemoryManager) newId() string {
	m.Lock()
	defer m.Unlock()

	return uuid.NewV1().String()
}
