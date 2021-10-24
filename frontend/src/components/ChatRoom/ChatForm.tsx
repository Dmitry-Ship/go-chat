import React, { FormEvent, useState } from "react";
import styles from "./ChatForm.module.css";
import { useParams } from "react-router-dom";
import Loader from "../common/Loader";
import { useAuth } from "../../authContext";
import { useWS } from "../../WSContext";

const ChatForm: React.FC<{
  loading: boolean;
  joined: boolean;
  onJoin: () => void;
  onSubmit: (message: string, roomId: string, userId: string) => void;
}> = ({ onSubmit, loading, joined, onJoin }) => {
  const [message, setMessage] = useState<string>("");

  const { roomId } = useParams<{ roomId: string }>();
  const auth = useAuth();
  const { sendNotification } = useWS();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSubmit(message, roomId, auth.user?.id || "");
    setMessage("");
  };

  const handleJoin = () => {
    sendNotification({
      type: "join",
      data: { room_id: roomId, user_id: auth.user?.id },
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
