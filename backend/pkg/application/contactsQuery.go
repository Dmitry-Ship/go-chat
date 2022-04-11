package application

import (
	"GitHub/go-chat/backend/pkg/readModel"
)

type ContactsQueryService interface {
	GetContacts() ([]*readModel.UserDTO, error)
}

type contactsQueryService struct {
	users readModel.UserQueryRepository
}

func NewContactsQueryService(users readModel.UserQueryRepository) *contactsQueryService {
	return &contactsQueryService{
		users: users,
	}
}

func (s *contactsQueryService) GetContacts() ([]*readModel.UserDTO, error) {
	return s.users.FindAll()
}
