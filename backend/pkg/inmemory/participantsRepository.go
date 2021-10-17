package inmemory

import (
	"GitHub/go-chat/backend/domain"
	"errors"

	"github.com/google/uuid"
)

type participantRepository struct {
	participants map[uuid.UUID]*domain.Participant
}

func NewParticipantRepository() *participantRepository {
	return &participantRepository{
		participants: make(map[uuid.UUID]*domain.Participant),
	}
}

func (r *participantRepository) Create(participant *domain.Participant) error {
	for _, currentParticipant := range r.participants {
		if currentParticipant.RoomId == participant.RoomId && currentParticipant.UserId == participant.UserId {
			return errors.New("participant already exists")
		}
	}

	r.participants[participant.Id] = participant

	return nil
}

func (r *participantRepository) FindAllByRoomID(roomID uuid.UUID) ([]*domain.Participant, error) {
	participants := make([]*domain.Participant, 0, len(r.participants))
	for _, participant := range r.participants {
		if participant.RoomId == roomID {
			participants = append(participants, participant)
		}
	}
	return participants, nil
}

func (r *participantRepository) FindByRoomIDAndUserID(roomID uuid.UUID, userID uuid.UUID) (*domain.Participant, error) {
	for _, participant := range r.participants {
		if participant.RoomId == roomID && participant.UserId == userID {
			return participant, nil
		}
	}
	return nil, errors.New("participant not found")
}

func (r *participantRepository) DeleteAllByRoomID(roomID uuid.UUID) error {
	for _, participant := range r.participants {
		if participant.RoomId == roomID {
			delete(r.participants, participant.Id)
		}
	}
	return nil
}

func (r *participantRepository) DeleteByRoomIDAndUserID(roomID uuid.UUID, userID uuid.UUID) error {
	for _, participant := range r.participants {
		if participant.RoomId == roomID && participant.UserId == userID {
			delete(r.participants, participant.Id)
			return nil
		}
	}
	return errors.New("participants not found")
}
