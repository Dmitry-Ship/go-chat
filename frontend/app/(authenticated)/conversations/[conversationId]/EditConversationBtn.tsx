import React, { FormEvent, useReducer, useState } from "react";
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

export function EditConversationBtn({
  onLeave,
  conversation,
  conversationId,
}: {
  conversation: Conversation & {
    joined: boolean;
    participants_count: number;
    is_owner: boolean;
  };
  conversationId: string;
  onLeave: () => void;
}) {
  const [isEditing, toggleEditing] = useReducer((editing) => !editing, false);
  const [newName, setNewName] = useState(conversation.name);
  const router = useRouter();

  const leaveConversationRequest = useMutation(leaveConversation, {
    onSuccess: (data) => {
      onLeave();
      router.push("/");
      toggleEditing();
    },
  });

  const deleteConversationRequest = useMutation(deleteConversation, {
    onSuccess: (data) => {
      router.push("/");
      toggleEditing();
    },
  });

  const renameConversationRequest = useMutation(renameConversation, {
    onSuccess: (data) => {
      toggleEditing();
    },
  });

  function handleLeave() {
    leaveConversationRequest.mutate({
      conversation_id: conversation.id,
    });
  }

  function handleDelete() {
    deleteConversationRequest.mutate({
      conversation_id: conversation.id,
    });
  }

  async function handleRename(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();

    renameConversationRequest.mutate({
      conversation_id: conversation.id,
      new_name: newName,
    });

    setNewName(conversation.name);
  }

  function handleStartEditing() {
    toggleEditing();
    setNewName(conversation.name);
  }

  return (
    <>
      <button onClick={handleStartEditing} className={styles.editButton}>
        ⚙️
      </button>

      <SlideIn onClose={toggleEditing} isOpen={isEditing}>
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
}
