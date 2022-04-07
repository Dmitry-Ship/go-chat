import React, { useState } from "react";
import styles from "./EditConversationBtn.module.css";
import SlideIn from "../common/SlideIn";
import { makeCommand } from "../../api/fetch";
import { useRouter } from "next/router";

const EditConversationBtn: React.FC<{
  joined: boolean;
  onLeave: () => void;
  conversationId: string;
}> = ({ joined, onLeave, conversationId }) => {
  const [isEditing, setIsEditing] = useState(false);
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

  return (
    <>
      <button onClick={() => setIsEditing(true)} className={styles.editButton}>
        ⚙️
      </button>
      <SlideIn onClose={handleClose} isOpen={isEditing}>
        <>
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
