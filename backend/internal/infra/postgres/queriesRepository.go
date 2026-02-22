package postgres

import (
	"context"
	"time"

	"GitHub/go-chat/backend/internal/domain"
	"GitHub/go-chat/backend/internal/infra/postgres/db"
	"GitHub/go-chat/backend/internal/presentation"
	"GitHub/go-chat/backend/internal/readModel"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type queriesRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func optionalString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
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
	user, err := r.queries.GetUserByID(context.Background(), uuidToPgtype(id))
	if err != nil {
		return readModel.UserDTO{}, err
	}

	return readModel.UserDTO{
		ID:     pgtypeToUUID(user.ID),
		Name:   user.Name,
		Avatar: user.Avatar.String,
	}, nil
}

func (r *queriesRepository) GetUsersByIDs(ids []uuid.UUID) ([]readModel.UserDTO, error) {
	if len(ids) == 0 {
		return []readModel.UserDTO{}, nil
	}

	pgIDs := make([]pgtype.UUID, len(ids))
	for i, id := range ids {
		pgIDs[i] = uuidToPgtype(id)
	}

	users, err := r.queries.GetUsersByIDs(context.Background(), pgIDs)
	if err != nil {
		return nil, err
	}

	usersDTO := make([]readModel.UserDTO, len(users))
	for i, user := range users {
		usersDTO[i] = readModel.UserDTO{
			ID:     pgtypeToUUID(user.ID),
			Name:   user.Name,
			Avatar: user.Avatar.String,
		}
	}

	return usersDTO, nil
}

func (r *queriesRepository) GetConversationMessages(conversationID uuid.UUID, cursor *readModel.MessageCursor, limit int) (readModel.MessagePageDTO, error) {
	pageLimit := limit
	if pageLimit <= 0 {
		pageLimit = 50
	} else if pageLimit > 100 {
		pageLimit = 100
	}

	queryLimit := int32(pageLimit + 1)
	conversationIDPg := uuidToPgtype(conversationID)

	var rawMessages []readModel.RawMessageDTO
	if cursor == nil {
		messages, err := r.queries.GetConversationMessagesFirstPageRaw(context.Background(), db.GetConversationMessagesFirstPageRawParams{
			ConversationID: conversationIDPg,
			PageLimit:      queryLimit,
		})
		if err != nil {
			return readModel.MessagePageDTO{}, err
		}

		rawMessages = make([]readModel.RawMessageDTO, len(messages))
		for i, msg := range messages {
			rawMessages[i] = readModel.RawMessageDTO{
				ID:             pgtypeToUUID(msg.ID),
				Type:           uint8(msg.Type),
				CreatedAt:      msg.CreatedAt.Time,
				ConversationID: pgtypeToUUID(msg.ConversationID),
				Content:        msg.Content,
				UserID:         pgtypeToUUID(msg.UserID),
			}
		}
	} else {
		messages, err := r.queries.GetConversationMessagesPageRaw(context.Background(), db.GetConversationMessagesPageRawParams{
			ConversationID:  conversationIDPg,
			CursorCreatedAt: pgtype.Timestamptz{Time: cursor.CreatedAt, Valid: true},
			CursorID:        uuidToPgtype(cursor.ID),
			PageLimit:       queryLimit,
		})
		if err != nil {
			return readModel.MessagePageDTO{}, err
		}

		rawMessages = make([]readModel.RawMessageDTO, len(messages))
		for i, msg := range messages {
			rawMessages[i] = readModel.RawMessageDTO{
				ID:             pgtypeToUUID(msg.ID),
				Type:           uint8(msg.Type),
				CreatedAt:      msg.CreatedAt.Time,
				ConversationID: pgtypeToUUID(msg.ConversationID),
				Content:        msg.Content,
				UserID:         pgtypeToUUID(msg.UserID),
			}
		}
	}

	hasMore := false
	if len(rawMessages) > pageLimit {
		hasMore = true
		rawMessages = rawMessages[:len(rawMessages)-1]
	}

	formatter := presentation.NewMessageFormatter()
	messageDTOs := make([]readModel.MessageDTO, 0, len(rawMessages))

	for i := len(rawMessages) - 1; i >= 0; i-- {
		messageDTOs = append(messageDTOs, formatter.FormatMessageDTO(rawMessages[i]))
	}

	response := readModel.MessagePageDTO{
		Messages: messageDTOs,
		HasMore:  hasMore,
	}

	if hasMore && len(messageDTOs) > 0 {
		oldest := messageDTOs[0]
		response.NextCursor = oldest.CreatedAt.UTC().Format(time.RFC3339Nano) + "|" + oldest.ID.String()
	}

	return response, nil
}

