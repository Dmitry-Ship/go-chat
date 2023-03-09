import React from "react";
import { createPortal } from "react-dom";
import { usePortal } from "./usePortal";

export function Portal({
  id,
  children,
}: {
  id: string;
  children: React.ReactNode;
}) {
  const target = usePortal(id);
  return createPortal(children, target);
}
