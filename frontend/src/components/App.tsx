import React from "react";
import styles from "./App.module.css";
import { ProvideAuth } from "../authContext";
import Routes from "./Routes";

function App() {
  return (
    <ProvideAuth>
      <div className={styles.app}>
        <Routes />
      </div>
    </ProvideAuth>
  );
}

export default App;
