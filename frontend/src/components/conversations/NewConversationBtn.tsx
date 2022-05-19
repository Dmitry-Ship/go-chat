import React, { useEffect, useRef, useState } from "react";
import styles from "./NewConversationBtn.module.css";
import SlideIn from "../common/SlideIn";
import { v4 as uuidv4 } from "uuid";
import { useRouter } from "next/router";
import { useAPI } from "../../contexts/apiContext";

function NewConversationBtn() {
  const [isCreating, setIsCreating] = useState(false);
  const [conversationName, setConversationName] = useState("");
  const router = useRouter();
  const input = useRef<HTMLInputElement>(null);
  const { makeCommand } = useAPI();

  useEffect(() => {
    if (isCreating) {
      input.current?.focus();
    }
  }, [isCreating]);

  const handleCreate = async () => {
    setConversationName("");
    setIsCreating(false);
    const conversationId = uuidv4();

    const result = await makeCommand("/createConversation", {
      conversation_name: conversationName,
      conversation_id: conversationId,
    });

    if (result.status) {
      router.push(`/conversations/${conversationId}`);
    }
  };

  return (
    <div>
      <button className={"btn"} onClick={() => setIsCreating(true)}>
        + 💬
      </button>
      <SlideIn isOpen={isCreating} onClose={() => setIsCreating(false)}>
        <form className={styles.form} onSubmit={handleCreate}>
          <input
            type="text"
            ref={input}
            placeholder="Conversation name"
            size={32}
            className={`${styles.input} input`}
            value={conversationName}
            onChange={(e) => setConversationName(e.target.value)}
          />
          <button type="submit" disabled={!conversationName} className={"btn"}>
            Create
          </button>
        </form>
      </SlideIn>
    </div>
  );
}

export default NewConversationBtn;
