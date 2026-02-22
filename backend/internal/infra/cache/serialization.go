package cache

import (
	"encoding/json"

	"GitHub/go-chat/backend/internal/domain"
)

type UserCache struct {
	ID     string
	Name   string
	Avatar string
}

func SerializeUser(user *domain.User) ([]byte, error) {
	cacheUser := UserCache{
		ID:     user.ID.String(),
		Name:   user.Name,
		Avatar: user.Avatar,
	}
	return json.Marshal(cacheUser)
}

func DeserializeUser(data []byte) (*UserCache, error) {
	var cacheUser UserCache
	err := json.Unmarshal(data, &cacheUser)
	if err != nil {
		return nil, err
	}
	return &cacheUser, nil
}
