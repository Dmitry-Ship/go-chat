import React, { useEffect, useState } from "react";
import { Room } from "../../types/coreTypes";
import styles from "./Chat.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { Link, useHistory, useParams } from "react-router-dom";
import { useQuery } from "../../api/hooks";
import EditRoomBtn from "./EditRoomBtn";
import { useAuth } from "../../authContext";
import { useWS } from "../../WSContext";

const Chat = () => {
  const { roomId } = useParams<{ roomId: string }>();
  const { user } = useAuth();
  const history = useHistory();
  const [room, setRoom] = useState<Room>();
  const [isJoined, setIsJoined] = useState(false);
  const { subscribe } = useWS();

  const roomQuery = useQuery<{
    room: Room;
    joined: boolean;
  }>(`/getRoom?room_id=${roomId}&user_id=${user?.id}`);

  useEffect(() => {
    if (roomQuery.status === "done" && roomQuery.data) {
      setRoom(roomQuery.data.room);
      setIsJoined(roomQuery.data.joined);
    }
  }, [roomQuery]);

  useEffect(() => {
    subscribe("room_deleted", (event) => {
      if (event.data.room_id === roomId) {
        history.push("/");
      }
    });
  }, []);

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
        <ChatLog />

        <ChatForm
          loading={roomQuery.status === "fetching"}
          joined={isJoined}
          onJoin={() => setIsJoined(true)}
        />
      </section>
    </>
  );
};

export default Chat;
