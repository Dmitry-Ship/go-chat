import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  createConversation,
  startDirectConversation,
  joinConversation,
  leaveConversation,
  deleteConversation,
  renameConversation,
  inviteUser,
  kickUser,
} from '@/lib/api'

export function useCreateConversation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ name, id }: { name: string; id: string }) =>
      createConversation(name, id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] })
    },
  })
}

export function useStartDirectConversation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (toUserId: string) => startDirectConversation(toUserId),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] })
      queryClient.setQueryData(['conversation', data.conversation_id], {
        id: data.conversation_id,
      })
    },
  })
}

export function useJoinConversation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (conversationId: string) => joinConversation(conversationId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] })
    },
  })
}

export function useLeaveConversation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (conversationId: string) => leaveConversation(conversationId),
    onSuccess: (_, conversationId) => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] })
      queryClient.removeQueries({ queryKey: ['conversation', conversationId] })
      queryClient.removeQueries({ queryKey: ['messages', conversationId] })
    },
  })
}

export function useDeleteConversation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (conversationId: string) => deleteConversation(conversationId),
    onSuccess: (_, conversationId) => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] })
      queryClient.removeQueries({ queryKey: ['conversation', conversationId] })
      queryClient.removeQueries({ queryKey: ['messages', conversationId] })
      queryClient.removeQueries({ queryKey: ['participants', conversationId] })
    },
  })
}

export function useRenameConversation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ conversationId, newName }: { conversationId: string; newName: string }) =>
      renameConversation(conversationId, newName),
    onSuccess: (_, { conversationId }) => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] })
      queryClient.invalidateQueries({ queryKey: ['conversation', conversationId] })
    },
  })
}

export function useInviteUser() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ conversationId, userId }: { conversationId: string; userId: string }) =>
      inviteUser(conversationId, userId),
    onSuccess: (_, { conversationId }) => {
      queryClient.invalidateQueries({ queryKey: ['participants', conversationId] })
    },
  })
}

export function useKickUser() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ conversationId, userId }: { conversationId: string; userId: string }) =>
      kickUser(conversationId, userId),
    onSuccess: (_, { conversationId }) => {
      queryClient.invalidateQueries({ queryKey: ['participants', conversationId] })
    },
  })
}
