import React, { useEffect, useState } from "react";
import { connect, sendMsg } from "../api/ws";
import { Message, Event } from "../types/coreTypes";
import styles from "./Chat.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { UserContext } from "../userContext";
import { Link } from "react-router-dom";

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
                user: event.data.user,
                created_at: event.data.created_at,
                roomId: event.data.room_id,
              },
            ]);
            break;

          case "user_id":
            setClientId(event.data.user_id);
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
      <Link to="/">leave</Link>

      <div className={styles.wrap}>
        <ChatLog logs={logs} />

        <ChatForm onSubmit={sendMsg} />
      </div>
    </UserContext.Provider>
  );
};

export default Chat;
