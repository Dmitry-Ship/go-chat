import React from "react";
import styles from "./Rooms.module.css";
import { Link } from "react-router-dom";
import { Room } from "../types/coreTypes";
import { useRequest } from "../api/hooks";
import NewRoomBtn from "./NewRoomBtn";

function Rooms() {
  const { data, loading } = useRequest<Room[]>("/getRooms");

  return (
    <>
      <header className={`${styles.header} header-for-scrollable`}>
        <h2>Rooms</h2>
      </header>
      <section className="wrap">
        <div className={styles.wrapper}>
          {loading
            ? [{}, {}, {}].map((_, i) => (
                <div key={i} className={styles.room}>
                  <div>
                    <h3>loading...</h3>
                  </div>
                </div>
              ))
            : data?.map((room, i) => (
                <Link key={i} to={"room/" + room.id} className={styles.room}>
                  <div>
                    <h3>{room.name}</h3>
                  </div>
                </Link>
              ))}
        </div>
        <NewRoomBtn />
      </section>
    </>
  );
}

export default Rooms;
