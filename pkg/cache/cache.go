package cache

import (
	"mstuca_schedule/internal/botErrors"
	"mstuca_schedule/internal/models"
	"sync"
)

type Cache interface {
	GetUser(id int64) (*models.User, error)
	SaveUserInfo(user *models.User)
	DeleteUser(id int64)
}

type cache struct {
	sync.RWMutex
	storage map[int64]*models.User
}

func New() Cache {
	mp := make(map[int64]*models.User)

	return &cache{
		storage: mp,
	}
}

func (c *cache) GetUser(id int64) (*models.User, error) {
	c.RLock()
	defer c.RUnlock()

	user, exists := c.storage[id]
	if !exists {
		return nil, botErrors.ErrUserNotPresents
	}

	return user, nil
}

func (c *cache) SaveUserInfo(user *models.User) {
	c.Lock()
	defer c.Unlock()

	c.storage[user.ID] = user
}

func (c *cache) DeleteUser(id int64) {
	delete(c.storage, id)
}
