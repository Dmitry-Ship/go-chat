import { useQuery } from '@tanstack/react-query'
import { getConversations } from '@/lib/api'

export function useConversations(page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ['conversations', page, pageSize],
    queryFn: () => getConversations(page, pageSize),
    staleTime: 5 * 60 * 1000,
  })
}
