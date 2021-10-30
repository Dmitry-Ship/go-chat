import React, { useState } from "react";
import styles from "./EditRoomBtn.module.css";
import SlideIn from "../common/SlideIn";
import { makeCommand } from "../../api/fetch";
import { useAuth } from "../../contexts/authContext";
import { useWS } from "../../contexts/WSContext";
import { useRouter } from "next/router";

const EditRoomBtn: React.FC<{
  joined: boolean;
  onLeave: () => void;
  roomId: string;
}> = ({ joined, onLeave, roomId }) => {
  const { user } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const router = useRouter();
  const { sendNotification } = useWS();

  const handleClose = () => {
    setIsEditing(false);
  };

  const handleLeave = () => {
    sendNotification("leave", { room_id: roomId, user_id: user?.id });
    onLeave();
    router.push("/");
    setIsEditing(false);
  };

  const handleDelete = async () => {
    await makeCommand("/deleteRoom", { room_id: roomId, user_id: user?.id });

    router.push("/");
    setIsEditing(false);
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
