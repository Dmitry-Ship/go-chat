import React, { useEffect, useRef, useState } from "react";
import styles from "./ChatLog.module.css";
import { Message, MessageRaw } from "../../types/coreTypes";
import MessageComponent from "./MessageComponent";
import Loader from "../common/Loader";
import { useQuery } from "../../api/hooks";
import { parseMessage } from "../../messages";
import { useWS } from "../../contexts/WSContext";

const ChatLog: React.FC<{ roomId: string }> = ({ roomId }) => {
  const { onNotification } = useWS();

  const [logs, setLogs] = useState<Message[]>([]);

  const appendLog = (items: Message[]) => {
    setLogs((oldLogs) => [...oldLogs, ...items]);
  };

  const messagesQuery = useQuery<{
    messages: MessageRaw[];
  }>(`/getRoomsMessages?room_id=${roomId}`);

  useEffect(() => {
    if (messagesQuery.status === "done" && messagesQuery.data) {
      appendLog(messagesQuery.data.messages?.map((m) => parseMessage(m)));
    }
  }, [messagesQuery]);

  const logComponent = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (logs.length > 0) {
      logComponent.current?.scrollIntoView();
    }
  }, [logs]);

  useEffect(() => {
    onNotification("message", (event) => {
      appendLog([parseMessage(event.data)]);
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <main className={`${styles.log} scrollable-content`}>
      {messagesQuery.status === "fetching" ? (
        <Loader />
      ) : (
        <>
          {logs.map((item, i) => {
            const previous = logs[i - 1];
            const isFistInAGroup =
              !previous ||
              previous?.type === "system" ||
              item.user.id !== previous?.user.id;

            const next = logs[i + 1];

            const isLastInAGroup =
              !next ||
              next?.type === "system" ||
              item.user.id !== next?.user.id;

            return (
              <MessageComponent
                key={i}
                message={item}
                isFistInAGroup={isFistInAGroup}
                isLastInAGroup={isLastInAGroup}
              />
            );
          })}
          <div ref={logComponent} />
        </>
      )}
    </main>
  );
};

export default ChatLog;