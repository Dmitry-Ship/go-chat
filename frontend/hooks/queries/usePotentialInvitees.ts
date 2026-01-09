import { useQuery } from '@tanstack/react-query'
import { getPotentialInvitees } from '@/lib/api'

export function usePotentialInvitees(conversationId: string | null, enabled = true) {
  return useQuery({
    queryKey: ['potentialInvitees', conversationId],
    queryFn: () => getPotentialInvitees(conversationId!),
    enabled: !!conversationId && enabled,
    staleTime: 5 * 60 * 1000,
  })
}
