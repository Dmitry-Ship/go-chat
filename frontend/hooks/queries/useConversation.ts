import { useQuery } from '@tanstack/react-query'
import { getConversation } from '@/lib/api'

export function useConversation(conversationId: string | null) {
  return useQuery({
    queryKey: ['conversation', conversationId],
    queryFn: () => getConversation(conversationId!),
    enabled: !!conversationId,
    staleTime: 10 * 60 * 1000,
  })
}
