import React, { useState } from "react";
import styles from "./ChatForm.module.css";

const ChatForm: React.FC<{
  onSubmit: (message: string, roomId: number) => void;
}> = ({ onSubmit }) => {
  const [message, setMessage] = useState<string>("");

  const handleSubmit = (e: React.MouseEvent<HTMLElement>) => {
    e.preventDefault();
    onSubmit(message, 123);
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
