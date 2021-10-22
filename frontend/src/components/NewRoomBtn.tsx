import React, { useEffect } from "react";
import { makeRequest } from "../api/fetch";
import styles from "./NewRoomBtn.module.css";
import SlideIn from "./SlideIn";
import { useHistory } from "react-router-dom";
import { v4 as uuidv4 } from "uuid";
import { useAuth } from "../authContext";

function NewRoomBtn() {
  const [isCreating, setIsCreating] = React.useState(false);
  const [roomName, setRoomName] = React.useState("");
  const user = useAuth().user;
  const history = useHistory();
  const input = React.useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isCreating) {
      input.current?.focus();
    }
  }, [isCreating]);

  const handleCreate = async () => {
    setRoomName("");
    setIsCreating(false);
    const roomId = uuidv4();

    const result = await makeRequest("/createRoom", {
      method: "POST",
      body: { room_name: roomName, user_id: user?.id, room_id: roomId },
    });

    if (result.status) {
      history.push(`/room/${roomId}`);
    }
  };

  return (
    <div className={`${styles.wrap} controls-for-scrollable`}>
      <button className={"btn"} onClick={() => setIsCreating(true)}>
        + New Room
      </button>
      <SlideIn isOpen={isCreating} onClose={() => setIsCreating(false)}>
        <form className={styles.form} onSubmit={handleCreate}>
          <input
            type="text"
            ref={input}
            placeholder="Room name"
            size={32}
            className={`${styles.input} input`}
            value={roomName}
            onChange={(e) => setRoomName(e.target.value)}
          />
          <button type="submit" disabled={!roomName} className={"btn"}>
            Create
          </button>
        </form>
      </SlideIn>
    </div>
  );
}

export default NewRoomBtn;
