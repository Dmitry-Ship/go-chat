import { useInfiniteQuery } from '@tanstack/react-query'
import type { InfiniteData } from '@tanstack/react-query'
import { getConversations } from '@/lib/api'
import type { ConversationDTO, ConversationPageResponse } from '@/lib/types'

export function useConversations(limit = 20) {
  const query = useInfiniteQuery<ConversationPageResponse, Error, InfiniteData<ConversationPageResponse>, (string | number)[], string | null>({
    queryKey: ['conversations', limit],
    queryFn: ({ pageParam }) => getConversations(pageParam, limit),
    initialPageParam: null as string | null,
    getNextPageParam: (lastPage) => lastPage.next_cursor ?? undefined,
    staleTime: 5 * 60 * 1000,
  })

  const conversations: ConversationDTO[] = query.data
    ? query.data.pages.flatMap((page) => page.conversations)
    : []

  return {
    ...query,
    conversations,
  }
}
