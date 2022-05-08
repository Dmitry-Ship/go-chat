import React, { FormEvent, useEffect, useState } from "react";
import styles from "./EditConversationBtn.module.css";
import SlideIn from "../common/SlideIn";
import { makeCommand } from "../../api/fetch";
import { useRouter } from "next/router";
import InviteMenu from "./InviteMenu";

const EditConversationBtn: React.FC<{
  joined: boolean;
  onLeave: () => void;
  conversationId: string;
}> = ({ joined, onLeave, conversationId }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [newName, setNewName] = useState("");
  const router = useRouter();

  const handleClose = () => {
    setIsEditing(false);
  };

  const handleLeave = async () => {
    await makeCommand("/leaveConversation", {
      conversation_id: conversationId,
    });
    onLeave();
    router.push("/");
    setIsEditing(false);
  };

  const handleDelete = async () => {
    const result = await makeCommand("/deleteConversation", {
      conversation_id: conversationId,
    });

    if (result.status) {
      router.push("/");
      setIsEditing(false);
    }
  };

  const handleRename = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const result = await makeCommand("/renameConversation", {
      conversation_id: conversationId,
      new_name: newName,
    });

    if (result.status) {
      setIsEditing(false);
    }

    setNewName("");
  };

  useEffect(() => {
    if (!isEditing) {
      setNewName("");
    }
  }, [isEditing]);

  return (
    <>
      <button onClick={() => setIsEditing(true)} className={styles.editButton}>
        ⚙️
      </button>

      <SlideIn onClose={handleClose} isOpen={isEditing}>
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

          <InviteMenu />

          <button onClick={handleDelete} className={`btn ${styles.menuItem}`}>
            🗑 Delete
          </button>
          {joined && (
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
