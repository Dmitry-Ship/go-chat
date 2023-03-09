import React from "react";
import styles from "./Avatar.module.css";

const sizesMap = {
  md: 65,
  lg: 100,
};

export function Avatar({
  src,
  size = "md",
}: {
  src: string;
  size?: keyof typeof sizesMap;
}) {
  return (
    <div className={styles.avatar} style={{ width: size, height: size }}>
      {src}
    </div>
  );
}
