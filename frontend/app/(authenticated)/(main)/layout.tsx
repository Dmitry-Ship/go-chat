"use client";
import { usePathname } from "next/navigation";
import { Navigation } from "./navigation";

export default function NavLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  return (
    <>
      {children}
      <Navigation current={pathname || ""} />
    </>
  );
}
