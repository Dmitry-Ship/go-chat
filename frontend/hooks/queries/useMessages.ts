import { useQuery } from '@tanstack/react-query'
import { getMessages } from '@/lib/api'

export function useMessages(conversationId: string | null, page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ['messages', conversationId, page, pageSize],
    queryFn: () => getMessages(conversationId!, page, pageSize),
    enabled: !!conversationId,
    staleTime: 10 * 60 * 1000,
  })
}
