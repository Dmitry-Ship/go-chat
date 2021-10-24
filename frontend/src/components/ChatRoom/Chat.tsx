import React, { useEffect, useState } from "react";
import { Message, MessageRaw, Room } from "../../types/coreTypes";
import styles from "./Chat.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { Link, useHistory, useParams } from "react-router-dom";
import { useRequest } from "../../api/hooks";
import { parseMessage } from "../../messages";
import EditRoomBtn from "./EditRoomBtn";
import { useAuth } from "../../authContext";
import { useWS } from "../../WSContext";

const Chat = () => {
  const { roomId } = useParams<{ roomId: string }>();
  const { user } = useAuth();
  const history = useHistory();

  const [logs, setLogs] = useState<Message[]>([]);
  const [room, setRoom] = useState<Room>();
  const [isJoined, setIsJoined] = useState(false);
  const { sendNotification, subscribe } = useWS();

  const appendLog = (items: Message[]) => {
    setLogs((oldLogs) => [...oldLogs, ...items]);
  };

  const { data: messagesData, loading: messagesLoading } = useRequest<{
    messages: MessageRaw[];
  }>(`/getRoomsMessages?room_id=${roomId}&user_id=${user?.id}`);

  const { data, loading } = useRequest<{
    room: Room;
    joined: boolean;
  }>(`/getRoom?room_id=${roomId}&user_id=${user?.id}`);

  useEffect(() => {
    if (messagesData && !messagesLoading) {
      appendLog(messagesData.messages.map((m) => parseMessage(m)));
    }
  }, [messagesData, messagesLoading]);

  useEffect(() => {
    if (data && !loading) {
      setRoom(data.room);
      setIsJoined(data.joined);
    }
  }, [data, loading]);

  useEffect(() => {
    subscribe("message", (event) => {
      appendLog([parseMessage(event.data)]);
    });

    subscribe("room_deleted", (event) => {
      if (event.data.room_id === roomId) {
        history.push("/");
      }
    });
  }, []);

  const sendMessage = (msg: string, roomId: string, userId: string) => {
    sendNotification({
      type: "message",
      data: { content: msg, room_id: roomId, user_id: userId },
    });
  };

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <Link className={styles.backButton} to="/rooms">
          ⏪
        </Link>
        <b>{room?.name}</b>

        <EditRoomBtn joined={isJoined} onLeave={() => setIsJoined(false)} />
      </header>

      <section className="wrap">
        <ChatLog logs={logs} loading={messagesLoading} />

        <ChatForm
          onSubmit={sendMessage}
          loading={loading}
          joined={isJoined}
          onJoin={() => setIsJoined(true)}
        />
      </section>
    </>
  );
};

export default Chat;
