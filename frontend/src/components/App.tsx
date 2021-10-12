import React, { useEffect, useState } from "react";
import styles from "./App.module.css";
import Chat from "./ChatRoom/Chat";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import Rooms from "./Rooms";
import { UserContext } from "../userContext";
import { onEvent } from "../api/ws";
import { useWS } from "../api/hooks";
import Loader from "./Loader";

function App() {
  const { status } = useWS();
  const [userId, setUserId] = useState<string | null>(null);

  useEffect(() => {
    onEvent("user_id", (event) => {
      setUserId(event.data.user_id);
    });
  }, []);

  return (
    <UserContext.Provider value={{ id: userId }}>
      <div className={styles.app}>
        {status === "connecting" ? (
          <Loader />
        ) : (
          <Router>
            <Switch>
              <Route path="/room/:roomId">
                <Chat />
              </Route>
              <Route path="/">
                <Rooms />
              </Route>
            </Switch>
          </Router>
        )}
      </div>
    </UserContext.Provider>
  );
}

export default App;
