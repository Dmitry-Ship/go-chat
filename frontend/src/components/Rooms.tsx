import React from "react";
import styles from "./Rooms.module.css";
import { Link } from "react-router-dom";
import { Room } from "../types/coreTypes";
import { useRequest } from "../api/hooks";
import NewRoomBtn from "./NewRoomBtn";
import Loader from "./common/Loader";
import AccountSettingsBtn from "./AccountSettingsBtn";

function Rooms() {
  const { data, loading } = useRequest<Room[]>("/getRooms");

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <h2>Rooms</h2>
        <NewRoomBtn />
      </header>
      <section className="wrap">
        <div className={`${styles.wrapper} scrollable-content`}>
          {loading ? (
            <Loader />
          ) : (
            data?.map((room, i) => (
              <Link
                key={i}
                to={"room/" + room.id}
                className={`${styles.room} rounded`}
              >
                <div>
                  <h3>{room.name}</h3>
                </div>
              </Link>
            ))
          )}
        </div>
      </section>
    </>
  );
}

export default Rooms;
