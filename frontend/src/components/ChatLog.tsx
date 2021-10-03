import React from "react";
import styles from "./ChatLog.module.css";
import { Message } from "../types/coreTypes";
import MessageComponent from "./Message";

const ChatLog: React.FC<{ logs: Message[] }> = ({ logs }) => {
  return (
    <div className={styles.log}>
      {logs.map((item, i) => {
        const previous = logs[i - 1];
        const isFistInAGroup =
          !previous ||
          previous?.type === "system" ||
          item.sender.id !== previous?.sender.id;

        const next = logs[i + 1];

        const isLastInAGroup =
          !next ||
          next?.type === "system" ||
          item.sender.id !== next?.sender.id;

        return (
          <MessageComponent
            key={i}
            message={item}
            isFistInAGroup={isFistInAGroup}
            isLastInAGroup={isLastInAGroup}
          />
        );
      })}
    </div>
  );
};

export default ChatLog;
