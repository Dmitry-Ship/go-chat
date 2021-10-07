import React, { useEffect } from "react";
import styles from "./Rooms.module.css";
import { Link } from "react-router-dom";
import { Room } from "../types/coreTypes";
import { makeRequest } from "../api/fetch";

function Rooms() {
  const [rooms, setRooms] = React.useState<Room[]>([]);

  useEffect(() => {
    const fetchRooms = async () => {
      const rooms = await makeRequest("/getRooms", "GET", null);
      setRooms(rooms);
    };

    fetchRooms();
  }, []);

  return (
    <Link to="room/123">
      {rooms.map((room, i) => (
        <div key={i} className={styles.room}>
          {room.name}
        </div>
      ))}
    </Link>
  );
}

export default Rooms;
