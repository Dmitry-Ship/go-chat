import React from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import styles from "./Navigation.module.css";

const links = [
  { href: "/", label: "💬" },
  { href: "/contacts", label: "👥" },
  { href: "/settings", label: "⚙️" },
];

const Navigation = () => {
  const router = useRouter();

  return (
    <div className="controls-for-scrollable">
      {links.map((link) => (
        <Link
          href={link.href}
          key={link.href}
          className={`${styles.navBtn} shadow ${
            router.pathname === link.href ? styles.active : ""
          }`}
        >
          {link.label}
        </Link>
      ))}
    </div>
  );
};

export default Navigation;
