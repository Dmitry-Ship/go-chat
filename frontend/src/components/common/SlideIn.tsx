import React, { useRef } from "react";
import { Portal } from "./Portal";
import styles from "./SlideIn.module.css";

export const SlideIn: React.FC<{
  children: React.ReactNode;
  isOpen: boolean;
  onClose: () => void;
}> = ({ children, isOpen, onClose }) => {
  const node = useRef(null);

  if (!isOpen) {
    return null;
  }

  const handleClick = (e: React.MouseEvent<HTMLDivElement>) => {
    // @ts-ignore
    if (!node?.current?.contains(e.target as Node)) {
      onClose();
    }
  };

  return (
    <Portal id="modal">
      <div className={styles.overlay} onClick={handleClick}>
        <div ref={node} className={styles.slideIn}>
          <button className={styles.closeBtn} onClick={onClose}>
            ❌
          </button>
          {children}
        </div>
      </div>
    </Portal>
  );
};
