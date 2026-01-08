import React from "react";
import styles from "./Avatar.module.css";
import { AVATAR_SIZES } from "../../constants";

const sizesMap = AVATAR_SIZES;

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
