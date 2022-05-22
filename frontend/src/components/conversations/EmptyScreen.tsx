import React from "react";
import NewConversationBtn from "./NewConversationBtn";
import styles from "./EmptyScreen.module.css";

function EmptyScreen() {
  return (
    <div className={styles.wrap}>
      <h3>No conversations yet</h3>
      <NewConversationBtn />
    </div>
  );
}

export default EmptyScreen;
