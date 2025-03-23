package sessions

import (
	"session-auth/internal/utils"
	"sync"
)

type SessionStore struct {
	mu      sync.RWMutex
	session map[string]string
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		session: make(map[string]string),
	}
}

func (s *SessionStore) GetSession(id string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.session[id]
}

func (s *SessionStore) CreateSession(username string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := utils.GenerateRadomId(16)

	if err != nil {
		return "", err
	}

	s.session[id] = username

	return id, nil
}
