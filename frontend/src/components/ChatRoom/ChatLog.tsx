import React from "react";
import styles from "./ChatLog.module.css";
import { Message } from "../../types/coreTypes";
import MessageComponent from "./Message";

const ChatLog: React.FC<{ logs: Message[] }> = ({ logs }) => {
  return (
    <main className={`${styles.log} scrollable-content`}>
      {logs.map((item, i) => {
        const previous = logs[i - 1];
        const isFistInAGroup =
          !previous ||
          previous?.type === "system" ||
          item.user.id !== previous?.user.id;

        const next = logs[i + 1];

        const isLastInAGroup =
          !next || next?.type === "system" || item.user.id !== next?.user.id;

        return (
          <MessageComponent
            key={i}
            message={item}
            isFistInAGroup={isFistInAGroup}
            isLastInAGroup={isLastInAGroup}
          />
        );
      })}
    </main>
  );
};

export default ChatLog;
