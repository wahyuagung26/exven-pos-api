package persistence

import (
	"context"
	"fmt"
	"sync"

	"github.com/exven/pos-system/modules/auth/domain"
)

// InMemorySessionRepository is a simple in-memory session repository for testing
type SessionRepository struct {
	sessions map[string]*domain.Session
	mu       sync.RWMutex
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[string]*domain.Session),
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID] = session
	return nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, sessionID)
	return nil
}

func (r *SessionRepository) FindByID(ctx context.Context, sessionID string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if session, exists := r.sessions[sessionID]; exists {
		return session, nil
	}

	return nil, fmt.Errorf("session not found")
}

func (r *SessionRepository) FindByUserID(ctx context.Context, userID uint64) ([]*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var sessions []*domain.Session
	for _, session := range r.sessions {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, session := range r.sessions {
		if session.UserID == userID {
			delete(r.sessions, id)
		}
	}

	return nil
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	// Not implemented for simple testing
	return nil
}
