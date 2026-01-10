package postgres

import (
	"context"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres/db"
	"GitHub/go-chat/backend/internal/presentation"
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

	messages, err := r.queries.GetConversationMessagesRaw(context.Background(), db.GetConversationMessagesRawParams{
		ConversationID: uuidToPgtype(conversationID),
		Limit:          limit,
		Offset:         offset,
	})

	if err != nil {
		return nil, err
	}

	formatter := presentation.NewMessageFormatter()
	messageDTOs := make([]readModel.MessageDTO, len(messages))

	for i, msg := range messages {
		rawMessage := readModel.RawMessageDTO{
			ID:             pgtypeToUUID(msg.ID),
			Type:           uint8(msg.Type),
			CreatedAt:      msg.CreatedAt.Time,
			ConversationID: pgtypeToUUID(msg.ConversationID),
			Content:        msg.Content,
			UserID:         pgtypeToUUID(msg.UserID),
			UserName:       msg.UserName.String,
			UserAvatar:     msg.UserAvatar.String,
		}

		messageDTOs[i] = formatter.FormatMessageDTO(rawMessage, requestUserID)
	}

	return messageDTOs, nil
}

func (r *queriesRepository) GetNotificationMessage(messageID uuid.UUID, requestUserID uuid.UUID) (readModel.MessageDTO, error) {
	msg, err := r.queries.GetNotificationMessageRaw(context.Background(), uuidToPgtype(messageID))
	if err != nil {
		return readModel.MessageDTO{}, err
	}

	formatter := presentation.NewMessageFormatter()
	rawMessage := readModel.RawMessageDTO{
		ID:             pgtypeToUUID(msg.ID),
		Type:           uint8(msg.Type),
		CreatedAt:      msg.CreatedAt.Time,
		ConversationID: pgtypeToUUID(msg.ConversationID),
		Content:        msg.Content,
		UserID:         pgtypeToUUID(msg.UserID),
		UserName:       msg.UserName.String,
		UserAvatar:     msg.UserAvatar.String,
	}

	return formatter.FormatMessageDTO(rawMessage, requestUserID), nil
}

func (r *queriesRepository) StoreMessageAndReturnWithUser(id uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, content string, messageType int32) (readModel.MessageDTO, error) {
	if err := r.queries.StoreMessage(context.Background(), db.StoreMessageParams{
		ID:             uuidToPgtype(id),
		ConversationID: uuidToPgtype(conversationID),
		UserID:         uuidToPgtype(userID),
		Content:        content,
		Type:           messageType,
	}); err != nil {
		return readModel.MessageDTO{}, err
	}

	msg, err := r.queries.GetMessageWithUser(context.Background(), uuidToPgtype(id))
	if err != nil {
		return readModel.MessageDTO{}, err
	}

	formatter := presentation.NewMessageFormatter()
	rawMessage := readModel.RawMessageDTO{
		ID:             pgtypeToUUID(msg.ID),
		Type:           uint8(msg.Type),
		CreatedAt:      msg.CreatedAt.Time,
		ConversationID: pgtypeToUUID(msg.ConversationID),
		Content:        msg.Content,
		UserID:         pgtypeToUUID(msg.UserID),
		UserName:       msg.UserName,
		UserAvatar:     msg.UserAvatar.String,
	}

	return formatter.FormatMessageDTO(rawMessage, userID), nil
}

