import React from "react";
import styles from "./Avatar.module.css";

const sizesMap = {
  md: 65,
  lg: 100,
};

const Avatar: React.FC<{
  src: string;
  size?: keyof typeof sizesMap;
}> = ({ src, size = "md" }) => {
  return (
    <div className={styles.avatar} style={{ width: size, height: size }}>
      {src}
    </div>
  );
};

export default Avatar;
