import React, { useContext } from "react";
import { makeRequest } from "../api/fetch";
import { UserContext } from "../userContext";
import styles from "./NewRoomBtn.module.css";
import SlideIn from "./SlideIn";

function NewRoomBtn() {
  const [isCreating, setIsCreating] = React.useState(false);
  const [roomName, setRoomName] = React.useState("");
  const user = useContext(UserContext);

  const handleCreate = async () => {
    setRoomName("");
    setIsCreating(false);

    const result = await makeRequest("/createRoom", {
      method: "POST",
      body: { room_name: roomName, user_id: user.id },
    });

    if (result.status) {
      window.location.href = `/room/${result.data.id}`;
    }
  };

  return (
    <>
      <button className={styles.newRoom} onClick={() => setIsCreating(true)}>
        + New Room
      </button>
      <SlideIn isOpen={isCreating} onClose={() => setIsCreating(false)}>
        <form className={styles.form}>
          <button
            type="submit"
            disabled={!roomName}
            className={styles.newRoom}
            onClick={handleCreate}
          >
            Create
          </button>
          <input
            type="text"
            placeholder="Room name"
            autoFocus
            size={32}
            className={styles.input}
            value={roomName}
            onChange={(e) => setRoomName(e.target.value)}
          />
        </form>
      </SlideIn>
    </>
  );
}

export default NewRoomBtn;
