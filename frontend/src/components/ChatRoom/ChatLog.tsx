import React, { useEffect, useRef } from "react";
import styles from "./ChatLog.module.css";
import { Message } from "../../types/coreTypes";
import MessageComponent from "./Message";
import Loader from "../common/Loader";

const ChatLog: React.FC<{ logs: Message[]; loading: boolean }> = ({
  logs,
  loading,
}) => {
  const logComponent = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (logs.length > 0) {
      logComponent.current?.scrollIntoView();
    }
  }, [logs]);

  return (
    <main className={`${styles.log} scrollable-content`}>
      {loading ? (
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
