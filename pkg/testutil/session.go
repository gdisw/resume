package testutil

import (
	"maps"
	"net/http"
	"sync"

	"github.com/gdisw/resume/pkg/http/session"
)

type SingleSessionStore struct {
	mutex  *sync.RWMutex
	values map[session.Key]any
}

func NewSingleSessionStore() *SingleSessionStore {
	return &SingleSessionStore{
		mutex:  &sync.RWMutex{},
		values: make(map[session.Key]any),
	}
}

func (s *SingleSessionStore) Put(_ *http.Request, _ http.ResponseWriter, k session.Key, v any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.values[k] = v

	return nil
}

func (s *SingleSessionStore) PutAll(_ *http.Request, _ http.ResponseWriter, vs map[session.Key]any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	maps.Copy(s.values, vs)

	return nil
}

func (s *SingleSessionStore) Get(_ *http.Request, k session.Key) (any, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.values[k], nil
}

func (s *SingleSessionStore) GetAll(_ *http.Request, ks []session.Key) (map[session.Key]any, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	out := make(map[session.Key]any)
	maps.Copy(out, s.values)

	return out, nil
}

func (s *SingleSessionStore) Delete(_ *http.Request, _ http.ResponseWriter) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.values = make(map[session.Key]any)

	return nil
}
