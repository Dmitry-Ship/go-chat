import React, { useEffect, useRef } from "react";
import styles from "./SlideIn.module.css";

const SlideIn: React.FC<{
  children: React.ReactNode;
  isOpen: boolean;
  onClose: () => void;
}> = ({ children, isOpen, onClose }) => {
  const node = useRef(null);

  if (!isOpen) {
    return null;
  }

  const handleClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!node?.current?.contains(e.target as Node)) {
      onClose();
    }
  };

  return (
    <div className={styles.overlay} onClick={handleClick}>
      <div ref={node} className={styles.slideIn}>
        <div onClick={onClose}>❌</div>
        {children}
      </div>
    </div>
  );
};

export default SlideIn;
