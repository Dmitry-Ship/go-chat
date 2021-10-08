import React, { useEffect, useState } from "react";
import { onEvent, sendMsg, sendNotification } from "../api/ws";
import { Message, Room } from "../types/coreTypes";
import styles from "./Chat.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { Link, useParams } from "react-router-dom";
import { useRequest } from "../api/hooks";

const Chat = () => {
  const [logs, setLogs] = useState<Message[]>([]);

  const appendLog = (items: Message[]) => {
    setLogs((oldLogs) => [...oldLogs, ...items]);
  };

  const { roomId } = useParams<{ roomId: string }>();

  const { data, loading } = useRequest<{ room: Room; messages: Message[] }>(
    "/getRoomsMessages?room_id=" + roomId
  );

  useEffect(() => {
    onEvent("message", (event) => {
      appendLog([
        {
          text: event.data.content,
          type: event.data.type,
          user: event.data.user,
          created_at: event.data.created_at,
          roomId: event.data.room_id,
        },
      ]);
    });

    sendNotification({ type: "join", data: { room_id: Number(roomId) } });

    let vh = window.innerHeight * 0.01;
    document.documentElement.style.setProperty("--vh", `${vh}px`);

    window.addEventListener("resize", () => {
      let vh = window.innerHeight * 0.01;
      document.documentElement.style.setProperty("--vh", `${vh}px`);
    });
    return () => {
      sendNotification({ type: "leave", data: { room_id: Number(roomId) } });
    };
  }, []);

  return (
    <>
      <div className={styles.header}>
        <Link className={styles.backButton} to="/">
          ⏪
        </Link>
        <b>{data?.room?.name}</b>
      </div>

      <ChatLog logs={logs} />

      <ChatForm onSubmit={sendMsg} />
    </>
  );
};

export default Chat;
