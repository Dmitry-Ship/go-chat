package application

import (
	"GitHub/go-chat/backend/domain"
)

type ContactsQueryService interface {
	GetContacts() ([]*domain.User, error)
}

type contactsQueryService struct {
	users domain.UserRepository
}

func NewContactsQueryService(users domain.UserRepository) *contactsQueryService {
	return &contactsQueryService{
		users: users,
	}
}

func (s *contactsQueryService) GetContacts() ([]*domain.User, error) {
	return s.users.FindAll()
}
