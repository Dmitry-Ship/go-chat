import React from "react";
import styles from "./EmptyScreen.module.css";

export function EmptyScreen({
  text,
  children,
}: {
  text: string;
  children?: React.ReactNode;
}) {
  return (
    <div className={styles.wrap}>
      <h3>{text}</h3>

      {children}
    </div>
  );
}
