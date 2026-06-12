package utils

import (
	"sync"

	"github.com/joyboy1210/tex/internal/models"
)

var (
	sessionMut   sync.RWMutex
	userSessions = make(map[string]*models.UserSession)
)

func SetSessionCategory(phone string, categoryId uint) {
	sessionMut.Lock()
	defer sessionMut.Unlock()
	if _, exists := userSessions[phone]; !exists {
		userSessions[phone] = &models.UserSession{}
	}
	userSessions[phone].CategoryId = categoryId
}

func SetSessionProduct(phone string, productId uint) {
	sessionMut.Lock()
	defer sessionMut.Unlock()
	if _, exists := userSessions[phone]; !exists {
		userSessions[phone] = &models.UserSession{}
	}
	userSessions[phone].ProductId = productId
}

func GetSession(phone string) (*models.UserSession, bool) {
	sessionMut.RLock()
	defer sessionMut.RUnlock()
	session, exists := userSessions[phone]
	return session, exists
}

func ClearSession(phone string) {
	sessionMut.Lock()
	defer sessionMut.Unlock()
	delete(userSessions, phone)
}
