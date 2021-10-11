import React, { useContext, useEffect, useState } from "react";
import styles from "./EditRoomBtn.module.css";
import SlideIn from "./SlideIn";

const EditRoomBtn = () => {
  const [isEditing, setIsEditing] = useState(false);

  const handleClose = () => {
    setIsEditing(false);
  };

  return (
    <>
      <button onClick={() => setIsEditing(true)} className={styles.editButton}>
        ⚙️
      </button>
      <SlideIn onClose={handleClose} isOpen={isEditing}>
        <div className={styles.menu}>
          <button
            onClick={() => setIsEditing(true)}
            className={styles.menuItem}
          >
            ✏️ Rename
          </button>

          <button
            onClick={() => setIsEditing(true)}
            className={styles.menuItem}
          >
            🗑 Delete
          </button>
        </div>
      </SlideIn>
    </>
  );
};

export default EditRoomBtn;
