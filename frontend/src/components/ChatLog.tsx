import React from "react";
import styles from "./ChatLog.module.css";
import { Message } from "../types/coreTypes";
import MessageComponent from "./Message";

const ChatLog: React.FC<{ logs: Message[] }> = ({ logs }) => {
  return (
    <div className={styles.log}>
      {logs.map((item, i) => (
        <MessageComponent key={i} message={item} />
      ))}
    </div>
  );
};

export default ChatLog;
