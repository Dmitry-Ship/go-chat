import React, { useContext } from "react";
import { makeRequest } from "../api/fetch";
import { UserContext } from "../userContext";
import styles from "./NewRoomBtn.module.css";
import SlideIn from "./SlideIn";
import { useHistory } from "react-router-dom";

function NewRoomBtn() {
  const [isCreating, setIsCreating] = React.useState(false);
  const [roomName, setRoomName] = React.useState("");
  const user = useContext(UserContext);
  const history = useHistory();

  const handleCreate = async () => {
    setRoomName("");
    setIsCreating(false);

    const result = await makeRequest("/createRoom", {
      method: "POST",
      body: { room_name: roomName, user_id: user.id },
    });

    if (result.status) {
      history.push(`/room/${result.data.id}`);
    }
  };

  return (
    <div className="controls-for-scrollable">
      <button className={styles.newRoom} onClick={() => setIsCreating(true)}>
        + New Room
      </button>
      <SlideIn isOpen={isCreating} onClose={() => setIsCreating(false)}>
        <form className={styles.form} onSubmit={handleCreate}>
          <button type="submit" disabled={!roomName} className={styles.newRoom}>
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
    </div>
  );
}

export default NewRoomBtn;
