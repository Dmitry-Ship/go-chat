import React, { useState } from "react";
import styles from "./EditRoomBtn.module.css";
import SlideIn from "../common/SlideIn";
import { makeCommand } from "../../api/fetch";
import { useRouter } from "next/router";

const EditRoomBtn: React.FC<{
  joined: boolean;
  onLeave: () => void;
  roomId: string;
}> = ({ joined, onLeave, roomId }) => {
  const [isEditing, setIsEditing] = useState(false);
  const router = useRouter();

  const handleClose = () => {
    setIsEditing(false);
  };

  const handleLeave = async () => {
    await makeCommand("/leaveRoom", {
      room_id: roomId,
    });
    onLeave();
    router.push("/");
    setIsEditing(false);
  };

  const handleDelete = async () => {
    const result = await makeCommand("/deleteRoom", {
      room_id: roomId,
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

export default EditRoomBtn;
