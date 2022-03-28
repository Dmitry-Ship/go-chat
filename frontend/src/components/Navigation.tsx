import React from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import styles from "./Navigation.module.css";

const Navigation = () => {
  const router = useRouter();

  const links = [
    { href: "/", label: "💬" },
    { href: "/contacts", label: "👥" },
    { href: "/settings", label: "⚙️" },
  ];

  return (
    <div className="controls-for-scrollable">
      {links.map((link) => (
        <Link href={link.href} key={link.href}>
          <a
            className={`${styles.navBtn} ${
              router.pathname === link.href ? styles.active : ""
            }`}
          >
            {link.label}
          </a>
        </Link>
      ))}
    </div>
  );
};

export default Navigation;
