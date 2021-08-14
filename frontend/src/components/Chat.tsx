import React, { useEffect, useState } from "react";
import { connect, sendMsg } from "../api";
import { Message } from "../types/coreTypes";
import styles from "./Chat.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";

const Chat = () => {
  const [logs, setLogs] = useState<Message[]>([]);
  const [message, setMessage] = useState<string>("");

  const appendLog = (item: Message) => {
    setLogs((oldLogs) => [...oldLogs, item]);
  };

  useEffect(() => {
    connect((msg) => {
      appendLog({
        text: msg.content,
        type: msg.type,
        sender: msg.sender,
      });
    });
  }, []);

  const handleSubmit = () => {
    sendMsg({
      content: message,
      type: "user",
    });

    setMessage("");
  };

  return (
    <div className={styles.wrap}>
      <ChatLog logs={logs} />

      <ChatForm
        message={message}
        onChange={setMessage}
        onSubmit={handleSubmit}
      />
    </div>
  );
};

export default Chat;
