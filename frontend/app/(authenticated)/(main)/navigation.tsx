import Link from "next/link";
import styles from "./layout.module.css";

const links = [
  { href: "/main", label: "💬" },
  { href: "/contacts", label: "👥" },
  { href: "/settings", label: "⚙️" },
];

export function Navigation({ current }: { current: string }) {
  return (
    <div className="controls-for-scrollable">
      {links.map((link) => (
        <Link
          href={link.href}
          key={link.href}
          className={`${styles.navBtn} shadow ${
            current === link.href ? styles.active : ""
          }`}
        >
          {link.label}
        </Link>
      ))}
    </div>
  );
}
