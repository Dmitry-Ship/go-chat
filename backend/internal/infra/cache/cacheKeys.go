package cache

import (
	"fmt"
	"time"
)

const (
	CachePrefixUser         = "user"
	CachePrefixConversation = "conv"
	CachePrefixParticipants = "participants"
	CachePrefixConvMeta     = "conv_meta"
	CachePrefixUserConvList = "user_conv_list"

	TTLUser         = 15 * time.Minute
	TTLConversation = 15 * time.Minute
	TTLParticipants = 10 * time.Minute
	TTLConvMeta     = 15 * time.Minute
	TTLUserConvList = 5 * time.Minute
)

func UserKey(id string) string {
	return fmt.Sprintf("%s:%s", CachePrefixUser, id)
}

func UsernameKey(username string) string {
	return fmt.Sprintf("%s:username:%s", CachePrefixUser, username)
}

func ConversationKey(id string) string {
	return fmt.Sprintf("%s:%s", CachePrefixConversation, id)
}

func ParticipantsKey(conversationID string) string {
	return fmt.Sprintf("%s:%s", CachePrefixParticipants, conversationID)
}

func ConvMetaKey(conversationID string) string {
	return fmt.Sprintf("%s:%s", CachePrefixConvMeta, conversationID)
}

func UserConvListKey(userID string) string {
	return fmt.Sprintf("%s:%s", CachePrefixUserConvList, userID)
}

func UserKeysPattern(userID string) string {
	return fmt.Sprintf("%s:%s*", CachePrefixUser, userID)
}

func ConversationKeysPattern(conversationID string) string {
	return fmt.Sprintf("%s:%s*", CachePrefixConversation, conversationID)
}
