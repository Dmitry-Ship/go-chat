import React from "react";
import { useAPI } from "../src/contexts/apiContext";
import styles from "./ErrorAlert.module.css";

const ErrorAlert: React.FC = () => {
  const { error, setError } = useAPI();

  if (!error) {
    return null;
  }

  return (
    <div className={styles.wrap}>
      <span>{error}</span>
      <button onClick={() => setError(null)}>⤫</button>
    </div>
  );
};

export default ErrorAlert;
