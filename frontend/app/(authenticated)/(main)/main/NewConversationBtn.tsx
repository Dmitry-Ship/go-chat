import React, { useRef, useState } from "react";
import styles from "./NewConversationBtn.module.css";
import { SlideIn } from "../../../../src/components/common/SlideIn";
import { v4 as uuidv4 } from "uuid";
import { useRouter } from "next/navigation";
import { useMutation } from "react-query";
import { createConversation } from "../../../../src/api/fetch";

export const NewConversationBtn: React.FC<{ text: string }> = ({ text }) => {
  const [isCreating, setIsCreating] = useState(false);
  const [conversationName, setConversationName] = useState("");
  const [newId, setNewId] = useState("");

  const router = useRouter();
  const input = useRef<HTMLInputElement>(null);

  const { mutate } = useMutation(createConversation, {
    onSuccess: (data) => {
      router.push(`/conversations/${newId}`);
      setConversationName("");
      setNewId("");
    },
  });

  const handleCreate = () => {
    const conversationId = uuidv4();
    setNewId(conversationId);

    mutate({
      conversation_name: conversationName,
      conversation_id: conversationId,
    });
  };

  const handleStartCreating = () => {
    setIsCreating(true);
    input.current?.focus();
  };

  return (
    <div>
      <button className={"btn"} onClick={handleStartCreating}>
        {text}
      </button>
      <SlideIn isOpen={isCreating} onClose={() => setIsCreating(false)}>
        <form className={styles.form} onSubmit={handleCreate}>
          <input
            type="text"
            ref={input}
            placeholder="Group name"
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
};
