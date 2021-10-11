package inmemory

import (
	"GitHub/go-chat/backend/domain"
	"errors"
)

type participantRepository struct {
	participants map[int32]*domain.Participant
}

func NewParticipantRepository() *participantRepository {
	return &participantRepository{
		participants: make(map[int32]*domain.Participant),
	}
}

func (r *participantRepository) FindByID(id int32) (*domain.Participant, error) {
	participant, ok := r.participants[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return participant, nil
}

func (r *participantRepository) FindAll() ([]*domain.Participant, error) {
	participants := make([]*domain.Participant, 0, len(r.participants))
	for _, participant := range r.participants {
		participants = append(participants, participant)
	}
	return participants, nil
}

func (r *participantRepository) Create(participant *domain.Participant) (*domain.Participant, error) {
	r.participants[participant.Id] = participant
	return participant, nil
}

func (r *participantRepository) Update(participant *domain.Participant) error {
	_, ok := r.participants[participant.Id]
	if !ok {
		return errors.New("not found")
	}
	r.participants[participant.Id] = participant
	return nil
}

func (r *participantRepository) Delete(id int32) error {
	_, ok := r.participants[id]
	if !ok {
		return errors.New("not found")
	}
	delete(r.participants, id)
	return nil
}

func (r *participantRepository) FindByRoomID(roomID int32) ([]*domain.Participant, error) {
	participants := make([]*domain.Participant, 0, len(r.participants))
	for _, participant := range r.participants {
		if participant.RoomId == roomID {
			participants = append(participants, participant)
		}
	}
	return participants, nil
}

func (r *participantRepository) FindByUserID(userID int32) (*domain.Participant, error) {
	for _, participant := range r.participants {
		if participant.UserId == userID {
			return participant, nil
		}
	}
	return nil, errors.New("not found")
}
