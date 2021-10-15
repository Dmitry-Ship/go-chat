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

func (r *participantRepository) FindByID(id uuid.UUID) (*domain.Participant, error) {
	participant, ok := r.participants[id]
	if !ok {
		return nil, errors.New("participant not found")
	}
	return participant, nil
}

func (r *participantRepository) Create(participant *domain.Participant) (*domain.Participant, error) {
	r.participants[participant.Id] = participant
	return participant, nil
}

func (r *participantRepository) Delete(id uuid.UUID) error {
	_, ok := r.participants[id]
	if !ok {
		return errors.New("participant not found")
	}
	delete(r.participants, id)
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
