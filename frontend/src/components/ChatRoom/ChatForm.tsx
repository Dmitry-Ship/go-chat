import React, { FormEvent, useContext, useState } from "react";
import styles from "./ChatForm.module.css";
import { useParams } from "react-router-dom";
import { UserContext } from "../../userContext";
import Loader from "../Loader";
import { sendNotification } from "../../api/ws";

const ChatForm: React.FC<{
  loading: boolean;
  joined: boolean;
  onJoin: () => void;
  onSubmit: (message: string, roomId: number, userId: number) => void;
}> = ({ onSubmit, loading, joined, onJoin }) => {
  const [message, setMessage] = useState<string>("");

  const { roomId } = useParams<{ roomId: string }>();
  const user = useContext(UserContext);

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSubmit(message, Number(roomId), Number(user.id));
    setMessage("");
  };

  const handleJoin = () => {
    sendNotification({
      type: "join",
      data: { room_id: Number(roomId), user_id: user.id },
    });
    onJoin();
  };

  return (
    <div className={"controls-for-scrollable"}>
      {loading ? (
        <Loader />
      ) : (
        <>
          {joined ? (
            <form className={styles.form} onSubmit={handleSubmit}>
              <input
                type="text"
                className={"input"}
                size={64}
                value={message}
                onChange={(e) => setMessage(e.target.value)}
              />
              <button
                disabled={!message}
                type="submit"
                className={styles.submitBtn}
              >
                ⬆️
              </button>
            </form>
          ) : (
            <button onClick={handleJoin} className="btn">
              Join
            </button>
          )}
        </>
      )}
    </div>
  );
};

export default ChatForm;
