import React, { FormEvent, useEffect, useState } from "react";
import styles from "./EditConversationBtn.module.css";
import SlideIn from "../common/SlideIn";
import { useRouter } from "next/router";
import InviteMenu from "./InviteMenu";
import { useAPI } from "../../contexts/apiContext";
import { Conversation } from "../../types/coreTypes";

const EditConversationBtn: React.FC<{
  conversation: Conversation & {
    joined: boolean;
    participants_count: number;
    is_owner: boolean;
  };
  onLeave: () => void;
}> = ({ onLeave, conversation }) => {
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
              <form className={styles.menuItem} onSubmit={handleRename}>
                <input
                  type="text"
                  placeholder="New name"
                  size={32}
                  className={`${styles.menuItem} input`}
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                />
                <button type="submit" disabled={!newName} className={`btn`}>
                  Rename
                </button>
              </form>
              <button
                onClick={handleDelete}
                className={`btn ${styles.menuItem}`}
              >
                🗑 Delete
              </button>
            </>
          )}

          <InviteMenu />

          {conversation.joined && (
            <button onClick={handleLeave} className={`btn ${styles.menuItem}`}>
              ✌️ Leave
            </button>
          )}
        </>
      </SlideIn>
    </>
  );
};

export default EditConversationBtn;
