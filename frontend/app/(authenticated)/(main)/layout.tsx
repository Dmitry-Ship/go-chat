"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import styles from "./layout.module.css";

const links = [
  { href: "/main", label: "💬" },
  { href: "/contacts", label: "👥" },
  { href: "/settings", label: "⚙️" },
];

export default function NavLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  return (
    <>
      {children}
      <div className="controls-for-scrollable">
        {links.map((link) => (
          <Link
            href={link.href}
            key={link.href}
            className={`${styles.navBtn} shadow ${
              pathname === link.href ? styles.active : ""
            }`}
          >
            {link.label}
          </Link>
        ))}
      </div>
    </>
  );
}
