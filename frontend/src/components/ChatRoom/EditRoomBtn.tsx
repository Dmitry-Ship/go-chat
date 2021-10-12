import React, { useContext, useState } from "react";
import styles from "./EditRoomBtn.module.css";
import SlideIn from "../SlideIn";
import { sendNotification } from "../../api/ws";
import { UserContext } from "../../userContext";
import { useHistory, useParams } from "react-router-dom";
import { makeRequest } from "../../api/fetch";

const EditRoomBtn: React.FC<{ joined: boolean; onLeave: () => void }> = ({
  joined,
  onLeave,
}) => {
  const { roomId } = useParams<{ roomId: string }>();
  const user = useContext(UserContext);
  const [isEditing, setIsEditing] = useState(false);
  const history = useHistory();

  const handleClose = () => {
    setIsEditing(false);
  };

  const handleLeave = () => {
    sendNotification({
      type: "leave",
      data: { room_id: roomId, user_id: user.id },
    });
    onLeave();
    history.push("/");
    setIsEditing(false);
  };

  const handleDelete = async () => {
    await makeRequest("/deleteRoom", {
      method: "POST",
      body: { room_id: roomId, user_id: user.id },
    });

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
