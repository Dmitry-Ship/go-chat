import React, { useEffect, useState } from "react";
import { connect, sendMsg } from "../api";
import { Message } from "../types/coreTypes";
import styles from "./Chat.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";

const Chat = () => {
  const [logs, setLogs] = useState<Message[]>([]);
  const [message, setMessage] = useState<string>("");

  const appendLog = (items: Message[]) => {
    setLogs((oldLogs) => [...oldLogs, ...items]);
  };

  useEffect(() => {
    connect((msgs) => {
      const messages = msgs.map((msg: any) => ({
        text: msg.content,
        type: msg.type,
        sender: msg.sender,
        created_at: msg.created_at,
      }));
      appendLog(messages);
    });
  }, []);

  const handleSubmit = (e: React.MouseEvent<HTMLElement>) => {
    e.preventDefault();

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
