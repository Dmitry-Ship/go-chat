import React from "react";
import styles from "./App.module.css";
import Chat from "./Chat";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import Rooms from "./Rooms";

function App() {
  return (
    <div className={styles.App}>
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
    </div>
  );
}

export default App;
