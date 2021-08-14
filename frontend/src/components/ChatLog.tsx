import React, { useEffect, useState } from "react";
import styles from "./ChatLog.module.css";
import { Message } from "../types/coreTypes";

const classes = {
  system: styles.systemMessage,
  outbound: styles.outboundMessage,
  user: "",
};

const ChatLog: React.FC<{ logs: Message[] }> = ({ logs }) => {
  return (
    <div className={styles.log}>
      {logs.map((item, i) => (
        <div key={i} className={`${classes[item.type]} ${styles.message}`}>
          {item.text}
        </div>
      ))}
    </div>
  );
};

export default ChatLog;
