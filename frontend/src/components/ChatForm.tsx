import React from "react";
import styles from "./ChatForm.module.css";

const ChatForm: React.FC<{
  message: string;
  onChange: (value: string) => void;
  onSubmit: (e: React.MouseEvent<HTMLElement>) => void;
}> = ({ message, onChange, onSubmit }) => {
  return (
    <form className={styles.form}>
      <input
        type="text"
        className={styles.input}
        size={64}
        autoFocus
        value={message}
        onChange={(e) => onChange(e.target.value)}
      />
      <button
        disabled={!message}
        type="submit"
        className={styles.submitBtn}
        onClick={onSubmit}
      >
        🚀
      </button>
    </form>
  );
};

export default ChatForm;
