"use client";

import { useState, useMemo } from "react";
import { usePotentialInvitees } from "@/hooks/queries/usePotentialInvitees";
import { useInviteUser } from "@/hooks/mutations/conversationMutations";
import {
  AlertDialog as Dialog,
  AlertDialogContent as DialogContent,
  AlertDialogDescription as DialogDescription,
  AlertDialogFooter as DialogFooter,
  AlertDialogHeader as DialogHeader,
  AlertDialogTitle as DialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface InviteUserDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  conversationId: string;
}

export const InviteUserDialog = ({
  open,
  onOpenChange,
  conversationId,
}: InviteUserDialogProps) => {
  const [search, setSearch] = useState("");
  const { data: users = [], refetch } = usePotentialInvitees(conversationId, open);
  const inviteUserMutation = useInviteUser();

  const filteredUsers = useMemo(() => 
    users.filter((user) =>
      user.name.toLowerCase().includes(search.toLowerCase())
    ),
    [users, search]
  );

  const handleInvite = async (userId: string) => {
    try {
      await inviteUserMutation.mutateAsync({ conversationId, userId });
      refetch();
    } catch (error) {
      console.error("Failed to invite user:", error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Invite User</DialogTitle>
          <DialogDescription>
            Select a user to invite to the conversation
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <Input
            placeholder="Search users..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />

          <div className="max-h-60 overflow-y-auto space-y-2">
            {filteredUsers.map((user) => (
              <div
                key={user.id}
                className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
              >
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white font-semibold">
                    {user.avatar}
                  </div>
                  <div>
                    <div className="font-medium">{user.name}</div>
                  </div>
                </div>

                <Button
                  size="sm"
                  onClick={() => handleInvite(user.id)}
                  disabled={inviteUserMutation.isPending}
                >
                  {inviteUserMutation.isPending ? "Inviting..." : "Invite"}
                </Button>
              </div>
            ))}
            {filteredUsers.length === 0 && (
              <div className="text-center text-gray-500 py-4">
                No users available to invite
              </div>
            )}
          </div>
        </div>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
          >
            Cancel
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
