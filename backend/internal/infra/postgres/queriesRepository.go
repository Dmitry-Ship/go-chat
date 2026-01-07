package postgres

import (
	"context"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres/db"
	"GitHub/go-chat/backend/internal/readModel"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type queriesRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewQueriesRepository(pool *pgxpool.Pool) *queriesRepository {
	return &queriesRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *queriesRepository) paginate(paginationInfo readModel.PaginationInfo) (limit int32, offset int32) {
	page := paginationInfo.GetPage()

	if page == 0 {
		page = 1
	}

	pageSize := paginationInfo.GetPageSize()

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 50
	}

	offsetVal := (page - 1) * pageSize
	return int32(pageSize), int32(offsetVal)
}

func (r *queriesRepository) GetContacts(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	limit, offset := r.paginate(paginationInfo)

	users, err := r.queries.GetContacts(context.Background(), db.GetContactsParams{
		ID:     uuidToPgtype(userID),
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, err
	}

	usersDTO := make([]readModel.ContactDTO, len(users))
	for i, user := range users {
		usersDTO[i] = readModel.ContactDTO{
			ID:     pgtypeToUUID(user.ID),
			Name:   user.Name,
			Avatar: user.Avatar.String,
		}
	}

	return usersDTO, nil
}

func (r *queriesRepository) GetParticipants(conversationID uuid.UUID, userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	limit, offset := r.paginate(paginationInfo)

	participants, err := r.queries.GetParticipantsByConversationID(context.Background(), db.GetParticipantsByConversationIDParams{
		ConversationID: uuidToPgtype(conversationID),
		Limit:          limit,
		Offset:         offset,
	})

	if err != nil {
		return nil, err
	}

	usersDTO := make([]readModel.ContactDTO, len(participants))
	for i, participant := range participants {
		usersDTO[i] = readModel.ContactDTO{
			ID:     pgtypeToUUID(participant.ID),
			Name:   participant.Name,
			Avatar: participant.Avatar.String,
		}
	}

	return usersDTO, nil
}

func (r *queriesRepository) GetPotentialInvitees(conversationID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ContactDTO, error) {
	limit, offset := r.paginate(paginationInfo)

	users, err := r.queries.GetPotentialInvitees(context.Background(), db.GetPotentialInviteesParams{
		ConversationID: uuidToPgtype(conversationID),
		Limit:          limit,
		Offset:         offset,
	})

	if err != nil {
		return nil, err
	}

	usersDTO := make([]readModel.ContactDTO, len(users))
	for i, user := range users {
		usersDTO[i] = readModel.ContactDTO{
			ID:     pgtypeToUUID(user.ID),
			Name:   user.Name,
			Avatar: user.Avatar.String,
		}
	}

	return usersDTO, nil
}

func (r *queriesRepository) GetUserByID(id uuid.UUID) (readModel.UserDTO, error) {
	user, err := r.queries.GetUserByIDDTO(context.Background(), uuidToPgtype(id))
	if err != nil {
		return readModel.UserDTO{}, err
	}

	return readModel.UserDTO{
		ID:     pgtypeToUUID(user.ID),
		Name:   user.Name,
		Avatar: user.Avatar.String,
	}, nil
}

func (r *queriesRepository) GetConversationMessages(conversationID uuid.UUID, requestUserID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.MessageDTO, error) {
	limit, offset := r.paginate(paginationInfo)

	messages, err := r.queries.GetConversationMessagesWithUser(context.Background(), db.GetConversationMessagesWithUserParams{
		ConversationID: uuidToPgtype(conversationID),
		Limit:          limit,
		Offset:         offset,
	})

	if err != nil {
		return nil, err
	}

	messageDTOs := make([]readModel.MessageDTO, len(messages))
	for i, msg := range messages {
		createdAt := msg.CreatedAt.Time

		text := ""
		switch messageTypesMap[uint8(msg.Type)] {
		case domain.MessageTypeText:
			text = msg.Content
		case domain.MessageTypeRenamedConversation:
			text = msg.UserName.String + " renamed chat to " + msg.Content
		case domain.MessageTypeJoinedConversation:
			text = msg.UserName.String + " joined"
		case domain.MessageTypeLeftConversation:
			text = msg.UserName.String + " left"
		case domain.MessageTypeInvitedConversation:
			text = msg.UserName.String + " was invited"
		}

		messageDTO := readModel.MessageDTO{
			ID:             pgtypeToUUID(msg.ID),
			CreatedAt:      createdAt,
			Text:           text,
			Type:           messageTypesMap[uint8(msg.Type)].String(),
			ConversationId: pgtypeToUUID(msg.ConversationID),
			User: readModel.UserDTO{
				ID:     pgtypeToUUID(msg.UserID),
				Avatar: msg.UserAvatar.String,
				Name:   msg.UserName.String,
			},
		}

		if messageTypesMap[uint8(msg.Type)] == domain.MessageTypeText {
			messageDTO.IsInbound = pgtypeToUUID(msg.UserID) != requestUserID
		}

		messageDTOs[i] = messageDTO
	}

	return messageDTOs, nil
}

func (r *queriesRepository) GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (readModel.MessageDTO, error) {
	msg, err := r.queries.GetNotificationMessageWithUser(context.Background(), uuidToPgtype(messageID))
	if err != nil {
		return readModel.MessageDTO{}, err
	}

	createdAt := msg.CreatedAt.Time

	text := ""
	switch messageTypesMap[uint8(msg.Type)] {
	case domain.MessageTypeText:
		text = msg.Content
	case domain.MessageTypeRenamedConversation:
		text = msg.UserName.String + " renamed chat to " + msg.Content
	case domain.MessageTypeJoinedConversation:
		text = msg.UserName.String + " joined"
	case domain.MessageTypeLeftConversation:
		text = msg.UserName.String + " left"
	case domain.MessageTypeInvitedConversation:
		text = msg.UserName.String + " was invited"
	}

	messageDTO := readModel.MessageDTO{
		ID:             pgtypeToUUID(msg.ID),
		CreatedAt:      createdAt,
		Text:           text,
		Type:           messageTypesMap[uint8(msg.Type)].String(),
		ConversationId: pgtypeToUUID(msg.ConversationID),
		User: readModel.UserDTO{
			ID:     pgtypeToUUID(msg.UserID),
			Avatar: msg.UserAvatar.String,
			Name:   msg.UserName.String,
		},
	}

	if messageTypesMap[uint8(msg.Type)] == domain.MessageTypeText {
		messageDTO.IsInbound = pgtypeToUUID(msg.UserID) != requestUserID
	}

	return messageDTO, nil
}

func (r *queriesRepository) GetUserConversations(userID uuid.UUID, paginationInfo readModel.PaginationInfo) ([]readModel.ConversationDTO, error) {
	limit, offset := r.paginate(paginationInfo)

	queryResults, err := r.queries.GetUserConversations(context.Background(), db.GetUserConversationsParams{
		UserID: uuidToPgtype(userID),
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, err
	}

	conversationDTOs := make([]readModel.ConversationDTO, len(queryResults))

	for i, result := range queryResults {
		conversationDTO := readModel.ConversationDTO{
			ID:   pgtypeToUUID(result.ConversationID),
			Type: conversationTypesMap[uint8(result.Type)].String(),
		}

		if result.MessageID.Valid {
			createdAt := result.MessageCreatedAt.Time

			text := ""
			switch messageTypesMap[uint8(result.MessageType.Int32)] {
			case domain.MessageTypeText:
				text = result.MessageContent.String
			case domain.MessageTypeRenamedConversation:
				text = result.MessageUserName.String + " renamed chat to " + result.MessageContent.String
			case domain.MessageTypeJoinedConversation:
				text = result.MessageUserName.String + " joined"
			case domain.MessageTypeLeftConversation:
				text = result.MessageUserName.String + " left"
			case domain.MessageTypeInvitedConversation:
				text = result.MessageUserName.String + " was invited"
			}

			msgID := uuid.UUID(result.MessageID.Bytes)
			msgUserID := uuid.UUID(result.MessageUserID.Bytes)

			messageDTO := readModel.MessageDTO{
				ID:             msgID,
				CreatedAt:      createdAt,
				Text:           text,
				Type:           messageTypesMap[uint8(result.MessageType.Int32)].String(),
				ConversationId: pgtypeToUUID(result.ConversationID),
				User: readModel.UserDTO{
					ID:     msgUserID,
					Avatar: result.MessageUserAvatar.String,
					Name:   result.MessageUserName.String,
				},
			}

			if messageTypesMap[uint8(result.MessageType.Int32)] == domain.MessageTypeText {
				messageDTO.IsInbound = msgUserID != userID
			}

			conversationDTO.LastMessage = messageDTO
		}

		switch conversationTypesMap[uint8(result.Type)] {
		case domain.ConversationTypeDirect:
			if result.OtherUserID.Valid {
				conversationDTO.Avatar = result.OtherUserAvatar.String
				conversationDTO.Name = result.OtherUserName.String
			}
		case domain.ConversationTypeGroup:
			if result.GroupAvatar.Valid && result.GroupName.Valid {
				conversationDTO.Avatar = result.GroupAvatar.String
				conversationDTO.Name = result.GroupName.String
			}
		}

		conversationDTOs[i] = conversationDTO
	}

	return conversationDTOs, nil
}

func (r *queriesRepository) GetConversation(id uuid.UUID, userID uuid.UUID) (readModel.ConversationFullDTO, error) {
	result, err := r.queries.GetConversationFull(context.Background(), db.GetConversationFullParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})

	if err != nil {
		return readModel.ConversationFullDTO{}, err
	}

	conversationDTO := readModel.ConversationFullDTO{
		ID:        pgtypeToUUID(result.ConversationID),
		CreatedAt: result.CreatedAt.Time,
		Type:      conversationTypesMap[uint8(result.Type)].String(),
	}

	switch conversationTypesMap[uint8(result.Type)] {
	case domain.ConversationTypeDirect:
		if result.OtherUserID.Valid {
			conversationDTO.Avatar = result.OtherUserAvatar.String
			conversationDTO.Name = result.OtherUserName.String
		}
	case domain.ConversationTypeGroup:
		if result.GroupAvatar.Valid && result.GroupName.Valid {
			conversationDTO.Avatar = result.GroupAvatar.String
			conversationDTO.Name = result.GroupName.String
		}
		if result.GroupOwnerID.Valid {
			conversationDTO.IsOwner = pgtypeToUUID(result.GroupOwnerID) == userID
		}
		conversationDTO.ParticipantsCount = result.ParticipantsCount
		conversationDTO.HasJoined = result.UserParticipantID.Valid
	}

	return conversationDTO, nil
}
