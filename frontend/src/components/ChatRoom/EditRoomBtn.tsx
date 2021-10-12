import React, { useContext, useState } from "react";
import styles from "./EditRoomBtn.module.css";
import SlideIn from "../SlideIn";
import { sendNotification } from "../../api/ws";
import { UserContext } from "../../userContext";
import { useParams } from "react-router-dom";

const EditRoomBtn: React.FC<{ joined: boolean; onLeave: () => void }> = ({
  joined,
  onLeave,
}) => {
  const { roomId } = useParams<{ roomId: string }>();
  const user = useContext(UserContext);
  const [isEditing, setIsEditing] = useState(false);

  const handleClose = () => {
    setIsEditing(false);
  };

  const handleLeave = () => {
    sendNotification({
      type: "leave",
      data: { room_id: Number(roomId), user_id: user.id },
    });
    setIsEditing(false);
    onLeave();
  };

  return (
    <>
      <button onClick={() => setIsEditing(true)} className={styles.editButton}>
        ⚙️
      </button>
      <SlideIn onClose={handleClose} isOpen={isEditing}>
        <>
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
