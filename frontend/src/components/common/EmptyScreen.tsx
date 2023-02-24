import React from "react";
import styles from "./EmptyScreen.module.css";

export const EmptyScreen: React.FC<{
  text: string;
  children?: React.ReactNode;
}> = ({ text, children }) => {
  return (
    <div className={styles.wrap}>
      <h3>{text}</h3>

      {children}
    </div>
  );
};
