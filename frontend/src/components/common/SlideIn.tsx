import React, { useRef } from "react";
import { Portal } from "./Portal";
import styles from "./SlideIn.module.css";

export function SlideIn({
  children,
  isOpen,
  onClose,
}: {
  children: React.ReactNode;
  isOpen: boolean;
  onClose: () => void;
}) {
  const node = useRef(null);

  if (!isOpen) {
    return null;
  }

  function handleClick(e: React.MouseEvent<HTMLDivElement>) {
    // @ts-ignore
    if (!node?.current?.contains(e.target as Node)) {
      onClose();
    }
  }

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
}
