import { useQuery } from '@tanstack/react-query'
import { getContacts } from '@/lib/api'

export function useContacts(page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ['contacts', page, pageSize],
    queryFn: () => getContacts(page, pageSize),
    staleTime: 15 * 60 * 1000,
  })
}
