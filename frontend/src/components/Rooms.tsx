import React from "react";
import styles from "./Rooms.module.css";
import { Link } from "react-router-dom";
import { Room } from "../types/coreTypes";
import { useQuery } from "../api/hooks";
import NewRoomBtn from "./NewRoomBtn";
import Loader from "./common/Loader";

function Rooms() {
  const response = useQuery<Room[]>("/getRooms");

  return (
    <>
      <header className={`header header-for-scrollable`}>
        <h2>Rooms</h2>
        <NewRoomBtn />
      </header>
      <section className="wrap">
        <div className={`${styles.wrapper} scrollable-content`}>
          {(() => {
            switch (response.status) {
              case "fetching":
                return <Loader />;
              case "done":
                return response.data?.map((room, i) => (
                  <Link
                    key={i}
                    to={"room/" + room.id}
                    className={`${styles.room} rounded`}
                  >
                    <div>
                      <h3>{room.name}</h3>
                    </div>
                  </Link>
                ));
              default:
                return null;
            }
          })()}
        </div>
      </section>
    </>
  );
}

export default Rooms;
