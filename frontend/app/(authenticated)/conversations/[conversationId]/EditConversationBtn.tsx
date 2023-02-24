import React, { FormEvent, useState } from "react";
import styles from "./EditConversationBtn.module.css";
import { SlideIn } from "../../../../src/components/common/SlideIn";
import { useRouter } from "next/navigation";
import { InviteMenu } from "./InviteMenu";
import { Conversation } from "../../../../src/types/coreTypes";
import { useMutation } from "react-query";
import {
  deleteConversation,
  leaveConversation,
  renameConversation,
} from "../../../../src/api/fetch";

export const EditConversationBtn: React.FC<{
  conversation: Conversation & {
    joined: boolean;
    participants_count: number;
    is_owner: boolean;
  };
  conversationId: string;
  onLeave: () => void;
}> = ({ onLeave, conversation, conversationId }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [newName, setNewName] = useState(conversation.name);
  const router = useRouter();

  const leaveConversationRequest = useMutation(leaveConversation, {
    onSuccess: (data) => {
      onLeave();
      router.push("/");
      setIsEditing(false);
    },
  });

  const deleteConversationRequest = useMutation(deleteConversation, {
    onSuccess: (data) => {
      router.push("/");
      setIsEditing(false);
    },
  });

  const renameConversationRequest = useMutation(renameConversation, {
    onSuccess: (data) => {
      setIsEditing(false);
    },
  });

  const handleClose = () => {
    setIsEditing(false);
  };

  const handleLeave = () => {
    leaveConversationRequest.mutate({
      conversation_id: conversation.id,
    });
  };

  const handleDelete = () => {
    deleteConversationRequest.mutate({
      conversation_id: conversation.id,
    });
  };

  const handleRename = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    renameConversationRequest.mutate({
      conversation_id: conversation.id,
      new_name: newName,
    });

    setNewName(conversation.name);
  };

  const handleStartEditing = () => {
    setIsEditing(true);
    setNewName(conversation.name);
  };

  return (
    <>
      <button onClick={handleStartEditing} className={styles.editButton}>
        ⚙️
      </button>

      <SlideIn onClose={handleClose} isOpen={isEditing}>
        <>
          {conversation.is_owner && (
            <>
              <form className={"m-b"} onSubmit={handleRename}>
                <input
                  type="text"
                  placeholder="New name"
                  size={32}
                  className={`m-b input`}
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                />
                <button
                  type="submit"
                  disabled={!newName || newName === conversation.name}
                  className={`btn`}
                >
                  Rename
                </button>
              </form>
              <button onClick={handleDelete} className={`btn m-b`}>
                🗑 Delete
              </button>
            </>
          )}

          <InviteMenu conversationId={conversationId} />

          {conversation.joined && (
            <button onClick={handleLeave} className={`btn m-b`}>
              ✌️ Leave
            </button>
          )}
        </>
      </SlideIn>
    </>
  );
};
