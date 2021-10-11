import React, { useContext, useEffect, useRef, useState } from "react";
import { onEvent, sendMsg, sendNotification } from "../api/ws";
import { Message, MessageRaw, Room } from "../types/coreTypes";
import styles from "./Chat.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { Link, useParams } from "react-router-dom";
import { useRequest } from "../api/hooks";
import { UserContext } from "../userContext";
import { parseMessage } from "../messages";
import EditRoomBtn from "./EditRoomBtn";

const Chat = () => {
  const { roomId } = useParams<{ roomId: string }>();

  const [logs, setLogs] = useState<Message[]>([]);

  const appendLog = (items: Message[]) => {
    setLogs((oldLogs) => [...oldLogs, ...items]);
  };

  const { data, loading } = useRequest<{ room: Room; messages: MessageRaw[] }>(
    "/getRoomsMessages?room_id=" + roomId
  );

  useEffect(() => {
    if (data && !loading) {
      appendLog(data.messages.map((m) => parseMessage(m)));
    }
  }, [data, loading]);

  const user = useContext(UserContext);

  useEffect(() => {
    onEvent("message", (event) => {
      appendLog([parseMessage(event.data)]);
    });

    sendNotification({
      type: "join",
      data: { room_id: Number(roomId), user_id: user.id },
    });

    return () => {
      sendNotification({
        type: "leave",
        data: { room_id: Number(roomId), user_id: user.id },
      });
    };
  }, []);

  const handleSubmit = (msg: string, roomId: number, userId: number) => {
    sendMsg(msg, roomId, userId);
  };

  return (
    <>
      <header className={`${styles.header} header-for-scrollable`}>
        <Link className={styles.backButton} to="/">
          ⏪
        </Link>
        <b>{data?.room?.name}</b>

        <EditRoomBtn />
      </header>
      <section className="wrap">
        <ChatLog logs={logs} />
        <ChatForm onSubmit={handleSubmit} />
      </section>
    </>
  );
};

export default Chat;
