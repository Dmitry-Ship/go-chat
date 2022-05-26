import React from "react";
import NewConversationBtn from "./NewConversationBtn";
import styles from "./EmptyScreen.module.css";

function EmptyScreen() {
  return (
    <div className={styles.wrap}>
      <h3>No one to talk to yet 🤷🏼</h3>
      <NewConversationBtn text={"+ New Group Chat"} />
    </div>
  );
}

export default EmptyScreen;
