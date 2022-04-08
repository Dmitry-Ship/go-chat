package application

import (
	"GitHub/go-chat/backend/domain"
	"GitHub/go-chat/backend/pkg/mappers"
)

type ContactsQueryService interface {
	GetContacts() ([]*mappers.UserDTO, error)
}

type contactsQueryService struct {
	users domain.UserRepository
}

func NewContactsQueryService(users domain.UserRepository) *contactsQueryService {
	return &contactsQueryService{
		users: users,
	}
}

func (s *contactsQueryService) GetContacts() ([]*mappers.UserDTO, error) {
	users, err := s.users.FindAll()

	if err != nil {
		return nil, err
	}

	var result []*mappers.UserDTO

	for _, user := range users {
		result = append(result, mappers.ToUserDTO(user))
	}

	return result, nil
}
