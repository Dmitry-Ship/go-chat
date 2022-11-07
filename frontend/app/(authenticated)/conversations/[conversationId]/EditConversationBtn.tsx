import React, { FormEvent, useEffect, useState } from "react";
import styles from "./EditConversationBtn.module.css";
import SlideIn from "../../../../src/components/common/SlideIn";
import { useRouter } from "next/router";
import InviteMenu from "./InviteMenu";
import { useAPI } from "../../../../src/contexts/apiContext";
import { Conversation } from "../../../../src/types/coreTypes";

const EditConversationBtn: React.FC<{
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
  const { makeCommand } = useAPI();

  const handleClose = () => {
    setIsEditing(false);
  };

  const handleLeave = async () => {
    await makeCommand("/leaveConversation", {
      conversation_id: conversation.id,
    });
    onLeave();
    router.push("/");
    setIsEditing(false);
  };

  const handleDelete = async () => {
    const result = await makeCommand("/deleteConversation", {
      conversation_id: conversation.id,
    });

    if (result.status) {
      router.push("/");
      setIsEditing(false);
    }
  };

  const handleRename = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const result = await makeCommand("/renameConversation", {
      conversation_id: conversation.id,
      new_name: newName,
    });

    if (result.status) {
      setIsEditing(false);
    }

    setNewName(conversation.name);
  };

  useEffect(() => {
    if (!isEditing) {
      setNewName(conversation.name);
    }
  }, [isEditing, conversation.name]);

  return (
    <>
      <button onClick={() => setIsEditing(true)} className={styles.editButton}>
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

export default EditConversationBtn;
