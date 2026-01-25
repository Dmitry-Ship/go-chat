import { useInfiniteQuery, useQuery } from '@tanstack/react-query'
import { getConversationUsers, getMessages } from '@/lib/api'

export function useMessages(conversationId: string | null, limit = 20) {
  return useInfiniteQuery({
    queryKey: ['messages', conversationId, limit],
    queryFn: ({ pageParam }) => getMessages(conversationId!, pageParam ?? null, limit),
    initialPageParam: null,
    getNextPageParam: (lastPage) => lastPage.next_cursor ?? undefined,
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
