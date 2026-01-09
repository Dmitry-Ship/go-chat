"use client";

import { useState } from "react";

export const useNotifications = () => {
  const [permission, setPermission] = useState<NotificationPermission>(
    typeof window !== "undefined" && "Notification" in window
      ? Notification.permission
      : "default"
  );
  const [enabled, setEnabled] = useState(
    typeof window !== "undefined"
      ? localStorage.getItem("notificationsEnabled") === "true"
      : false
  );

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
