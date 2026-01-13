import { useQuery } from '@tanstack/react-query'
import { getConversationUsers, getMessages } from '@/lib/api'

export function useMessages(conversationId: string | null, page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ['messages', conversationId, page, pageSize],
    queryFn: () => getMessages(conversationId!, page, pageSize),
    enabled: !!conversationId,
    staleTime: 10 * 60 * 1000,
  })
}

export function useConversationUsers(conversationId: string | null, userIds: string[]) {
  const idsKey = userIds.join(',');
  return useQuery({
    queryKey: ['conversation-users', conversationId, idsKey],
    queryFn: () => getConversationUsers(conversationId!, userIds),
    enabled: !!conversationId && userIds.length > 0,
    staleTime: 10 * 60 * 1000,
  })
}
