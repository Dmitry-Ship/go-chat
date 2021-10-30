import React, { useEffect, useState } from "react";
import { Room } from "../../types/coreTypes";
import styles from "./ChatRoom.module.css";
import ChatForm from "./ChatForm";
import ChatLog from "./ChatLog";
import { useQuery } from "../../api/hooks";
import EditRoomBtn from "./EditRoomBtn";
import { useAuth } from "../../contexts/authContext";
import { useWS } from "../../contexts/WSContext";
import { useRouter } from "next/router";
import Link from "next/link";

const ChatRoom: React.FC = () => {
  const { user } = useAuth();
  const router = useRouter();
  const roomId = router.query.roomId as string;
  const [room, setRoom] = useState<Room>();
  const [isJoined, setIsJoined] = useState(false);
  const { subscribe } = useWS();

  useEffect(() => {
    subscribe("room_deleted", (event) => {
      if (event.data.room_id === roomId) {
        router.push("/");
      }
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [roomId, router]);

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

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <Link href="/">
          <a className={styles.backButton}>⏪</a>
        </Link>
        <b>{room?.name}</b>

        <EditRoomBtn
          roomId={roomId}
          joined={isJoined}
          onLeave={() => setIsJoined(false)}
        />
      </header>

      <section className="wrap">
        <ChatLog roomId={roomId} />

        <ChatForm
          roomId={roomId}
          loading={roomQuery.status === "fetching"}
          joined={isJoined}
          onJoin={() => setIsJoined(true)}
        />
      </section>
    </>
  );
};

export default ChatRoom;
