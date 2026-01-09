"use client";

import { useState } from "react";
import { useConversation } from "@/hooks/queries/useConversation";
import { useParticipants } from "@/hooks/queries/useParticipants";
import { useKickUser } from "@/hooks/mutations/conversationMutations";
import { useAuth } from "@/contexts/AuthContext";
import {
  AlertDialog as Dialog,
  AlertDialogContent as DialogContent,
  AlertDialogDescription as DialogDescription,
  AlertDialogFooter as DialogFooter,
  AlertDialogHeader as DialogHeader,
  AlertDialogTitle as DialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { UserPlus, UserMinus } from "lucide-react";
import { InviteUserDialog } from "./InviteUserDialog";

interface ParticipantsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  conversationId: string;
}

export const ParticipantsDialog = ({
  open,
  onOpenChange,
  conversationId,
}: ParticipantsDialogProps) => {
  const { data: participants = [] } = useParticipants(conversationId);
  const { data: conversation } = useConversation(conversationId);
  const { user: currentUser } = useAuth();
  const kickUserMutation = useKickUser();
  const [showInvite, setShowInvite] = useState(false);

  const handleKick = async (userId: string) => {
    try {
      await kickUserMutation.mutateAsync({ conversationId, userId });
    } catch (error) {
      console.error("Failed to kick user:", error);
    }
  };

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Participants</DialogTitle>
            <DialogDescription>
              Manage participants for this conversation
            </DialogDescription>
          </DialogHeader>

          <div className="max-h-80 overflow-y-auto space-y-2 py-4">
            {participants.map((participant) => (
              <div
                key={participant.id}
                className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
              >
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center text-white font-semibold">
                    {participant.avatar}
                  </div>
                  <div>
                    <div className="font-medium">{participant.name}</div>
                    {conversation?.is_owner && (
                      <Badge variant="outline" className="text-xs">
                        Owner
                      </Badge>
                    )}
                  </div>
                </div>

                {conversation?.is_owner && participant.id !== currentUser?.id && (
                  <Button
                    size="icon"
                    variant="ghost"
                    onClick={() => handleKick(participant.id)}
                    disabled={kickUserMutation.isPending}
                  >
                    <UserMinus className="h-4 w-4 text-red-600" />
                  </Button>
                )}
              </div>
            ))}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setShowInvite(true)}
              disabled={!conversation?.is_owner}
            >
              <UserPlus className="mr-2 h-4 w-4" />
              Invite User
            </Button>
            <Button type="button" onClick={() => onOpenChange(false)}>
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <InviteUserDialog
        open={showInvite}
        onOpenChange={setShowInvite}
        conversationId={conversationId}
      />
    </>
  );
};
