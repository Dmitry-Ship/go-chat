import React, { useEffect, useRef, useState } from "react";
import { makeCommand } from "../../api/fetch";
import styles from "./NewRoomBtn.module.css";
import SlideIn from "../common/SlideIn";
import { v4 as uuidv4 } from "uuid";
import { useAuth } from "../../contexts/authContext";
import { useRouter } from "next/router";

function NewRoomBtn() {
  const [isCreating, setIsCreating] = useState(false);
  const [roomName, setRoomName] = useState("");
  const { user } = useAuth();
  const router = useRouter();
  const input = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isCreating) {
      input.current?.focus();
    }
  }, [isCreating]);

  const handleCreate = async () => {
    setRoomName("");
    setIsCreating(false);
    const roomId = uuidv4();

    const result = await makeCommand("/createRoom", {
      room_name: roomName,
      user_id: user?.id,
      room_id: roomId,
    });

    if (result.status) {
      router.push(`/rooms/${roomId}`);
    }
  };

  return (
    <div>
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
