import React from "react";
import { createPortal } from "react-dom";
import { usePortal } from "./usePortal";

export const Portal = ({
  id,
  children,
}: {
  id: string;
  children: React.ReactNode;
}) => {
  const target = usePortal(id);
  return createPortal(children, target);
};
