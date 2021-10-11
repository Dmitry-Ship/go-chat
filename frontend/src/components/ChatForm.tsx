import React, { FormEvent, useContext, useState } from "react";
import styles from "./ChatForm.module.css";
import { useParams } from "react-router-dom";
import { UserContext } from "../userContext";

const ChatForm: React.FC<{
  onSubmit: (message: string, roomId: number, userId: number) => void;
}> = ({ onSubmit }) => {
  const [message, setMessage] = useState<string>("");

  const { roomId } = useParams<{ roomId: string }>();
  const user = useContext(UserContext);

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSubmit(message, Number(roomId), Number(user.id));
    setMessage("");
  };

  return (
    <div className={"controls-for-scrollable"}>
      <form className={styles.form} onSubmit={handleSubmit}>
        <input
          type="text"
          className={styles.input}
          size={64}
          value={message}
          onChange={(e) => setMessage(e.target.value)}
        />
        <button disabled={!message} type="submit" className={styles.submitBtn}>
          ⬆️
        </button>
      </form>
    </div>
  );
};

export default ChatForm;