func (r *queriesRepository) GetUserConversations(userID uuid.UUID, cursor *readModel.ConversationCursor, limit int) (readModel.ConversationPageDTO, error) {
	pageLimit := limit
	if pageLimit <= 0 {
		pageLimit = 50
	} else if pageLimit > 100 {
		pageLimit = 100
	}

	queryLimit := int32(pageLimit + 1)
	userIDPg := uuidToPgtype(userID)
	formatter := presentation.NewMessageFormatter()

	toConversationDTO := func(
		conversationID pgtype.UUID,
		conversationType int32,
		messageID pgtype.UUID,
		messageCreatedAt pgtype.Timestamptz,
		messageContent string,
		messageType int32,
		messageUserID pgtype.UUID,
		messageUserName pgtype.Text,
		messageUserAvatar pgtype.Text,
		groupAvatar pgtype.Text,
		groupName pgtype.Text,
		otherUserID pgtype.UUID,
		otherUserName string,
		otherUserAvatar pgtype.Text,
	) readModel.ConversationDTO {
		conversationDTO := readModel.ConversationDTO{
			ID:   pgtypeToUUID(conversationID),
			Type: conversationTypesMap[uint8(conversationType)].String(),
		}

		if messageID.Valid {
			rawLastMessage := readModel.RawLastMessageDTO{
				MessageID:         pgtypeToUUID(messageID),
				MessageCreatedAt:  messageCreatedAt.Time,
				MessageContent:    messageContent,
				MessageType:       messageType,
				MessageUserID:     pgtypeToUUID(messageUserID),
				MessageUserName:   messageUserName.String,
				MessageUserAvatar: messageUserAvatar.String,
				ConversationID:    pgtypeToUUID(conversationID),
			}
			conversationDTO.LastMessage = formatter.FormatConversationLastMessage(rawLastMessage)
		}

		switch conversationTypesMap[uint8(conversationType)] {
		case domain.ConversationTypeDirect:
			if otherUserID.Valid {
				conversationDTO.Avatar = otherUserAvatar.String
				conversationDTO.Name = optionalString(otherUserName)
			}
		case domain.ConversationTypeGroup:
			if groupAvatar.Valid && groupName.Valid {
				conversationDTO.Avatar = groupAvatar.String
				conversationDTO.Name = groupName.String
			}
		}

		return conversationDTO
	}

	if cursor == nil {
		queryResults, err := r.queries.GetUserConversationsFirstPage(context.Background(), db.GetUserConversationsFirstPageParams{
			UserID:    userIDPg,
			PageLimit: queryLimit,
		})
		if err != nil {
			return readModel.ConversationPageDTO{}, err
		}

		hasMore := false
		if len(queryResults) > pageLimit {
			hasMore = true
			queryResults = queryResults[:len(queryResults)-1]
		}

		conversationDTOs := make([]readModel.ConversationDTO, len(queryResults))
		for i, result := range queryResults {
			conversationDTOs[i] = toConversationDTO(
				result.ConversationID,
				result.Type,
				result.MessageID,
				result.MessageCreatedAt,
				result.MessageContent,
				result.MessageType,
				result.MessageUserID,
				result.MessageUserName,
				result.MessageUserAvatar,
				result.GroupAvatar,
				result.GroupName,
				result.OtherUserID,
				result.OtherUserName,
				result.OtherUserAvatar,
			)
		}

		response := readModel.ConversationPageDTO{
			Conversations: conversationDTOs,
			HasMore:       hasMore,
		}

		if hasMore && len(queryResults) > 0 {
			oldest := queryResults[len(queryResults)-1]
			response.NextCursor = oldest.CreatedAt.Time.UTC().Format(time.RFC3339Nano) + "|" + pgtypeToUUID(oldest.ConversationID).String()
		}

		return response, nil
	}

	queryResults, err := r.queries.GetUserConversationsPage(context.Background(), db.GetUserConversationsPageParams{
		UserID:          userIDPg,
		CursorCreatedAt: pgtype.Timestamptz{Time: cursor.CreatedAt, Valid: true},
		CursorID:        uuidToPgtype(cursor.ID),
		PageLimit:       queryLimit,
	})
	if err != nil {
		return readModel.ConversationPageDTO{}, err
	}

	hasMore := false
	if len(queryResults) > pageLimit {
		hasMore = true
		queryResults = queryResults[:len(queryResults)-1]
	}

	conversationDTOs := make([]readModel.ConversationDTO, len(queryResults))
	for i, result := range queryResults {
		conversationDTOs[i] = toConversationDTO(
			result.ConversationID,
			result.Type,
			result.MessageID,
			result.MessageCreatedAt,
			result.MessageContent,
			result.MessageType,
			result.MessageUserID,
			result.MessageUserName,
			result.MessageUserAvatar,
			result.GroupAvatar,
			result.GroupName,
			result.OtherUserID,
			result.OtherUserName,
			result.OtherUserAvatar,
		)
	}

	response := readModel.ConversationPageDTO{
		Conversations: conversationDTOs,
		HasMore:       hasMore,
	}

	if hasMore && len(queryResults) > 0 {
		oldest := queryResults[len(queryResults)-1]
		response.NextCursor = oldest.CreatedAt.Time.UTC().Format(time.RFC3339Nano) + "|" + pgtypeToUUID(oldest.ConversationID).String()
	}

	return response, nil
}

