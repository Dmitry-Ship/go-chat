import React, { useContext, useEffect, useState } from "react";
import { onEvent, sendMsg } from "../../api/ws";
import { Message, MessageRaw, Room } from "../../types/coreTypes";
import styles from "./Chat.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { Link, useHistory, useParams } from "react-router-dom";
import { useRequest } from "../../api/hooks";
import { UserContext } from "../../userContext";
import { parseMessage } from "../../messages";
import EditRoomBtn from "./EditRoomBtn";

const Chat = () => {
  const { roomId } = useParams<{ roomId: string }>();
  const user = useContext(UserContext);
  const history = useHistory();

  const [logs, setLogs] = useState<Message[]>([]);
  const [room, setRoom] = useState<Room>();
  const [isJoined, setIsJoined] = useState(false);

  const appendLog = (items: Message[]) => {
    setLogs((oldLogs) => [...oldLogs, ...items]);
  };

  const { data, loading } = useRequest<{
    room: Room;
    messages: MessageRaw[];
    joined: boolean;
  }>(`/getRoomsMessages?room_id=${roomId}&user_id=${user.id}`);

  useEffect(() => {
    if (data && !loading) {
      appendLog(data.messages.map((m) => parseMessage(m)));
      setRoom(data.room);
      setIsJoined(data.joined);
    }
  }, [data, loading]);

  useEffect(() => {
    onEvent("message", (event) => {
      appendLog([parseMessage(event.data)]);
    });

    onEvent("room_deleted", (event) => {
      if (event.data.room_id === Number(roomId)) {
        history.push("/");
      }
    });
  }, []);

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <Link className={styles.backButton} to="/">
          ⏪
        </Link>
        <b>{room?.name}</b>

        <EditRoomBtn joined={isJoined} onLeave={() => setIsJoined(false)} />
      </header>
      <section className="wrap">
        <ChatLog logs={logs} loading={loading} />

        <ChatForm
          onSubmit={sendMsg}
          loading={loading}
          joined={isJoined}
          onJoin={() => setIsJoined(true)}
        />
      </section>
    </>
  );
};

export default Chat;
