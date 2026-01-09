import { useQuery } from '@tanstack/react-query'
import { getParticipants } from '@/lib/api'

export function useParticipants(conversationId: string | null, page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ['participants', conversationId, page, pageSize],
    queryFn: () => getParticipants(conversationId!, page, pageSize),
    enabled: !!conversationId,
    staleTime: 10 * 60 * 1000,
  })
}
