package repository

import (
	"sync"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type UserRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
}

// InMemoryUserRepo stores users in memory
type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*entity.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]*entity.User),
	}
}

func (r *InMemoryUserRepo) AddUser(u *entity.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.Email] = u
}

func (r *InMemoryUserRepo) GetUserByEmail(email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[email]
	if !ok {
		return nil, entity.ErrorNotFound("user not found")
	}
	return u, nil
}
