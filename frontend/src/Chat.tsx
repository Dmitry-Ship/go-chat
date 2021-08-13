import React, { useEffect, useState } from "react";
import styles from "./Chat.module.css";

type Message = {
  text: string;
  type: "user" | "system";
};

const Chat = () => {
  const [logs, setLogs] = useState<Message[]>([]);
  const [message, setMessage] = useState<string>("");

  const [conn, setConn] = useState<any>(null);

  const appendLog = (item: Message) => {
    setLogs((oldLogs) => [...oldLogs, item]);
  };

  useEffect(() => {
    // const connection = new WebSocket("ws://" + location.host + "/ws");
    const connection = new WebSocket("ws://localhost:8080/ws");

    connection.onerror = (evt) => {
      console.log("Connection error:", evt);
    };

    connection;
    connection.onclose = (evt) => {
      appendLog({
        text: "Connection closed.",
        type: "system",
      });
    };

    connection.onmessage = (evt) => {
      const messages = evt.data.split("\n");

      messages.forEach((message: string) => {
        appendLog({
          text: message,
          type: "user",
        });
      });
    };

    setConn(connection);
  }, []);

  const handleSubmit = (e: React.MouseEvent<HTMLElement>) => {
    e.preventDefault();

    conn.send(message);

    setMessage("");
  };

  return (
    <div className={styles.wrap}>
      <div className={styles.log}>
        {logs.map((item, i) => (
          <p
            key={i}
            className={item.type === "system" ? styles.systemMessage : ""}
          >
            {item.text}
          </p>
        ))}
      </div>

      <form className={styles.form}>
        <input
          type="text"
          className={styles.input}
          size={64}
          autoFocus
          value={message}
          onChange={(e) => setMessage(e.target.value)}
        />
        <button
          disabled={!conn || conn.readyState !== 1 || !message}
          type="submit"
          className={styles.submitBtn}
          onClick={handleSubmit}
        >
          Send
        </button>
      </form>
    </div>
  );
};

export default Chat;