func (r *queriesRepository) GetConversation(id uuid.UUID, userID uuid.UUID) (readModel.ConversationFullDTO, error) {
	result, err := r.queries.GetConversationBase(context.Background(), db.GetConversationBaseParams{
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
		otherUsers, err := r.queries.GetDirectConversationOtherUser(context.Background(), db.GetDirectConversationOtherUserParams{
			ConversationID: uuidToPgtype(id),
			UserID:         uuidToPgtype(userID),
		})
		if err != nil {
			return readModel.ConversationFullDTO{}, err
		}
		if len(otherUsers) > 0 {
			conversationDTO.Avatar = otherUsers[0].Avatar.String
			conversationDTO.Name = optionalString(otherUsers[0].Name)
		}
	case domain.ConversationTypeGroup:
		if result.GroupAvatar.Valid && result.GroupName.Valid {
			conversationDTO.Avatar = result.GroupAvatar.String
			conversationDTO.Name = result.GroupName.String
		}
		if result.GroupOwnerID.Valid {
			conversationDTO.IsOwner = pgtypeToUUID(result.GroupOwnerID) == userID
		}
		participantsCount, err := r.queries.CountConversationParticipants(context.Background(), uuidToPgtype(id))
		if err != nil {
			return readModel.ConversationFullDTO{}, err
		}
		conversationDTO.ParticipantsCount = participantsCount
		conversationDTO.HasJoined = true
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

func (r *queriesRepository) InviteToConversationAtomic(conversationID uuid.UUID, inviteeID uuid.UUID, participantID uuid.UUID) (uuid.UUID, error) {
	result, err := r.queries.InviteToConversationAtomic(context.Background(), db.InviteToConversationAtomicParams{
		ConversationID: uuidToPgtype(conversationID),
		InviteeID:      uuidToPgtype(inviteeID),
		ParticipantID:  uuidToPgtype(participantID),
	})

	if err != nil {
		return uuid.Nil, err
	}

	return pgtypeToUUID(result), nil
}

func (r *queriesRepository) LeaveConversationAtomic(conversationID uuid.UUID, userID uuid.UUID) (int64, error) {
	return r.queries.LeaveConversationAtomic(context.Background(), db.LeaveConversationAtomicParams{
		ConversationID: uuidToPgtype(conversationID),
		UserID:         uuidToPgtype(userID),
	})
}
