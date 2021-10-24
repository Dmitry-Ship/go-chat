import React from "react";
import styles from "./Avatar.module.css";

const Avatar: React.FC<{
  src: string;
}> = ({ src }) => {
  return <div className={styles.avatar}>{src}</div>;
};

export default Avatar;