func (r *queriesRepository) StoreSystemMessageAndReturn(id uuid.UUID, conversationID uuid.UUID, userID uuid.UUID, content string, messageType int32) (readModel.MessageDTO, error) {
	msg, err := r.queries.StoreSystemMessageAndReturn(context.Background(), db.StoreSystemMessageAndReturnParams{
		ID:             uuidToPgtype(id),
		ConversationID: uuidToPgtype(conversationID),
		ID_2:           uuidToPgtype(userID),
		Content:        content,
		Type:           messageType,
	})
	if err != nil {
		return readModel.MessageDTO{}, err
	}

	formatter := presentation.NewMessageFormatter()
	rawMessage := readModel.RawMessageDTO{
		ID:             pgtypeToUUID(msg.ID),
		Type:           uint8(msg.Type),
		CreatedAt:      msg.CreatedAt.Time,
		ConversationID: pgtypeToUUID(msg.ConversationID),
		Content:        msg.FormattedText,
		UserID:         pgtypeToUUID(msg.UserID),
		UserName:       msg.UserName,
		UserAvatar:     msg.UserAvatar.String,
	}

	return formatter.FormatMessageDTO(rawMessage, userID), nil
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

	formatter := presentation.NewMessageFormatter()
	conversationDTOs := make([]readModel.ConversationDTO, len(queryResults))

	for i, result := range queryResults {
		conversationDTO := readModel.ConversationDTO{
			ID:   pgtypeToUUID(result.ConversationID),
			Type: conversationTypesMap[uint8(result.Type)].String(),
		}

		if result.MessageID.Valid {
			rawLastMessage := readModel.RawLastMessageDTO{
				MessageID:         pgtypeToUUID(result.MessageID),
				MessageCreatedAt:  result.MessageCreatedAt.Time,
				MessageContent:    result.MessageContent.String,
				MessageType:       result.MessageType.Int32,
				MessageUserID:     pgtypeToUUID(result.MessageUserID),
				MessageUserName:   result.MessageUserName.String,
				MessageUserAvatar: result.MessageUserAvatar.String,
				ConversationID:    pgtypeToUUID(result.ConversationID),
			}
			conversationDTO.LastMessage = formatter.FormatConversationLastMessage(rawLastMessage, userID)
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

func (r *queriesRepository) IsMember(conversationID uuid.UUID, userID uuid.UUID) (bool, error) {
	return r.queries.IsMember(context.Background(), db.IsMemberParams{
		ConversationID: uuidToPgtype(conversationID),
		UserID:         uuidToPgtype(userID),
	})
}

func (r *queriesRepository) IsMemberOwner(conversationID uuid.UUID, userID uuid.UUID) (bool, error) {
	return r.queries.IsMemberOwner(context.Background(), db.IsMemberOwnerParams{
		ConversationID: uuidToPgtype(conversationID),
		UserID:         uuidToPgtype(userID),
	})
}

func (r *queriesRepository) RenameConversationAndReturn(conversationID uuid.UUID, name string) error {
	rowsAffected, err := r.queries.RenameConversationAndReturn(context.Background(), db.RenameConversationAndReturnParams{
		ConversationID: uuidToPgtype(conversationID),
		Name:           name,
	})

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrorUserNotOwner
	}

	return nil
}

func (r *queriesRepository) InviteToConversationAtomic(conversationID uuid.UUID, inviteeID uuid.UUID, participantID uuid.UUID) (uuid.UUID, error) {
	result, err := r.queries.InviteToConversationAtomic(context.Background(), db.InviteToConversationAtomicParams{
		ConversationID: uuidToPgtype(conversationID),
		ID:             uuidToPgtype(inviteeID),
		ID_2:           uuidToPgtype(participantID),
	})

	if err != nil {
		return uuid.Nil, err
	}

	return pgtypeToUUID(result), nil
}

func (r *queriesRepository) KickParticipantAtomic(conversationID uuid.UUID, targetID uuid.UUID) (int64, error) {
	return r.queries.KickParticipantAtomic(context.Background(), db.KickParticipantAtomicParams{
		ConversationID: uuidToPgtype(conversationID),
		UserID:         uuidToPgtype(targetID),
	})
}

func (r *queriesRepository) LeaveConversationAtomic(conversationID uuid.UUID, userID uuid.UUID) (int64, error) {
	return r.queries.LeaveConversationAtomic(context.Background(), db.LeaveConversationAtomicParams{
		ConversationID: uuidToPgtype(conversationID),
		UserID:         uuidToPgtype(userID),
	})
}
