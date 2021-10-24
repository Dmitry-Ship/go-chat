import React, { useState } from "react";
import styles from "./EditRoomBtn.module.css";
import SlideIn from "../common/SlideIn";
import { useHistory, useParams } from "react-router-dom";
import { makeCommand } from "../../api/fetch";
import { useAuth } from "../../authContext";
import { useWS } from "../../WSContext";

const EditRoomBtn: React.FC<{ joined: boolean; onLeave: () => void }> = ({
  joined,
  onLeave,
}) => {
  const { roomId } = useParams<{ roomId: string }>();
  const { user } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const history = useHistory();
  const { sendNotification } = useWS();

  const handleClose = () => {
    setIsEditing(false);
  };

  const handleLeave = () => {
    sendNotification("leave", { room_id: roomId, user_id: user?.id });
    onLeave();
    history.push("/");
    setIsEditing(false);
  };

  const handleDelete = async () => {
    await makeCommand("/deleteRoom", { room_id: roomId, user_id: user?.id });

    history.push("/");
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
