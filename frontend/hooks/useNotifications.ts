"use client";

import { useState, useEffect } from "react";
import { useAuth } from "./useAuth";
import { useChat } from "@/contexts/ChatContext";

export const useNotifications = () => {
  const { user } = useAuth();
  const { activeConversationId } = useChat();
  const [permission, setPermission] = useState<NotificationPermission>("default");
  const [enabled, setEnabled] = useState(false);

  useEffect(() => {
    if ("Notification" in window) {
      setPermission(Notification.permission);
      const savedEnabled = localStorage.getItem("notificationsEnabled") === "true";
      setEnabled(savedEnabled);
    }
  }, []);

  const requestPermission = async () => {
    if ("Notification" in window) {
      const result = await Notification.requestPermission();
      setPermission(result);
      return result === "granted";
    }
    return false;
  };

  const showNotification = (title: string, body: string) => {
    if (enabled && permission === "granted" && document.hidden) {
      new Notification(title, {
        body,
        icon: "/favicon.ico",
      });
    }
  };

  const toggleEnabled = async () => {
    if (!enabled) {
      const granted = await requestPermission();
      if (granted) {
        setEnabled(true);
        localStorage.setItem("notificationsEnabled", "true");
      }
    } else {
      setEnabled(false);
      localStorage.setItem("notificationsEnabled", "false");
    }
  };

  return {
    permission,
    enabled,
    requestPermission,
    showNotification,
    toggleEnabled,
  };
};
