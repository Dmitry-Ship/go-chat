import React, { useEffect, useState } from "react";
import { connect, sendMsg } from "../api";
import { Message, Event } from "../types/coreTypes";
import styles from "./Chat.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { UserContext } from "../userContext";

const Chat = () => {
  const [logs, setLogs] = useState<Message[]>([]);
  const [clientId, setClientId] = useState<string | null>(null);

  const appendLog = (items: Message[]) => {
    setLogs((oldLogs) => [...oldLogs, ...items]);
  };

  useEffect(() => {
    connect((events) => {
      events.forEach((event: Event) => {
        switch (event.type) {
          case "message":
            appendLog([
              {
                text: event.data.content,
                type: event.data.type,
                sender: event.data.sender,
                created_at: event.data.created_at,
              },
            ]);
            break;

          case "client_id":
            setClientId(event.data.client_id);
            break;
          default:
            break;
        }
      });
    });

    let vh = window.innerHeight * 0.01;
    document.documentElement.style.setProperty("--vh", `${vh}px`);

    window.addEventListener("resize", () => {
      let vh = window.innerHeight * 0.01;
      document.documentElement.style.setProperty("--vh", `${vh}px`);
    });
  }, []);

  return (
    <UserContext.Provider value={{ id: clientId }}>
      <div className={styles.wrap}>
        <ChatLog logs={logs} />

        <ChatForm onSubmit={sendMsg} />
      </div>
    </UserContext.Provider>
  );
};

export default Chat;
