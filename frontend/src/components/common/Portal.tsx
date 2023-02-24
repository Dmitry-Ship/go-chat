import React from "react";
import { createPortal } from "react-dom";
import { usePortal } from "./usePortal";

export const Portal: React.FC<{ id: string; children: React.ReactNode }> = ({
  id,
  children,
}) => {
  const target = usePortal(id);
  return createPortal(children, target);
};
