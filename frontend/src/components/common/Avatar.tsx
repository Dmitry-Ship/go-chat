import React from "react";
import styles from "./Avatar.module.css";

const Avatar: React.FC<{
  src: string;
  size?: number;
}> = ({ src, size }) => {
  return (
    <div className={styles.avatar} style={{ width: size, height: size }}>
      {src}
    </div>
  );
};

export default Avatar;
