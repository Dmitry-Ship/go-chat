import React, { useContext, useState } from "react";
import styles from "./ChatForm.module.css";
import { useParams } from "react-router-dom";
import { UserContext } from "../userContext";

const ChatForm: React.FC<{
  onSubmit: (message: string, roomId: number, userId: number) => void;
}> = ({ onSubmit }) => {
  const [message, setMessage] = useState<string>("");

  const { roomId } = useParams<{ roomId: string }>();
  const user = useContext(UserContext);

  const handleSubmit = (e: React.MouseEvent<HTMLElement>) => {
    e.preventDefault();
    onSubmit(message, Number(roomId), Number(user.id));
    setMessage("");
  };

  return (
    <form className={styles.form}>
      <input
        type="text"
        className={styles.input}
        size={64}
        autoFocus
        value={message}
        onChange={(e) => setMessage(e.target.value)}
      />
      <button
        disabled={!message}
        type="submit"
        className={styles.submitBtn}
        onClick={handleSubmit}
      >
        ⬆️
      </button>
    </form>
  );
};

export default ChatForm;
