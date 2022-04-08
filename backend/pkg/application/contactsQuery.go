package application

import (
	"GitHub/go-chat/backend/domain"
)

type ContactsQueryService interface {
	GetContacts() ([]*domain.UserDTO, error)
}

type contactsQueryService struct {
	users domain.UserQueryRepository
}

func NewContactsQueryService(users domain.UserQueryRepository) *contactsQueryService {
	return &contactsQueryService{
		users: users,
	}
}

func (s *contactsQueryService) GetContacts() ([]*domain.UserDTO, error) {
	return s.users.FindAll()
}
